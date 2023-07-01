package tests

import (
	"flag"
	"fmt"
	"os"
	"receipt-processor-challenge/config"
	. "receipt-processor-challenge/internal"
	"receipt-processor-challenge/internal/cache"
	"receipt-processor-challenge/internal/entity"
	"receipt-processor-challenge/internal/handler"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var (
	configPath   string
	log          *logrus.Logger
	ReceiptCache = make(map[string]*entity.Receipt, 0)
)

type ReceiptId struct {
	Id string `json:"id"`
}

func TestMain(m *testing.M) {
	flag.StringVar(&configPath, "config", "../config.yaml", "config file path")
	flag.Parse()
	os.Exit(m.Run())
}

func run(cfg config.Config) *fasthttp.Client {

	log = logrus.New()
	log.SetLevel(logrus.DebugLevel)
	h := handler.NewHandler(log)
	router := NewProcessorService(h)
	server := fasthttp.Server{
		ReadBufferSize: 10240,
		Handler:        router.Handler,
		Logger:         log,
	}

	go func() {
		if err := server.ListenAndServe(":" + cfg.App.Port); err != nil {
			log.Fatalf("Ошибка при запуске сервера: %v", err)
		}
	}()
	client := &fasthttp.Client{}
	return client
}
func TestProcessReceipt(t *testing.T) {
	cfg, err := config.Setup(configPath)
	if err != nil {
		panic(err)
	}
	client := run(cfg)

	type expectation struct {
		statusCode int
	}

	tests := map[string]struct {
		in   string
		want expectation
	}{
		"Sucsess. Full JSON receipt struct": {
			in: `{
				"retailer": "M&M Corner Market",
				"purchaseDate": "2022-03-20",
				"purchaseTime": "14:33",
				"items": [
				  {
					"shortDescription": "Gatorade",
					"price": "2.25"
				  },{
					"shortDescription": "Gatorade",
					"price": "2.25"
				  },{
					"shortDescription": "Gatorade",
					"price": "2.25"
				  },{
					"shortDescription": "Gatorade",
					"price": "2.25"
				  }
				],
				"total": "9.00"
			  }`,
			want: expectation{
				statusCode: 200,
			},
		},
		"Fail. Empty request body": {
			in: `{
			  }`,
			want: expectation{
				statusCode: 400,
			},
		},
	}

	for caseName, testCase := range tests {
		t.Run(caseName, func(t *testing.T) {
			response, request := ProcessReciept(testCase.in)
			defer fasthttp.ReleaseRequest(request)
			defer fasthttp.ReleaseResponse(response)

			request.Header.SetHost(cfg.App.Host + ":" + cfg.App.Port)
			err = client.Do(request, response)
			if err != nil {
				t.Fatalf("response error : %v", err)
			}

			if response.StatusCode() != testCase.want.statusCode {
				t.Errorf("Out -> \nWant: %q\nGot : %q , response body %q", testCase.want.statusCode, response.StatusCode(), string(response.Body()))
			}
		})
	}
}

func ProcessReciept(jsonBody string) (*fasthttp.Response, *fasthttp.Request) {
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	request.Header.SetMethod("POST")
	request.SetRequestURI("/receipts/process")
	request.SetBody([]byte(jsonBody))
	return response, request
}

func GetPoints(id string) (*fasthttp.Response, *fasthttp.Request) {
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	request.Header.SetMethod("GET")
	url := fmt.Sprint("/receipts/", id, "/points")
	request.SetRequestURI(url)
	return response, request
}

func TestProcessReceiptAPI(t *testing.T) {

	cfg, err := config.Setup(configPath)
	if err != nil {
		panic(err)
	}
	client := run(cfg)

	type expectation struct {
		PostMethod struct {
			statusCode int
		}
		GetMethod struct {
			statusCode int
			points     string
		}
	}
	type testIn struct {
		receipt string
		id      string
	}

	tests := map[string]struct {
		in   testIn
		want expectation
	}{
		"SUCCSESS CASE_Full JSON receipt struct": {
			in: testIn{
				receipt: `{
					"retailer": "M&M Corner Market",
					"purchaseDate": "2022-03-20",
					"purchaseTime": "14:33",
					"items": [
					  {
						"shortDescription": "Gatorade",
						"price": "2.25"
					  },{
						"shortDescription": "Gatorade",
						"price": "2.25"
					  },{
						"shortDescription": "Gatorade",
						"price": "2.25"
					  },{
						"shortDescription": "Gatorade",
						"price": "2.25"
					  }
					],
					"total": "9.00"
				  }`,
			},
			want: expectation{
				PostMethod: struct{ statusCode int }{
					statusCode: 200,
				},
				GetMethod: struct {
					statusCode int
					points     string
				}{
					statusCode: 200,
					points:     "109",
				},
			},
		},
		"FAIL CASE_Empty request body": {
			in: testIn{receipt: `{
			  }`,
			},
			want: expectation{
				PostMethod: struct{ statusCode int }{
					statusCode: 400,
				},
			},
		},
		"FAIL CASE_No receipt with this ID": {
			in: testIn{id: uuid.New().String()},
			want: expectation{
				GetMethod: struct {
					statusCode int
					points     string
				}{
					statusCode: 400,
					points:     "",
				},
			},
		},
	}

	for caseName, testCase := range tests {
		t.Run(caseName, func(t *testing.T) {

			cache := cache.ReceiptCache

			//POST
			if caseName != "FAIL CASE_No receipt with this ID" {
				response, request := ProcessReciept(testCase.in.receipt)
				defer fasthttp.ReleaseRequest(request)
				defer fasthttp.ReleaseResponse(response)
				request.Header.SetHost(cfg.App.Host + ":" + cfg.App.Port)
				err = client.Do(request, response)
				if err != nil {
					t.Fatalf("response error : %v", err)
				}

				if response.StatusCode() != testCase.want.PostMethod.statusCode {
					t.Errorf("Out -> \nWant: %q\nGot : %q , response body %q", testCase.want.PostMethod.statusCode, response.StatusCode(), string(response.Body()))
				}
			}

			//GET
			if caseName != "FAIL CASE_Empty request body" {
				for id, _ := range cache {
					response, request := GetPoints(id)
					request.Header.SetHost(cfg.App.Host + ":" + cfg.App.Port)
					err = client.Do(request, response)
					if err != nil {
						t.Fatalf("response error : %v", err)
					}
					if response.StatusCode() != testCase.want.GetMethod.statusCode || !strings.Contains(string(response.Body()), testCase.want.GetMethod.points) {
						t.Errorf("Out -> \nWant: %d\nGot : %d , response body %q", testCase.want.GetMethod.statusCode, response.StatusCode(), string(response.Body()))
					}
					t.Log(string(response.Body()))
					delete(cache, id)
				}
			}
		})
	}
}
