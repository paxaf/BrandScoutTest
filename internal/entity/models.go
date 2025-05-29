package entity

type Quote struct {
	Id     string `json:"id"`
	Author string `json:"author"`
	Phrase string `json:"quote"`
}

type QuoteResponse struct {
	Quotes []Quote `json:"quotes"`
}
