package main

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spcoder/bitcoin/main_test"
)

func TestRunPrice(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	json := `{"trade_id":32865293,"price":"13960.87000000","size":"0.00019568","bid":"13962.76","ask":"13970.8","volume":"21704.54760458","time":"2018-01-10T13:09:52.893000Z"}`

	mockHttpRequester := main_test.NewMockHttpRequester(mockCtrl)
	mockHttpRequester.EXPECT().Get().Return(nil, json, nil)

	if price, err := runPrice(mockHttpRequester, false); err != nil {
		t.Fatal(err)
	} else {
		if price != "$13,960.87 USD\n" {
			t.Fatal("Invalid price")
		}
	}
}
