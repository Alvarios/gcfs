package responses

type Responses struct {
	Code   int32  `json:"code"`
	Status string `json:"status"`
	Body   string `json:"body"`
}
