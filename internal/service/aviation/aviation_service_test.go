package service_aviation

import (
	"context"
	"encoding/json"
	"errors"
	"flight-api/config"
	airport_dto "flight-api/internal/dto/airport"
	aviation_dto "flight-api/internal/dto/aviation"
	"flight-api/pkg/logger"
	"flight-api/util"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Dummy Data Aviation API
var DummyData = map[string][]map[string]string{
	"WIII": {},
	"WAHI": {},
	"KJFK": {
		{
			"site_number":              "15793.*A",
			"type":                     "AIRPORT",
			"facility_name":            "JOHN F KENNEDY INTL",
			"faa_ident":                "JFK",
			"icao_ident":               "KJFK",
			"region":                   "AEA",
			"district_office":          "NYC",
			"state":                    "NY",
			"state_full":               "NEW YORK",
			"county":                   "QUEENS",
			"city":                     "NEW YORK",
			"ownership":                "PU",
			"use":                      "PU",
			"manager":                  "CHARLES EVERETT",
			"manager_phone":            "(718) 244-3501",
			"latitude":                 "40-38-23.7400N",
			"latitude_sec":             "146303.7400N",
			"longitude":                "073-46-43.2930W",
			"longitude_sec":            "265603.2930W",
			"elevation":                "13",
			"magnetic_variation":       "13W",
			"tpa":                      "",
			"vfr_sectional":            "NEW YORK",
			"boundary_artcc":           "ZNY",
			"boundary_artcc_name":      "NEW YORK",
			"responsible_artcc":        "ZNY",
			"responsible_artcc_name":   "NEW YORK",
			"fss_phone_number":         "",
			"fss_phone_numer_tollfree": "1-800-WX-BRIEF",
			"notam_facility_ident":     "JFK",
			"status":                   "O",
			"certification_typedate":   "I E S 05/1973",
			"customs_airport_of_entry": "N",
			"military_joint_use":       "N",
			"military_landing":         "Y",
			"lighting_schedule":        "",
			"beacon_schedule":          "SS-SR",
			"control_tower":            "Y",
			"unicom":                   "122.950",
			"ctaf":                     "",
			"effective_date":           "11/04/2021",
		},
	},
	"KATL": {
		{
			"site_number":              "03640.*A",
			"type":                     "AIRPORT",
			"facility_name":            "HARTSFIELD - JACKSON ATLANTA INTL",
			"faa_ident":                "ATL",
			"icao_ident":               "KATL",
			"region":                   "ASO",
			"district_office":          "ATL",
			"state":                    "GA",
			"state_full":               "GEORGIA",
			"county":                   "FULTON",
			"city":                     "ATLANTA",
			"ownership":                "PU",
			"use":                      "PU",
			"manager":                  "BALRAM BHEODARI",
			"manager_phone":            "404-530-6600",
			"latitude":                 "33-38-12.1186N",
			"latitude_sec":             "121092.1186N",
			"longitude":                "084-25-40.3104W",
			"longitude_sec":            "303940.3104W",
			"elevation":                "1026",
			"magnetic_variation":       "05W",
			"tpa":                      "",
			"vfr_sectional":            "ATLANTA",
			"boundary_artcc":           "ZTL",
			"boundary_artcc_name":      "ATLANTA",
			"responsible_artcc":        "ZTL",
			"responsible_artcc_name":   "ATLANTA",
			"fss_phone_number":         "",
			"fss_phone_numer_tollfree": "1-800-WX-BRIEF",
			"notam_facility_ident":     "ATL",
			"status":                   "O",
			"certification_typedate":   "I E S 05/1973",
			"customs_airport_of_entry": "N",
			"military_joint_use":       "N",
			"military_landing":         "Y",
			"lighting_schedule":        "",
			"beacon_schedule":          "SS-SR",
			"control_tower":            "Y",
			"unicom":                   "122.950",
			"ctaf":                     "",
			"effective_date":           "11/04/2021",
		},
	},
}

var InvalidJSON1 = []byte("not-json")
var InvalidJSON2 = []byte(`{"KKKK":[{"site_number":"15793.*A"}`)

// Mock Aviation API
func newAviationMockServer(routes map[string][]map[string]string, delay time.Duration) *httptest.Server {
	log := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	log.Info("Mock Aviation API server started")

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apt := r.URL.Query().Get("apt")

		// Handle not has apt
		if apt == "" {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(
				map[string]string{
					"status":      "error",
					"status_code": "404",
					"message":     "Not Found",
				},
			)
			return
		}

		// Handle has apt
		resBody := make(map[string][]map[string]string)

		if delay > 0 {
			time.Sleep(delay)
		}

		w.Header().Set("Content-Type", "application/json")

		switch strings.ToUpper(apt) {
		case "KKKK":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(InvalidJSON2)
			return
		case "BADJSON":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(InvalidJSON1)
			return
		case "BADREQUEST":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(InvalidJSON1)
			return
		}

		// split apt
		codes := strings.Split(apt, ",")

		for _, code := range codes {
			code = strings.ToUpper(strings.TrimSpace(code))
			if code == "" {
				resBody[code] = []map[string]string{}
				continue
			}

			switch code {
			case "KKKK":
				resBody[code] = []map[string]string{}
				continue
			case "BADJSON":
				resBody[code] = []map[string]string{}
				continue
			}

			resp, ok := routes[code]
			if !ok {
				resBody[code] = []map[string]string{}
				continue
			}

			resBody[code] = resp
		}

		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(resBody)
	})

	return httptest.NewServer(h)
}

