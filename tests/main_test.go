package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/satori/go.uuid"
)

const (
	paymentURL = "http://localhost:8085/v1/payments"
	balanceURL = "http://localhost:8085/v1/payments/balance"
)

type payload struct {
	State         string `json:"state"`
	Amount        string `json:"amount"`
	TransactionID string `json:"transactionId"`
}

type response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func TestRequests(t *testing.T) {
	client := http.Client{Timeout: time.Duration(10) * time.Second}

	testSuite := []struct {
		testName       string
		state          string
		amount         string
		transactionID  string
		contentType    string
		sourceType     string
		expectedStatus int
		expectedResult string
	}{
		{
			testName:       "Test incorrect Content-Type",
			state:          "win",
			amount:         "10",
			transactionID:  uuid.NewV4().String(),
			contentType:    "plain/text",
			sourceType:     "payment",
			expectedStatus: 400,
			expectedResult: "{\"status\":false,\"error\":\"incorrect Content-Type, required application/json, got plain/text\"}",
		},
		{
			testName:       "Test incorrect Source-Type",
			state:          "win",
			amount:         "10",
			transactionID:  uuid.NewV4().String(),
			contentType:    "application/json",
			sourceType:     "client",
			expectedStatus: 400,
			expectedResult: "{\"status\":false,\"error\":\"wrong header Source-Type\"}",
		},
		{
			testName:       "Test incorrect amount",
			state:          "win",
			amount:         "aaaa",
			transactionID:  uuid.NewV4().String(),
			contentType:    "application/json",
			sourceType:     "payment",
			expectedStatus: 400,
			expectedResult: "{\"status\":false,\"error\":\"incorrect amount value\"}",
		},
		{
			testName:       "Test negative amount",
			state:          "win",
			amount:         "-1000",
			transactionID:  uuid.NewV4().String(),
			contentType:    "application/json",
			sourceType:     "payment",
			expectedStatus: 400,
			expectedResult: "{\"status\":false,\"error\":\"incorrect amount value\"}",
		},
		{
			testName:       "Test incorrect state",
			state:          "winwin",
			amount:         "10",
			transactionID:  uuid.NewV4().String(),
			contentType:    "application/json",
			sourceType:     "payment",
			expectedStatus: 400,
			expectedResult: "{\"status\":false,\"error\":\"State:\\n stateValidator: state field did not pass validation\"}",
		},
		{
			testName:       "Test incorrect payload",
			state:          "",
			amount:         "",
			transactionID:  "",
			contentType:    "application/json",
			sourceType:     "payment",
			expectedStatus: 400,
			expectedResult: "",
		},
		{
			testName:       "Test Success transaction",
			state:          "win",
			amount:         "10",
			transactionID:  uuid.NewV4().String(),
			contentType:    "application/json",
			sourceType:     "payment",
			expectedStatus: 200,
			expectedResult: "{\"status\":true,\"message\":\"payment was successfully proceed\"}",
		},
		{
			testName:       "Test Transaction ID idempotence",
			state:          "win",
			amount:         "10",
			transactionID:  "",
			contentType:    "application/json",
			sourceType:     "payment",
			expectedStatus: 400,
			expectedResult: "{\"status\":false,\"error\":\"failed to proceed payment: transaction_id already processed\"}",
		},
		{
			testName:       "Test payment greater than balance",
			state:          "lost",
			amount:         "10000000000",
			transactionID:  uuid.NewV4().String(),
			contentType:    "application/json",
			sourceType:     "payment",
			expectedStatus: 400,
			expectedResult: "{\"status\":false,\"error\":\"failed to proceed payment: insufficient funds\"}",
		},
	}

	for k, ts := range testSuite {
		t.Run(ts.testName, func(t *testing.T) {
			payload := payload{
				State:         ts.state,
				Amount:        ts.amount,
				TransactionID: ts.transactionID,
			}
			if ts.transactionID == "" {
				payload.TransactionID = testSuite[k-1].transactionID
			}
			tpBytes, err := json.Marshal(payload)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest("POST", paymentURL, bytes.NewBuffer(tpBytes))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", ts.contentType)
			req.Header.Set("Source-Type", ts.sourceType)

			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				errB := resp.Body.Close()
				if errB != nil {
					t.Error(errB)
				}
			}()

			if ts.expectedStatus != resp.StatusCode {
				t.Errorf("wrong status code: expected: %d, actual: %d", ts.expectedStatus, resp.StatusCode)
			}

			if ts.expectedResult != "" && strings.TrimSpace(ts.expectedResult) != strings.TrimSpace(string(body)) {
				t.Errorf("wrong body: expected: %s, actual: %s", ts.expectedResult, string(body))
			}
		})
	}
}

