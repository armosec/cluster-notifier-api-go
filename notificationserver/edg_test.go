package notificationserver

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPushNotificationServer(t *testing.T) {
	type args struct {
		edgeURL    string
		targetMap  map[string]string
		message    interface{}
		jsonFormat bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Successful POST request to server",
			args: args{
				edgeURL: "http://example.com",
				targetMap: map[string]string{
					TargetCustomer:  "123",
					TargetCluster:   "my-cluster",
					TargetComponent: TargetComponentPostureValue,
				},
				message: struct {
					Text string `json:"text"`
				}{Text: "Hello, world!"},
				jsonFormat: true,
			},
			wantErr: false,
		},
		{
			name: "Failed POST request to server",
			args: args{
				edgeURL:    "http://bad-url",
				targetMap:  nil,
				message:    nil,
				jsonFormat: true,
			},
			wantErr: true,
		},
		{
			name: "Panic when unable to send request after three tries",
			args: args{
				edgeURL: "http://example.com",
				targetMap: map[string]string{
					TargetCustomer:  "123",
					TargetCluster:   "my-cluster",
					TargetComponent: TargetComponentPostureValue,
				},
				message: struct {
					Text string `json:"text"`
				}{Text: "Hello, world!"},
				jsonFormat: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := PushNotificationServer(tt.args.edgeURL, tt.args.targetMap, tt.args.message, tt.args.jsonFormat)
			if (err != nil) != tt.wantErr {
				t.Errorf("PushNotificationServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_sendCommandToEdge(t *testing.T) {
	// Start mock server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Assert the request method is POST
		if req.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", req.Method)
		}

		// Assert the request body matches the test message
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			// Return a test response
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(rw, "test response")
		}
		if string(body) != "testMessage" {
			// Return a test response
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(rw, "test response")
		}

		// Return a test response
		rw.WriteHeader(http.StatusOK)
		fmt.Fprint(rw, "test response")
	}))
	defer server.Close()

	// Create mock client
	client := &http.Client{
		// Transport: server.Transport,
	}

	// Define test cases
	testCases := []struct {
		name    string
		message []byte
		wantErr bool
	}{
		{
			name:    "valid request",
			message: []byte("testMessage"),
			wantErr: false,
		},
		{
			name:    "invalid request",
			message: []byte("invalid message"),
			wantErr: true,
		},
	}

	// Loop over test cases and run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := sendCommandToEdge(client, server.URL, tc.message)
			if (err != nil) != tc.wantErr {
				t.Errorf("sendCommandToEdge() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func Test_setNotification(t *testing.T) {
	testCases := []struct {
		name       string
		targetMap  map[string]string
		message    interface{}
		jsonFormat bool
		expected   []byte
		expectErr  bool
	}{
		{
			name: "marshal to json",
			targetMap: map[string]string{
				TargetCustomer: "testCustomer",
			},
			message: map[string]string{
				"key": "value",
			},
			jsonFormat: true,
			expected:   []byte(`{"target":{"customerGUID":"testCustomer"},"sendSynchronicity":false,"notification":{"key":"value"}}`),
			expectErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := setNotification(tc.targetMap, tc.message, tc.jsonFormat)
			if tc.expectErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("expected no error but got %v", err)
			}
			if !bytes.Equal(actual, tc.expected) {
				t.Errorf("expected %s, but got %s", tc.expected, actual)
			}
		})
	}
}

// func Test_httpRespToString(t *testing.T) {
// 	testCases := []struct {
// 		name           string
// 		resp           *http.Response
// 		expectedResult string
// 		expectedErr    error
// 	}{
// 		{
// 			name: "successful response",
// 			resp: &http.Response{
// 				StatusCode: 200,
// 				Body:       ioutil.NopCloser(bytes.NewBufferString("response body")),
// 			},
// 			expectedResult: "response body",
// 			expectedErr:    nil,
// 		},
// 		{
// 			name: "unsuccessful response",
// 			resp: &http.Response{
// 				StatusCode: 500,
// 				Body:       ioutil.NopCloser(bytes.NewBufferString("error message")),
// 			},
// 			expectedResult: "error message",
// 			expectedErr:    fmt.Errorf("response status: 500. content: error message"),
// 		},
// 		{
// 			name:           "empty response",
// 			resp:           nil,
// 			expectedResult: "",
// 			expectedErr:    fmt.Errorf("empty response"),
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			result, err := httpRespToString(tc.resp)
// 			if result != tc.expectedResult {
// 				t.Errorf("Unexpected result. Expected: %v, but got: %v", tc.expectedResult, result)
// 			}
// 			if !errors.Is(err, tc.expectedErr) {
// 				t.Errorf("Unexpected error. Expected: %v, but got: %v", tc.expectedErr, err)
// 			}
// 		})
// 	}
// }
