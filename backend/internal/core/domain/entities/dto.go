package entities

type ApiResponse struct {
	Success    bool                `json:"success"`
	Message    string              `json:"message"`
	Data       interface{}         `json:"data,omitempty"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
}

type PaginationResponse struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
}
