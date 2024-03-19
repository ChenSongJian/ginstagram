package utils

type PageResponse struct {
	PageNum      int         `json:"page_num"`
	PageSize     int         `json:"page_size"`
	TotalPages   int         `json:"total_pages"`
	TotalRecords int         `json:"total_records"`
	Data         interface{} `json:"data"`
}
