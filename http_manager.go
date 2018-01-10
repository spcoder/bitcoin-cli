package main

import "github.com/parnurzeal/gorequest"

type HttpRequester interface {
	Get() (gorequest.Response, string, []error)
}

type HttpManager struct {
}

func (HttpManager) Get() (gorequest.Response, string, []error) {
	return gorequest.New().Get("https://api.gdax.com/products/BTC-USD/ticker").End()
}
