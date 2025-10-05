package enum

type FasilityTypeEnum *string

var (
	airportStr  = "airport"
	heliportStr = "heliport"
)

var (
	AIRPORT  FasilityTypeEnum = &airportStr
	HELIPORT FasilityTypeEnum = &heliportStr
	NIL      FasilityTypeEnum = nil
)

func ToFacilityType(s string) FasilityTypeEnum {
	switch s {
	case "airport":
		return AIRPORT
	case "heliport":
		return HELIPORT
	default:
		return NIL
	}
}
