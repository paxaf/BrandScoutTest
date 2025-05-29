package entity

type Quote struct {
	Author string `json:"author"`
	Phrase string `json:"quote"`
}

type QuoteResponse struct {
	Quotes []Quote `json:"quotes"`
}
