package core

type Response struct {
	Success  bool `json:"success"`
	Response any  `json:"response"`
}

type PaginateResponse struct {
	Success    bool `json:"success"`
	ResultFrom int  `json:"result_from"`
	ResultTo   int  `json:"result_to"`
	Data       any  `json:"data"`
}

type Paginate struct {
	StartWith int
	Limit     int
}
