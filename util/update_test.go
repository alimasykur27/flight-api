package util_test

import (
	airport_dto "flight-api/internal/dto/airport"
	"flight-api/internal/enum"
	"flight-api/internal/model"
	"flight-api/util"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// --- helpers ---

// --- UpdateString ---

func TestUpdateString(t *testing.T) {
	t.Run("valid - update with non-nil src", func(t *testing.T) {
		dst := new(string)
		*dst = "old"

		src := "new"
		util.UpdateString(&dst, &src)

		assert.Equal(t, "new", *dst)
	})

	t.Run("valid - src is nil (dst unchanged)", func(t *testing.T) {
		dst := new(string)
		*dst = "keep"

		util.UpdateString(&dst, nil)

		assert.Equal(t, "keep", *dst) // tidak berubah
	})

	t.Run("invalid - wrong dereference level", func(t *testing.T) {
		str := "oops"
		// dst seharusnya **string, tapi kita coba pakai *string
		// ini tidak akan compile kalau langsung dipanggil,
		// jadi kita uji dengan interface{} + reflect, atau tandai sebagai negative test.
		// Di sini cukup tunjukkan bahwa tidak bisa compile.
		_ = str
		// 👉 Compile-time error: UpdateString(str, &str) // uncomment to see the error
	})
}

// --- UpdateBool ---
func TestUpdateBool(t *testing.T) {
	t.Run("valid - update with non-nil src", func(t *testing.T) {
		dst := new(bool)
		*dst = false

		src := true
		util.UpdateBool(&dst, &src)

		if *dst != true {
			t.Fatalf("expected dst=true, got %v", *dst)
		}
	})

	t.Run("valid - src is nil (dst unchanged)", func(t *testing.T) {
		dst := new(bool)
		*dst = false

		util.UpdateBool(&dst, nil)

		if *dst != false {
			t.Fatalf("expected dst remain false, got %v", *dst)
		}
	})

	t.Run("invalid - dst is nil pointer to pointer", func(t *testing.T) {
		var dstPtr **bool = nil
		src := true

		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic, got none")
			}
		}()

		util.UpdateBool(dstPtr, &src) // panic karena dstPtr == nil
	})
}

// --- TestInt ---

func TestUpdateInt(t *testing.T) {
	t.Run("valid - update with non-nil src - int", func(t *testing.T) {
		dst := []*int{new(int), nil}

		src := 42

		// dst has value
		*dst[0] = 2
		old := *dst[0]
		util.UpdateInt(&dst[0], &src)
		assert.IsType(t, 42, *dst[0])
		assert.Equal(t, 42, *dst[0])
		assert.NotEqual(t, old, *dst[0])

		// dst is nil
		util.UpdateInt(&dst[1], &src)
		assert.IsType(t, 42, *dst[1])
		assert.Equal(t, 42, *dst[1])
		assert.NotEqual(t, nil, *dst[1])
	})

	t.Run("valid - update with non-nil src - int32", func(t *testing.T) {
		dst := []*int32{new(int32), nil}

		src := int32(42)

		// dst has value
		*dst[0] = 2
		old := *dst[0]
		util.UpdateInt(&dst[0], &src)
		assert.IsType(t, int32(42), *dst[0])
		assert.Equal(t, int32(42), *dst[0])
		assert.NotEqual(t, old, *dst[0])

		// dst is nil
		util.UpdateInt(&dst[1], &src)
		assert.IsType(t, int32(42), *dst[1])
		assert.Equal(t, int32(42), *dst[1])
		assert.NotEqual(t, nil, *dst[1])
	})

	t.Run("valid - update with non-nil src - int64", func(t *testing.T) {
		dst := []*int64{new(int64), nil}

		src := int64(410000123456789)

		// dst has value
		*dst[0] = 2
		old := *dst[0]
		util.UpdateInt(&dst[0], &src)
		assert.IsType(t, int64(410000123456789), *dst[0])
		assert.Equal(t, int64(410000123456789), *dst[0])
		assert.NotEqual(t, old, *dst[0])

		// dst is nil
		util.UpdateInt(&dst[1], &src)
		assert.IsType(t, int64(410000123456789), *dst[1])
		assert.Equal(t, int64(410000123456789), *dst[1])
		assert.NotEqual(t, nil, *dst[1])
	})

	t.Run("valid - src is nil (dst unchanged)", func(t *testing.T) {
		dst := new(int)
		*dst = 42

		util.UpdateInt(&dst, nil)

		assert.IsType(t, 42, *dst)
		assert.Equal(t, 42, *dst)
		assert.NotEqual(t, nil, *dst)
	})

	t.Run("invalid - dst is nil pointer to pointer", func(t *testing.T) {
		var dstPtr **int = nil
		src := 42

		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic, got none")
			}
		}()

		util.UpdateInt(dstPtr, &src)

		assert.IsType(t, nil, dstPtr)
		assert.Equal(t, nil, dstPtr)
	})
}

