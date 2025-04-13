package models

type ListWithTotal[T any] struct {
	Total int `form:"total" json:"total"`
	Items []T `form:"items" json:"items"`
}
