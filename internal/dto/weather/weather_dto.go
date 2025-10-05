package weather_dto

import location_dto "flight-api/internal/dto/location"

type WeatherDto struct {
	Object   *string                   `json:"object"`
	Location *location_dto.LocationDto `json:"location"`
	Current  *CurrentWeatherDto        `json:"current"`
}
