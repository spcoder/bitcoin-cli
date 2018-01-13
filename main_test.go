package main

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/spcoder/bitcoin-cli/main_test"
)

func Test_main(t *testing.T) {
	badCommandUsage, _ := ioutil.ReadFile(path.Join("testdata", "bad_command.txt"))
	tests := []struct {
		name    string
		osArgs  []string
		outExpr *regexp.Regexp
		wantErr string
	}{
		{"RunPrice", []string{"bitcoin", "price"}, regexp.MustCompile("\\$.*USD"), ""},
		{"BadCommand", []string{"bitcoin", "bad"}, nil, string(badCommandUsage)},
	}
	for _, tt := range tests {
		os.Args = tt.osArgs
		t.Run(tt.name, func(t *testing.T) {
			stderr := os.Stderr
			stdout := os.Stdout
			errReader, errWriter, _ := os.Pipe()
			outReader, outWriter, _ := os.Pipe()
			os.Stderr = errWriter
			os.Stdout = outWriter
			main()
			errWriter.Close()
			outWriter.Close()
			errData, _ := ioutil.ReadAll(errReader)
			outData, _ := ioutil.ReadAll(outReader)
			os.Stderr = stderr
			os.Stdout = stdout
			if tt.outExpr != nil {
				if !tt.outExpr.MatchString(string(outData)) {
					t.Errorf("main() error = %v, wanted %v", string(outData), "/"+tt.outExpr.String()+"/")
				}
			}
			if strings.Trim(string(errData), "\n") != strings.Trim(tt.wantErr, "\n") {
				t.Errorf("main() error = %v, wanted %v", string(errData), tt.wantErr)
			}
		})
	}
}

func Test_runPrice(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type args struct {
		httpManager HTTPRequester
		rawFormat   bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		errStr  string
	}{
		{"Success", args{httpRequester(mockCtrl), false}, "$13,960.87 USD\n", false, ""},
		{"SuccessRaw", args{httpRequester(mockCtrl), true}, "13960.87000000\n", false, ""},
		{"ErrorWhenGettingTickerData", args{errorHTTPRequester(mockCtrl), false}, "", true, "unable to get bitcoin price in USD: [some error]"},
		{"ErrorWhenParsingTickerData", args{badJSONHTTPRequester(mockCtrl), false}, "", true, "unable to parse ticker data: invalid character 't' looking for beginning of object key string"},
		{"ErrorWhenParsingPrice", args{badPriceHTTPRequester(mockCtrl), false}, "", true, "unable to parse the ticker price: syntax error scanning number"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runPrice(tt.args.httpManager, tt.args.rawFormat)
			if (err != nil) != tt.wantErr {
				t.Errorf("runPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errStr != err.Error() {
				t.Errorf("runPrice() error `%v`, want `%v`", err.Error(), tt.errStr)
			}
			if got != tt.want {
				t.Errorf("runPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func httpRequester(mockCtrl *gomock.Controller) HTTPRequester {
	json := `{"trade_id":32865293,"price":"13960.87000000","size":"0.00019568","bid":"13962.76","ask":"13970.8","volume":"21704.54760458","time":"2018-01-10T13:09:52.893000Z"}`
	mockHTTPRequester := main_test.NewMockHttpRequester(mockCtrl)
	mockHTTPRequester.EXPECT().Get().Return(nil, json, nil)
	return mockHTTPRequester
}

func errorHTTPRequester(mockCtrl *gomock.Controller) HTTPRequester {
	mockHTTPRequester := main_test.NewMockHttpRequester(mockCtrl)
	mockHTTPRequester.EXPECT().Get().Return(nil, "", []error{errors.New("some error")})
	return mockHTTPRequester
}

func badJSONHTTPRequester(mockCtrl *gomock.Controller) HTTPRequester {
	json := `{trade_id:32865293,"price":"13960.87000000","size":"0.00019568","bid":"13962.76","ask":"13970.8","volume":"21704.54760458","time":"2018-01-10T13:09:52.893000Z"}`
	mockHTTPRequester := main_test.NewMockHttpRequester(mockCtrl)
	mockHTTPRequester.EXPECT().Get().Return(nil, json, nil)
	return mockHTTPRequester
}

func badPriceHTTPRequester(mockCtrl *gomock.Controller) HTTPRequester {
	json := `{"trade_id":32865293,"price":"abc.87000000","size":"0.00019568","bid":"13962.76","ask":"13970.8","volume":"21704.54760458","time":"2018-01-10T13:09:52.893000Z"}`
	mockHTTPRequester := main_test.NewMockHttpRequester(mockCtrl)
	mockHTTPRequester.EXPECT().Get().Return(nil, json, nil)
	return mockHTTPRequester
}
