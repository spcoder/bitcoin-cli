package main

import "github.com/parnurzeal/gorequest"

// HTTPRequester is an interface for getting BTC ticker data
type HTTPRequester interface {
	Get() (gorequest.Response, string, []error)
}

// HTTPManager is a type for getting BTC ticker data
type HTTPManager struct {
}

// Get retrieves BTC ticker data from gdax
func (HTTPManager) Get() (gorequest.Response, string, []error) {
	return gorequest.New().Get("https://api.gdax.com/products/BTC-USD/ticker").End()
}
