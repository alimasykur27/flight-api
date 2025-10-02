package airport_dto

import weather_dto "flight-api/internal/dto/weather"

type AirportWeatherDto struct {
	Object  string                         `json:"object"`
	Code    *string                        `json:"code"`
	Airport *AirportDto                    `json:"airport"`
	Weather *weather_dto.CurrentWeatherDto `json:"weather"`
}
