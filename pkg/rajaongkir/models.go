package rajaongkir

type GetFeeRequest struct {
	Origin      string
	Destination string
	Courier     string
}

type FeeResponse struct {
	Meta Meta   `json:"meta"`
	Data []Data `json:"data"`
}

type Data struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Service     string `json:"service"`
	Description string `json:"description"`
	Cost        int64  `json:"cost"`
	Etd         string `json:"etd"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int64  `json:"code"`
	Status  string `json:"status"`
}
