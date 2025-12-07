package request

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequest_Logger(t *testing.T) {
	t.Parallel()

	logWriter := bytes.Buffer{}

	testLogger := logger{
		zerolog: h.zerolog.Output(&logWriter),
	}

	h = &testLogger

	reqId := "test"

	mockHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "http://localhost", nil)
	req.Header.Add(requestId, reqId)

	rr := httptest.NewRecorder()

	requestIdHandler := RequestId(Logger(mockHandler))

	requestIdHandler.ServeHTTP(rr, req)

	assert.Equal(
		t,
		http.StatusOK,
		rr.Code,
		"expected status OK 200, got = %v", rr.Code,
	)

	parsedLog := map[string]interface{}{}
	_ = json.Unmarshal(logWriter.Bytes(), &parsedLog)

	assert.Equal(
		t,
		reqId,
		parsedLog["request_id"],
		"expected = %v, got = %v", reqId, parsedLog["request_id"],
	)

	assert.Equal(
		t,
		"request",
		parsedLog["level"],
		"expected = request, got = %v", parsedLog["level"],
	)

	assert.Equal(
		t,
		"GET",
		parsedLog["method"],
		"expected = GET, got = %v", parsedLog["method"],
	)

	assert.Equal(
		t,
		float64(200),
		parsedLog["status"],
		"expected = 200, got = %v", parsedLog["status"],
	)

	assert.Equal(
		t,
		"HTTP/1.1",
		parsedLog["proto"],
		"expected = HTTP/1.1, got = %v", parsedLog["proto"],
	)

	assert.Equal(
		t,
		"http://localhost",
		parsedLog["uri"],
		"expected = http://localhost, got = %v", parsedLog["uri"],
	)
}
