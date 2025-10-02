package weather_dto

import location_dto "flight-api/internal/dto/location"

type CurrentWeatherDto struct {
	LastUpdatedEpoch *int     `json:"last_updated_epoch"`
	LastUpdated      *string  `json:"last_updated"`
	TempC            *float64 `json:"temp_c"`
	TempF            *float64 `json:"temp_f"`
	IsDay            *uint8   `json:"is_day"`
	Condition        *struct {
		Text *string `json:"text"`
		Icon *string `json:"icon"`
		Code *int    `json:"code"`
	} `json:"condition"`
	WindMph    *float64 `json:"wind_mph"`
	WindKph    *float64 `json:"wind_kph"`
	WindDegree *int     `json:"wind_degree"`
	WindDir    *string  `json:"wind_dir"`
	PressureMb *float64 `json:"pressure_mb"`
	PressureIn *float64 `json:"pressure_in"`
	PrecipMm   *float64 `json:"precip_mm"`
	PrecipIn   *float64 `json:"precip_in"`
	Humidity   *int     `json:"humidity"`
	Cloud      *int     `json:"cloud"`
	FeelslikeC *float64 `json:"feelslike_c"`
	FeelslikeF *float64 `json:"feelslike_f"`
	WindchillC *float64 `json:"windchill_c"`
	WindchillF *float64 `json:"windchill_f"`
	HeatindexC *float64 `json:"heatindex_c"`
	HeatindexF *float64 `json:"heatindex_f"`
	DewpointC  *float64 `json:"dewpoint_c"`
	DewpointF  *float64 `json:"dewpoint_f"`
	VisKm      *float64 `json:"vis_km"`
	VisMiles   *float64 `json:"vis_miles"`
	Uv         *float64 `json:"uv"`
	GustMph    *float64 `json:"gust_mph"`
	GustKph    *float64 `json:"gust_kph"`
	ShortRad   *float64 `json:"short_rad"`
	DiffRad    *float64 `json:"diff_rad"`
	DNI        *float64 `json:"dni"`
	GTI        *float64 `json:"gti"`
}

type WeatherDto struct {
	Object   *string                   `json:"object"`
	Location *location_dto.LocationDto `json:"location"`
	Current  *CurrentWeatherDto        `json:"current"`
}
