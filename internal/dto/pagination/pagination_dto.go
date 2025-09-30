package pagination_dto

type PaginationMetaDto struct {
	Page  int  `json:"page"`
	Limit int  `json:"limit"`
	Next  bool `json:"next"`
}

type PaginationDto struct {
	Object  string            `json:"object"`
	Records []interface{}     `json:"records"`
	Total   int               `json:"total"`
	Meta    PaginationMetaDto `json:"meta"`
}