func TestUpdateTime(t *testing.T) {
	t.Run("valid - update with non-nil src", func(t *testing.T) {
		now := time.Now()
		dst := &now

		src := time.Now().Add(2 * time.Hour)
		util.UpdateTime(&dst, &src)

		assert.IsType(t, src, *dst)
		assert.Equal(t, src, *dst)
		assert.NotEqual(t, now, *dst)
	})

	t.Run("valid - src is nil (dst unchanged)", func(t *testing.T) {
		now := time.Now()
		dst := &now

		util.UpdateTime(&dst, nil)

		assert.IsType(t, now, *dst)
		assert.Equal(t, now, *dst)
		assert.NotEqual(t, nil, *dst)
	})

	t.Run("invalid - dst is nil pointer to pointer", func(t *testing.T) {
		var dstPtr **time.Time = nil
		src := time.Now()

		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic, got none")
			}
		}()

		util.UpdateTime(dstPtr, &src)

		assert.IsType(t, nil, dstPtr)
		assert.Equal(t, nil, dstPtr)
	})
}

// helper: deep copy airport
func deepCopyAirport(a model.Airport) model.Airport {
	cp := a
	return cp
}

func TestFillUpdatableFields(t *testing.T) {
	newId := uuid.New()
	timeNow := time.Now()

	existingAirport := model.Airport{
		ID:            &newId,
		SiteNumber:    util.Ptr("67890"),
		ICAOID:        util.Ptr("KLAX"),
		FAAID:         util.Ptr("FAA1"),
		IATAID:        util.Ptr("LAX"),
		Name:          util.Ptr("Los Angeles Intl"),
		Type:          enum.AIRPORT, // enum (underlying string)
		Status:        util.Ptr(true),
		Country:       util.Ptr("US"),
		State:         util.Ptr("CA"),
		StateFull:     util.Ptr("California"),
		County:        util.Ptr("Los Angeles"),
		City:          util.Ptr("Los Angeles"),
		Ownership:     enum.OWN_PUBLIC, // enum (underlying string)
		Use:           enum.USE_PUBLIC, // enum (underlying string)
		Manager:       util.Ptr("John Smith"),
		ManagerPhone:  util.Ptr("+1-555-0200"),
		Latitude:      util.Ptr("33.9416"),
		LatitudeSec:   util.Ptr("00.0"),
		Longitude:     util.Ptr("-118.4085"),
		LongitudeSec:  nil,
		Elevation:     util.Ptr(int64(125)),
		ControlTower:  util.Ptr(true),
		Unicom:        util.Ptr("122.80"),
		CTAF:          util.Ptr("119.80"),
		EffectiveDate: nil,
		CreatedAt:     &timeNow,
		UpdatedAt:     &timeNow,
	}

	type want struct {
		// isi hanya yang kamu mau verifikasi berubah
		Manager       *string
		Type          enum.FasilityTypeEnum
		Ownership     enum.UseTypeEnum
		Use           enum.UseTypeEnum
		Status        *bool
		LatitudeSec   *string
		Elevation     *int64
		EffectiveDate *time.Time
	}

	tests := []struct {
		name         string
		payload      airport_dto.AirportUpdateDto
		want         want
		verifyOthers bool // jika true, pastikan field selain yang diubah tetap sama
	}{
		{
			name: "only Manager changes",
			payload: airport_dto.AirportUpdateDto{
				Manager: util.Ptr("Manager New"),
			},
			want:         want{Manager: util.Ptr("Manager New")},
			verifyOthers: true,
		},
		{
			name: "change enums and bool (Type, Ownership, Use, Status)",
			payload: func() airport_dto.AirportUpdateDto {
				typ := enum.HELIPORT // contoh enum lain
				own := enum.OWN_PRIVATE
				use := enum.USE_PRIVATE
				return airport_dto.AirportUpdateDto{
					Type:      typ,
					Ownership: own,
					Use:       use,
					Status:    util.Ptr(false),
				}
			}(),
			want: want{
				Type:      enum.HELIPORT,
				Ownership: enum.USE_PRIVATE,
				Use:       enum.USE_PRIVATE,
				Status:    util.Ptr(false),
			},
			verifyOthers: true,
		},
		{
			name: "set previously-nil field (LatitudeSec) and numeric (Elevation)",
			payload: airport_dto.AirportUpdateDto{
				LatitudeSec: util.Ptr("59.9"),
				Elevation:   util.Ptr(int64(130)),
			},
			want: want{
				LatitudeSec: util.Ptr("59.9"),
				Elevation:   util.Ptr(int64(130)),
			},
			verifyOthers: true,
		},
		{
			name: "set EffectiveDate",
			payload: func() airport_dto.AirportUpdateDto {
				ts := timeNow.Add(24 * time.Hour)
				return airport_dto.AirportUpdateDto{
					EffectiveDate: &ts,
				}
			}(),
			want: want{
				EffectiveDate: util.Ptr(timeNow.Add(24 * time.Hour)),
			},
			verifyOthers: true,
		},
		{
			name:    "nil payload (no-op)",
			payload: airport_dto.AirportUpdateDto{
				// semua nil
			},
			want:         want{}, // tidak ada yang berubah
			verifyOthers: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := deepCopyAirport(existingAirport)

			// ACT
			util.FillUpdatableFields(&got, tc.payload)

			// ASSERT: field yang diharapkan berubah
			if tc.want.Manager != nil {
				assert.Equal(t, *tc.want.Manager, util.DerefPtr(got.Manager))
			}
			if tc.want.Type != nil {
				// model.Airport.Type adalah enum (underlying string)
				assert.Equal(t, *tc.want.Type, *got.Type)
			}
			if tc.want.Ownership != nil {
				assert.Equal(t, *tc.want.Ownership, *got.Ownership)
			}
			if tc.want.Use != nil {
				assert.Equal(t, *tc.want.Use, *got.Use)
			}
			if tc.want.Status != nil {
				assert.Equal(t, *tc.want.Status, util.DerefPtr(got.Status))
			}
			if tc.want.LatitudeSec != nil {
				assert.Equal(t, *tc.want.LatitudeSec, util.DerefPtr(got.LatitudeSec))
			}
			if tc.want.Elevation != nil {
				assert.Equal(t, *tc.want.Elevation, util.DerefPtr(got.Elevation))
			}
			if tc.want.EffectiveDate != nil {
				assert.WithinDuration(t, *tc.want.EffectiveDate, util.DerefPtr(got.EffectiveDate), time.Second)
			}

			// ASSERT: yang lain tidak berubah (optional)
			if tc.verifyOthers {
				// contoh beberapa spot-check field penting yang tidak ikut berubah
				assert.Equal(t, util.DerefPtr(existingAirport.SiteNumber), util.DerefPtr(got.SiteNumber))
				assert.Equal(t, util.DerefPtr(existingAirport.ICAOID), util.DerefPtr(got.ICAOID))
				assert.Equal(t, util.DerefPtr(existingAirport.FAAID), util.DerefPtr(got.FAAID))
				assert.Equal(t, util.DerefPtr(existingAirport.IATAID), util.DerefPtr(got.IATAID))
				if tc.want.Manager == nil {
					assert.Equal(t, util.DerefPtr(existingAirport.Manager), util.DerefPtr(got.Manager))
				}
				// CreatedAt/UpdatedAt tidak diubah oleh FillUpdatableFields
				assert.WithinDuration(t, *existingAirport.CreatedAt, *got.CreatedAt, time.Second)
				assert.WithinDuration(t, *existingAirport.UpdatedAt, *got.UpdatedAt, time.Second)
			}
		})
	}
}
