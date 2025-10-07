package weather_dto

type ConditionDto struct {
	Text *string `json:"text"`
	Icon *string `json:"icon"`
	Code *int    `json:"code"`
}
