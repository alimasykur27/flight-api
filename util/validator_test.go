package util_test

import (
	"flight-api/internal/enum"
	"flight-api/util"
	"testing"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Sample struct yang memakai tag custom kamu.
type FacilityInput struct {
	FacilityType string `validate:"facility"`
	Ownership    string `validate:"ownership"`
	Use          string `validate:"use"`
}

func TestNewValidator_ReturnsValidate(t *testing.T) {
	t.Parallel()

	v := util.NewValidator()

	require.NotNil(t, v, "validator should not be nil")
	_, ok := interface{}(v).(*validator.Validate)
	assert.True(t, ok, "should return *validator.Validate")
}

func TestStructValidation_ValidCases(t *testing.T) {
	t.Parallel()

	v := util.NewValidator()

	t.Run("all_empty_values_are_valid", func(t *testing.T) {
		in := FacilityInput{
			FacilityType: "",
			Ownership:    "",
			Use:          "",
		}
		err := v.Struct(in)
		assert.NoError(t, err)
	})

	t.Run("all_allowed_enum_values_are_valid", func(t *testing.T) {
		in := FacilityInput{
			FacilityType: string(*enum.AIRPORT),
			Ownership:    string(*enum.OWN_PUBLIC),
			Use:          string(*enum.USE_PRIVATE),
		}
		err := v.Struct(in)
		assert.NoError(t, err)
	})
}

func TestStructValidation_InvalidCases(t *testing.T) {
	t.Parallel()

	v := util.NewValidator()

	t.Run("single_invalid_field", func(t *testing.T) {
		in := FacilityInput{
			FacilityType: "SPACEPORT", // bukan salah satu dari enum
			Ownership:    "",          // valid (kosong diperbolehkan)
			Use:          "",          // valid
		}
		err := v.Struct(in)
		if assert.Error(t, err) {
			verrs, ok := err.(validator.ValidationErrors)
			require.True(t, ok, "should be ValidationErrors")
			assert.Len(t, verrs, 1)
			assert.Equal(t, "FacilityType", verrs[0].Field())
			assert.Equal(t, "facility", verrs[0].Tag())
		}
	})

	t.Run("multiple_invalid_fields", func(t *testing.T) {
		in := FacilityInput{
			FacilityType: "HELIPAD",    // invalid (misal enum hanya AIRPORT/HELIPORT)
			Ownership:    "GOV",        // invalid
			Use:          "COMMERCIAL", // invalid
		}
		err := v.Struct(in)
		if assert.Error(t, err) {
			verrs, ok := err.(validator.ValidationErrors)
			require.True(t, ok, "should be ValidationErrors")
			// Bisa 3 error, urutan bisa beda—pakai set-like checks
			fields := map[string]string{}
			for _, e := range verrs {
				fields[e.Field()] = e.Tag()
			}
			assert.Equal(t, "facility", fields["FacilityType"])
			assert.Equal(t, "ownership", fields["Ownership"])
			assert.Equal(t, "use", fields["Use"])
			assert.Equal(t, 3, len(fields))
		}
	})
}
