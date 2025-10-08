package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flight-api/config"
	airport_dto "flight-api/internal/dto/airport"
	repository_airport "flight-api/internal/repository/airport"
	service_airport "flight-api/internal/service/airport"
	service_weather "flight-api/internal/service/weather"
	"flight-api/pkg/database"
	"flight-api/pkg/logger"
	"flight-api/util"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// readLines reads a text file into a slice of strings (trimmed, skip blank & comments).
func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	var out []string
	sc := bufio.NewScanner(f)
	const maxCap = 1024 * 1024 // 1MB per line guard
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, maxCap)

	for lineNo := 1; sc.Scan(); lineNo++ {
		s := strings.TrimSpace(sc.Text())
		if s == "" {
			continue
		}
		// optional: skip comment lines starting with '#'
		if strings.HasPrefix(s, "#") {
			continue
		}
		out = append(out, s)
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("scan %s: %w", path, err)
	}
	if len(out) == 0 {
		return nil, errors.New("seed file is empty after trimming")
	}
	return out, nil
}

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func main() {
	// Create a context that will be canceled on interrupt signals
	log := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Info("Starting sync seed data...")

	// Initialize Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	// Connect to database
	log.Info("Connecting to database ...")
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return
	} else {
		log.Info("Succesfully connected to database!")
	}
	defer db.Close()

	// Initialize validator
	validate := util.NewValidator()

	// Initialize repository
	airportRepository := repository_airport.NewAirportRepository(log)

	// Initialize service
	weatherService := service_weather.NewWeatherService(log, &cfg)
	airportService := service_airport.NewAirportService(log, validate, db, airportRepository, weatherService)

	// Read seed data from file
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
		return
	}

	migrationsPath := filepath.Join(wd, "migrations")
	seedPath := filepath.Join(migrationsPath, "seed.txt")

	if len(os.Args) > 1 && os.Args[1] != "" {
		seedPath = os.Args[1]
	}

	icaoCodes, err := readLines(seedPath)
	if err != nil {
		log.Fatalf("Failed read seed file: %v", err)
	}
	log.Infof("✅ Successfully read %d ICAO codes from %s", len(icaoCodes), seedPath)

	// Sync airports
	log.Info("Seeding data from file:", seedPath)

	start := time.Now()
	// Use a goroutine to avoid blocking main thread
	result := make(chan []airport_dto.AirportDto)
	defer close(result)

	go func() {
		res, err := airportService.Seeding(ctx, icaoCodes)
		if err != nil {
			log.Errorf("Failed to seed data: %v", err)
			result <- nil
		}
		result <- res
	}()

	var airports []airport_dto.AirportDto = <-result

	if err != nil {
		log.Fatalf("Failed to sync seed data: %v", err)
	}

	diff := time.Since(start)
	log.Infof("Synced %d records in %v", len(airports), diff)

	// Write to JSON file
	binPath := filepath.Join(wd, "bin")
	outPath := filepath.Join(binPath, "seed_result.json")

	// Check if binPath exists, if not create it
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		err := os.MkdirAll(binPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed create bin directory: %v", err)
		}
	}

	err = writeJSON(outPath, airports)
	if err != nil {
		log.Fatalf("Failed write json %v", err)
	}

	log.Infof("✅ Successfully wrote %d records to %s", len(airports), outPath)
	log.Infof("✅ Successfully Run Seeding Data")
}
