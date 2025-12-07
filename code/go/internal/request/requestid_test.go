package request

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRequest_RequestId_SetInHeader(t *testing.T) {
	t.Parallel()

	reqId := "test"

	mockHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ret := GetReqIdCtx(req.Context())
		assert.Equal(
			t,
			reqId,
			ret,
			"expected = %v, got = %v", reqId, ret,
		)

		rw.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "http://localhost", nil)
	req.Header.Add(requestId, reqId)

	rr := httptest.NewRecorder()

	requestIdHandler := RequestId(mockHandler)

	requestIdHandler.ServeHTTP(rr, req)

	assert.Equal(
		t,
		http.StatusOK,
		rr.Code,
		"expected status OK 200, got = %v", rr.Code,
	)
}

func TestRequest_RequestId_NotSet(t *testing.T) {
	t.Parallel()

	mockHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ret := GetReqIdCtx(req.Context())
		assert.NotNil(
			t,
			ret,
			"request id should not be empty",
		)

		parsedId, err := uuid.Parse(ret)
		assert.NotEqual(
			t,
			uuid.Nil,
			parsedId,
			"Parsed Id should no be nil uuid",
		)

		assert.NoError(
			t,
			err,
		)

		rw.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "http://localhost", nil)

	rr := httptest.NewRecorder()

	requestIdHandler := RequestId(mockHandler)

	requestIdHandler.ServeHTTP(rr, req)

	assert.Equal(
		t,
		http.StatusOK,
		rr.Code,
		"expected status OK 200, got = %v", rr.Code,
	)
}

func TestRequest_GetReqId(t *testing.T) {
	t.Parallel()

	// Case 1: No request headers.
	testRequest := httptest.NewRequest("GET", "http://localhost", nil)

	ret := GetReqId(testRequest)
	assert.Empty(
		t,
		ret,
		"Case 1: expected empty value",
	)

	// Case 2: Header set.
	testHeader := "test"
	testRequest = httptest.NewRequest("GET", "http://localhost", nil)
	testRequest.Header.Add(requestId, testHeader)

	ret = GetReqId(testRequest)
	assert.Equal(
		t,
		testHeader,
		ret,
		"Case 2: expected = %v, got = %v", testHeader, ret,
	)

	// Case 3: nil request.
	ret = GetReqId(nil)
	assert.Empty(
		t,
		ret,
		"Case 3: expected empty value",
	)
}
