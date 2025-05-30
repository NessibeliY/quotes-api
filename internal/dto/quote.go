package dto

import "fmt"

type AddQuoteRequest struct {
	Author string `json:"author"`
	Quote  string `json:"quote"`
}

func (r *AddQuoteRequest) Validate() error {
	if r.Author == "" {
		return fmt.Errorf("author is empty")
	}
	if r.Quote == "" {
		return fmt.Errorf("quote is empty")
	}
	return nil
}
