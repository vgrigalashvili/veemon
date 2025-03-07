package dto

// StandardResponse represents a generic API response.
type StandardResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}