// Mock Tripper
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

type errReadCloser struct{ err error }

func (e errReadCloser) Read(p []byte) (int, error) { return 0, e.err }
func (e errReadCloser) Close() error               { return nil }

// ---- Test Initialize Aviation Service ----
func TestNewAviationService_InitializesFields(t *testing.T) {
	// Arrange
	log := logger.NewLogger(logger.INFO_DEBUG_LEVEL)
	cfg := &config.Config{AviationURL: "http://example.com"}

	// Act
	svc := NewAviationService(log, cfg)

	// Assert 1: implement IAviationService (compile-time + runtime)
	var _ IAviationService = svc // compile-time assertion

	impl, ok := svc.(*AviationService)
	require.True(t, ok, "should return *AviationService as IAviationService")

	// Assert 2: fields terpasang benar
	require.NotNil(t, impl)
	require.Equal(t, log, impl.logger) // kalau ada getter; atau langsung impl.logger kalau di paket sama
	require.Equal(t, cfg, impl.cfg)    // kalau ada getter/field exported; sesuaikan aksesnya
}

// ---- Test FetchAirportData ----
func TestFetchAirportData_Success(t *testing.T) {
	srv := newAviationMockServer(DummyData, 1*time.Second)
	defer srv.Close()

	svc := &AviationService{
		cfg: &config.Config{
			AviationURL: srv.URL,
		},
		logger: logger.NewLogger(logger.INFO_DEBUG_LEVEL),
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}

	tests := []struct {
		name     string
		request  []string
		expected map[string]airport_dto.AirportRequestDto
	}{
		{
			name: "success - 1 ICAO code - a",
			request: []string{
				"KJFK",
			},
			expected: map[string]airport_dto.AirportRequestDto{
				"KJFK": aviation_dto.ToAirportRequestDto(
					aviation_dto.AviationAirportDto{
						SiteNumber:              "15793.*A",
						Type:                    "AIRPORT",
						FacilityName:            "JOHN F KENNEDY INTL",
						FAAIdentifier:           "JFK",
						ICAOIdentifier:          "KJFK",
						Region:                  "AEA",
						DistrictOffice:          "NYC",
						State:                   "NY",
						StateFull:               "NEW YORK",
						County:                  "QUEENS",
						City:                    "NEW YORK",
						Ownership:               "PU",
						Use:                     "PU",
						Manager:                 "CHARLES EVERETT",
						ManagerPhone:            "(718) 244-3501",
						Latitude:                "40-38-23.7400N",
						LatitudeSec:             "146303.7400N",
						Longitude:               "073-46-43.2930W",
						LongitudeSec:            "265603.2930W",
						Elevation:               "13",
						MagneticVariation:       "13W",
						TPA:                     "",
						VFRSectional:            "NEW YORK",
						NotamFacilityIdentifier: "JFK",
						Status:                  "O",
						ControlTower:            "Y",
						UNICOM:                  "122.950",
						CTAF:                    "",
						EffectiveDate:           "11/04/2021",
					},
				),
			},
		},
		{
			name: "success - 1 ICAO code - b",
			request: []string{
				"KATL",
			},
			expected: map[string]airport_dto.AirportRequestDto{
				"KATL": aviation_dto.ToAirportRequestDto(
					aviation_dto.AviationAirportDto{
						SiteNumber:              "03640.*A",
						Type:                    "AIRPORT",
						FacilityName:            "HARTSFIELD - JACKSON ATLANTA INTL",
						FAAIdentifier:           "ATL",
						ICAOIdentifier:          "KATL",
						Region:                  "ASO",
						DistrictOffice:          "ATL",
						State:                   "GA",
						StateFull:               "GEORGIA",
						County:                  "FULTON",
						City:                    "ATLANTA",
						Ownership:               "PU",
						Use:                     "PU",
						Manager:                 "BALRAM BHEODARI",
						ManagerPhone:            "404-530-6600",
						Latitude:                "33-38-12.1186N",
						LatitudeSec:             "121092.1186N",
						Longitude:               "084-25-40.3104W",
						LongitudeSec:            "303940.3104W",
						Elevation:               "1026",
						MagneticVariation:       "05W",
						TPA:                     "",
						VFRSectional:            "ATLANTA",
						NotamFacilityIdentifier: "ATL",
						Status:                  "O",
						ControlTower:            "Y",
						UNICOM:                  "122.950",
						CTAF:                    "",
						EffectiveDate:           "11/04/2021",
					},
				),
			},
		},
		{
			name: "success - 1 ICAO code - c",
			request: []string{
				"WIII",
			},
			expected: map[string]airport_dto.AirportRequestDto{
				"WIII": {},
			},
		},
		{
			name: "success - 1 ICAO code - d",
			request: []string{
				"WAHI",
			},
			expected: map[string]airport_dto.AirportRequestDto{
				"WAHI": {},
			},
		},
		{
			name: "success - 2 ICAO codes",
			request: []string{
				"KATL",
				"WIII",
			},
			expected: map[string]airport_dto.AirportRequestDto{
				"KATL": aviation_dto.ToAirportRequestDto(
					aviation_dto.AviationAirportDto{
						SiteNumber:              "03640.*A",
						Type:                    "AIRPORT",
						FacilityName:            "HARTSFIELD - JACKSON ATLANTA INTL",
						FAAIdentifier:           "ATL",
						ICAOIdentifier:          "KATL",
						Region:                  "ASO",
						DistrictOffice:          "ATL",
						State:                   "GA",
						StateFull:               "GEORGIA",
						County:                  "FULTON",
						City:                    "ATLANTA",
						Ownership:               "PU",
						Use:                     "PU",
						Manager:                 "BALRAM BHEODARI",
						ManagerPhone:            "404-530-6600",
						Latitude:                "33-38-12.1186N",
						LatitudeSec:             "121092.1186N",
						Longitude:               "084-25-40.3104W",
						LongitudeSec:            "303940.3104W",
						Elevation:               "1026",
						MagneticVariation:       "05W",
						TPA:                     "",
						VFRSectional:            "ATLANTA",
						NotamFacilityIdentifier: "ATL",
						Status:                  "O",
						ControlTower:            "Y",
						UNICOM:                  "122.950",
						CTAF:                    "",
						EffectiveDate:           "11/04/2021",
					},
				),
				"WIII": {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.FetchAirportData(context.Background(), tt.request)
			assert.NoError(t, err)

			for expectedCode, expectedAirport := range tt.expected {
				actualAirport, ok := result[expectedCode]
				assert.True(t, ok)

				v1 := reflect.ValueOf(expectedAirport)
				v2 := reflect.ValueOf(actualAirport)

				for i := 0; i < v1.NumField(); i++ {
					f1 := v1.Field(i)
					f2 := v2.Field(i)

					assert.IsType(t, f1.Interface(), f2.Interface(), "field %s not same type", v1.Type().Field(i).Name)
					assert.Equal(t, f1.Interface(), f2.Interface(), "field %s not equal", v1.Type().Field(i).Name)
				}
			}
		})
	}
}

func TestFetchAirportData_BadRequest(t *testing.T) {
	srv := newAviationMockServer(map[string][]map[string]string{
		"BADREQUEST": {},
	}, 0)
	defer srv.Close()

	svc := &AviationService{
		cfg: &config.Config{
			AviationURL: srv.URL,
		},
		logger: logger.NewLogger(logger.DEBUG_LEVEL),
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}

	_, err := svc.FetchAirportData(context.Background(), []string{"BADREQUEST"})
	assert.ErrorIs(t, err, util.ErrBadRequest)
}

func TestFetchAirportData_InternalServerError(t *testing.T) {
	srv := newAviationMockServer(map[string][]map[string]string{
		"KKKK": {},
	}, 0)
	defer srv.Close()

	svc := &AviationService{
		cfg: &config.Config{
			AviationURL: srv.URL,
		},
		logger: logger.NewLogger(logger.DEBUG_LEVEL),
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}

	_, err := svc.FetchAirportData(context.Background(), []string{"KKKK"})
	assert.ErrorIs(t, err, util.ErrInternalServer)
}

func TestFetchAirportData_Error_Timeout(t *testing.T) {
	srv := newAviationMockServer(DummyData, 2*time.Second)
	defer srv.Close()

	svc := &AviationService{
		cfg: &config.Config{
			AviationURL: srv.URL,
		},
		logger: logger.NewLogger(logger.DEBUG_LEVEL),
		client: &http.Client{
			Timeout: 50 * time.Millisecond,
		},
	}

	_, err := svc.FetchAirportData(context.Background(), []string{"KJFK"})
	assert.ErrorIs(t, err, util.ErrGatewayTimeout)
}

func TestFetchAirportData_Error_InternalServer(t *testing.T) {
	// Transport mengembalikan error biasa (bukan *url.Error, bukan net.Error)
	rt := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("weird client error")
	})
	client := &http.Client{Transport: rt}

	svc := &AviationService{
		cfg:    &config.Config{AviationURL: "http://mock"},
		client: client,
		logger: logger.NewLogger(logger.INFO_DEBUG_LEVEL),
	}

	_, err := svc.FetchAirportData(context.Background(), []string{"WIII"})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, util.ErrInternalServer) {
		t.Fatalf("expected ErrInternalServer, got %v", err)
	}
}

func TestFetchAirportData_ReadBodyError(t *testing.T) {
	rt := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       errReadCloser{err: io.ErrUnexpectedEOF}, // simulasi body korup/terputus
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})

	client := &http.Client{Transport: rt}

	svc := &AviationService{
		cfg:    &config.Config{AviationURL: "http://mock"},
		client: client,
		logger: logger.NewLogger(logger.INFO_DEBUG_LEVEL),
	}

	_, err := svc.FetchAirportData(context.Background(), []string{"WIII"})
	if err == nil {
		t.Fatalf("expected error")
	}

	if !errors.Is(err, util.ErrInternalServer) {
		t.Fatalf("expected ErrInternalServer, got %v", err)
	}
}
