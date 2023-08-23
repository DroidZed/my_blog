package tests

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

// executeRequest, creates a new ResponseRecorder
// then executes the request by calling ServeHTTP in the router
// after which the handler writes the response to the response recorder
// which we can then inspect.
func executeRequest(req *http.Request, s *TestServer) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

// checkResponseCode is a simple utility to check the response code
// of the response
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestHelloWorld(t *testing.T) {
	// Create a New Server Struct
	s := CreateNewTestServer()
	// Mount Handlers
	s.MountTestHandlers()

	// Create a New Request
	req, _ := http.NewRequest("GET", "/", nil)

	// Execute Request
	response := executeRequest(req, s)

	// Check the response code
	checkResponseCode(t, http.StatusOK, response.Code)

	// We can use testify/require to assert values, as it is more convenient
	require.Equal(t, "Hello World!", response.Body.String())
}
