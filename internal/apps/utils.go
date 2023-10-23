package core

import (
	"encoding/json"
	"net/http"
)

func NewPaginate(page int, limit int) *Paginate {
	startWith := limit*page - limit

	return &Paginate{
		StartWith: startWith,
		Limit:     limit,
	}
}

func SendResponse(writer http.ResponseWriter, data any, statusCode int) {
	resp := new(Response)

	resp.Response = data

	if statusCode < 400 {
		resp.Success = true
	} else {
		resp.Success = false
	}

	jsonData, _ := json.Marshal(resp)
	writer.WriteHeader(statusCode)
	writer.Write(jsonData)
}

func SendPaginateResponse(writer http.ResponseWriter, data any, paginate *Paginate) {
	resp := new(PaginateResponse)

	resp.Success = true
	resp.ResultFrom = paginate.StartWith
	resp.ResultTo = paginate.StartWith + paginate.Limit
	resp.Data = data

	jsonData, _ := json.Marshal(resp)
	writer.Write(jsonData)
}
