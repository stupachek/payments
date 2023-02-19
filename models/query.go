package models

type QueryParams struct {
	Limit  uint   `json:"limit"`
	Offset uint   `json:"offset"`
	Sort   string `json:"sort"`
}
