package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flight-api/config"
	sync_dto "flight-api/internal/dto/sync"
	repository_airport "flight-api/internal/repository/airport"
	service_aviation "flight-api/internal/service/aviation"
	service_sync "flight-api/internal/service/sync"
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
	util.PanicIfError(err)

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
	aviationService := service_aviation.NewAviationService(log, &cfg)
	syncService := service_sync.NewSyncService(
		log,
		validate,
		db,
		airportRepository,
		aviationService,
	)

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

	req := sync_dto.SyncAirportRequest{
		ICAOCodes: icaoCodes,
	}

	// Sync airports
	log.Info("Seeding data from file:", seedPath)

	start := time.Now()
	result, err := syncService.SyncAirports(ctx, req)

	if err != nil {
		log.Fatalf("Failed to sync seed data: %v", err)
	}

	diff := time.Since(start)
	log.Infof("✅ Successfully synced %d records in %v", len(result), diff)

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

	err = writeJSON(outPath, result)
	if err != nil {
		log.Fatalf("Failed write json %v", err)
	}

	log.Infof("✅ Successfully wrote %d records to %s", len(result), outPath)
	log.Infof("✅ Successfully Run Sync Seed Data")
}