func TestConcurrentRequests(t *testing.T) {
	client := http.Client{Timeout: time.Duration(10) * time.Second}

	testSuite := []struct {
		testName       string
		state          string
		amount         string
		transactionID  string
		contentType    string
		sourceType     string
		expectedStatus int
		expectedResult string
	}{
		{
			testName:       "Request 1",
			state:          "win",
			amount:         "100",
			contentType:    "application/json",
			sourceType:     "payment",
			expectedStatus: 200,
			expectedResult: "{\"status\":true,\"message\":\"payment was successfully proceed\"}",
		},
		{
			testName:       "Request 2",
			state:          "lost",
			amount:         "10",
			contentType:    "application/json",
			sourceType:     "payment",
			expectedStatus: 200,
			expectedResult: "{\"status\":true,\"message\":\"payment was successfully proceed\"}",
		},
		{
			testName:       "Request 3",
			state:          "win",
			amount:         "100",
			contentType:    "application/json",
			sourceType:     "payment",
			expectedStatus: 200,
			expectedResult: "{\"status\":true,\"message\":\"payment was successfully proceed\"}",
		},
		{
			testName:       "Request 4",
			state:          "lost",
			amount:         "10",
			contentType:    "application/json",
			sourceType:     "payment",
			expectedStatus: 200,
			expectedResult: "{\"status\":true,\"message\":\"payment was successfully proceed\"}",
		},
		{
			testName:       "Request 5",
			state:          "win",
			amount:         "100",
			transactionID:  "",
			contentType:    "application/json",
			sourceType:     "payment",
			expectedStatus: 200,
			expectedResult: "{\"status\":true,\"message\":\"payment was successfully proceed\"}",
		},
		{
			testName:       "Request 6",
			state:          "lost",
			amount:         "10",
			contentType:    "application/json",
			sourceType:     "payment",
			expectedStatus: 200,
			expectedResult: "{\"status\":true,\"message\":\"payment was successfully proceed\"}",
		},
		{
			testName:       "Request 7",
			state:          "lost",
			amount:         "0.1",
			contentType:    "application/json",
			sourceType:     "payment",
			expectedStatus: 200,
			expectedResult: "{\"status\":true,\"message\":\"payment was successfully proceed\"}",
		},
	}

	balanceBefore, err := getBalance(&client)
	if err != nil {
		t.Fatal(err)
	}

	for _, ts := range testSuite {
		t.Run(ts.testName, func(t *testing.T) {
			u := uuid.NewV4()
			payload := payload{
				State:         ts.state,
				Amount:        ts.amount,
				TransactionID: u.String(),
			}
			tpBytes, err := json.Marshal(payload)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest("POST", paymentURL, bytes.NewBuffer(tpBytes))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", ts.contentType)
			req.Header.Set("Source-Type", ts.sourceType)

			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				errB := resp.Body.Close()
				if errB != nil {
					t.Error(errB)
				}
			}()

			if ts.expectedStatus != resp.StatusCode {
				t.Errorf("wrong status code: expected: %d, actual: %d", ts.expectedStatus, resp.StatusCode)
			}

			if ts.expectedResult != "" && strings.TrimSpace(ts.expectedResult) != strings.TrimSpace(string(body)) {
				t.Errorf("wrong body: expected: %s, actual: %s", ts.expectedResult, string(body))
			}
		})
	}

	balanceAfter, err := getBalance(&client)
	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprintf("%.2f", balanceBefore+269.9) != fmt.Sprintf("%.2f", balanceAfter) {
		t.Errorf(
			"incorrect concurrent behavior: expected balance: %s, actual: %s",
			fmt.Sprintf("%.2f", balanceBefore+269.9), fmt.Sprintf("%.2f", balanceAfter),
		)
	}
}

