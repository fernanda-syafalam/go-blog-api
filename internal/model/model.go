package model

type WebResponse[T any] struct {
	Data T  `json:"data"`
	Paging *PageMetadata `json:"paging,omitempty"`
	Errors string `json:"errors,omitempty"`
}

type PageResponse[T any] struct {
	Data []T `json:"data"`
	PageMetadata *PageMetadata `json:"paging,omitempty"`
}

type PageMetadata struct {
	Page int `json:"page"`
	Size int `json:"size"`
	TotalItem int `json:"totalItem"`
	TotalPage int `json:"totalPage"`
}