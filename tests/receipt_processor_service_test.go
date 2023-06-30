package tests

import (
	"flag"
	"os"
	"receipt-processor-challenge/config"
	. "receipt-processor-challenge/internal"
	"receipt-processor-challenge/internal/handler"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var (
	configPath string
	log        *logrus.Logger
)

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
	// ctx := context.Background()

	// type expectation struct {
	// 	out entity.ReceiptResponse
	// 	err error
	// }

	// tests := map[string]struct {
	// 	in   string
	// 	want expectation
	// }{
	// 	"Sucsess. Full JSON receipt struct": {
	// 		in: `{
	// 			"retailer": "M&M Corner Market",
	// 			"purchaseDate": "2022-03-20",
	// 			"purchaseTime": "14:33",
	// 			"items": [
	// 			  {
	// 				"shortDescription": "Gatorade",
	// 				"price": "2.25"
	// 			  },{
	// 				"shortDescription": "Gatorade",
	// 				"price": "2.25"
	// 			  },{
	// 				"shortDescription": "Gatorade",
	// 				"price": "2.25"
	// 			  },{
	// 				"shortDescription": "Gatorade",
	// 				"price": "2.25"
	// 			  }
	// 			],
	// 			"total": "9.00"
	// 		  }`,
	// 		want: expectation{
	// 			out: entity.ReceiptResponse{Id: ""},
	// 			err: nil,
	// 		},
	// 	},
	// }

	// for caseName, testCase := range tests {
	// 	t.Run(caseName, func(t *testing.T) {
	// 		resp, err := client.GetEmployeeById(ctx, testCase.id)
	// 		if err != nil {
	// 			if testCase.want.err.Error() != err.Error() {
	// 				t.Errorf("Err -> \nWant: %q\nGot: %q\n", testCase.want.err, err)
	// 			}
	// 		} else {
	// 			if testCase.want.out.Id != resp.FirstName ||
	// 				testCase.want.out.LastName != resp.LastName ||
	// 				testCase.want.out.MiddleName != resp.MiddleName {
	// 				t.Errorf("Out -> \nWant: %q\nGot : %q", testCase.want.out, resp)
	// 			}
	// 		}
	// 	})
	// }

	testCases := []string{
		`{
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
		`{}`,
	}

	for _, c := range testCases {
		response, request := processReciept(c)
		defer fasthttp.ReleaseRequest(request)
		defer fasthttp.ReleaseResponse(response)

		request.Header.SetHost(cfg.App.Host + ":" + cfg.App.Port)
		err = client.Do(request, response)
		if err != nil {
			t.Fatalf("response error : %v", err)
		}

		// Checking response code
		if response.StatusCode() != fasthttp.StatusOK {
			t.Errorf("Expected status code %d, actual %d", fasthttp.StatusOK, response.StatusCode())
		}

		// Checking response
		actualResponse := string(response.Body())
		if strings.Contains(actualResponse, `{"id":"`) {
			t.Logf("Actual response '%s'", actualResponse)
		} else {
			t.Errorf("Actual response '%s'", actualResponse)
		}
	}

}

func processReciept(jsonBody string) (*fasthttp.Response, *fasthttp.Request) {
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	request.Header.SetMethod("POST")
	request.SetRequestURI("/receipts/process")
	request.SetBody([]byte(jsonBody))
	return response, request
}
func TestGetProcessReceipt(t *testing.T) {

}