func TestCustomConcurrentRequests(t *testing.T) {
	client := http.Client{Timeout: time.Duration(10) * time.Second}

	type testData struct {
		testName       string
		state          string
		amount         string
		transactionID  string
		contentType    string
		sourceType     string
		expectedStatus int
		expectedResult string
	}

	testSuite := []testData{
		{
			testName:       "Request 7",
			state:          "win",
			amount:         "100",
			contentType:    "application/json",
			sourceType:     "payment",
			expectedStatus: 200,
			expectedResult: "{\"status\":true,\"message\":\"payment was successfully proceed\"}",
		},
		{
			testName:       "Request 8",
			state:          "lost",
			amount:         "10",
			contentType:    "application/json",
			sourceType:     "payment",
			expectedStatus: 200,
			expectedResult: "{\"status\":true,\"message\":\"payment was successfully proceed\"}",
		},
		{
			testName:       "Request 9",
			state:          "win",
			amount:         "100",
			contentType:    "application/json",
			sourceType:     "payment",
			expectedStatus: 200,
			expectedResult: "{\"status\":true,\"message\":\"payment was successfully proceed\"}",
		},
		{
			testName:       "Request 10",
			state:          "lost",
			amount:         "10",
			contentType:    "application/json",
			sourceType:     "payment",
			expectedStatus: 200,
			expectedResult: "{\"status\":true,\"message\":\"payment was successfully proceed\"}",
		},
		{
			testName:       "Request 11",
			state:          "win",
			amount:         "100",
			transactionID:  "",
			contentType:    "application/json",
			sourceType:     "payment",
			expectedStatus: 200,
			expectedResult: "{\"status\":true,\"message\":\"payment was successfully proceed\"}",
		},
		{
			testName:       "Request 12",
			state:          "lost",
			amount:         "10",
			contentType:    "application/json",
			sourceType:     "payment",
			expectedStatus: 200,
			expectedResult: "{\"status\":true,\"message\":\"payment was successfully proceed\"}",
		},
		{
			testName:       "Request 13",
			state:          "lost",
			amount:         "0.1",
			contentType:    "application/json",
			sourceType:     "payment",
			expectedStatus: 200,
			expectedResult: "{\"status\":true,\"message\":\"payment was successfully proceed\"}",
		},
	}

	balanceBefore, err := getBalance(&client)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(testSuite))
	for _, ts := range testSuite {
		go func(t *testing.T, td testData) {
			defer wg.Done()
			u := uuid.NewV4()
			payload := payload{
				State:         td.state,
				Amount:        td.amount,
				TransactionID: u.String(),
			}
			tpBytes, err := json.Marshal(payload)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest("POST", paymentURL, bytes.NewBuffer(tpBytes))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", td.contentType)
			req.Header.Set("Source-Type", td.sourceType)

			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				errB := resp.Body.Close()
				if errB != nil {
					t.Error(errB)
				}
			}()

			if td.expectedStatus != resp.StatusCode {
				t.Errorf("wrong status code: expected: %d, actual: %d", td.expectedStatus, resp.StatusCode)
			}

			if td.expectedResult != "" && strings.TrimSpace(td.expectedResult) != strings.TrimSpace(string(body)) {
				t.Errorf("wrong body: expected: %s, actual: %s", td.expectedResult, string(body))
			}
		}(t, ts)
	}

	wg.Wait()

	balanceAfter, err := getBalance(&client)
	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprintf("%.2f", balanceBefore+269.9) != fmt.Sprintf("%.2f", balanceAfter) {
		t.Errorf(
			"incorrect concurrent behavior: expected balance: %s, actual: %s",
			fmt.Sprintf("%.2f", balanceBefore+269.9), fmt.Sprintf("%.2f", balanceAfter),
		)
	}
}

func getBalance(client *http.Client) (float64, error) {
	req, err := http.NewRequest("GET", balanceURL, nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var r response
	if err := json.Unmarshal(body, &r); err != nil {
		return 0, err
	}

	return strconv.ParseFloat(r.Message, 64)
}
