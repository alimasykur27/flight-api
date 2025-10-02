package location_dto

type LocationDto struct {
	Name           *string  `json:"name"`
	Eegion         *string  `json:"region"`
	Country        *string  `json:"country"`
	Lat            *float64 `json:"lat"`
	Lon            *float64 `json:"lon"`
	TzId           *string  `json:"tz_id"`
	LocaltimeEpoch *int     `json:"localtime_epoch"`
	Localtime      *string  `json:"localtime"`
}
