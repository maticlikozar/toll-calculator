package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequest_ParseIp(t *testing.T) {
	t.Parallel()

	// Case 1: Empty value.
	testIp := ""
	ret, err := ParseIp(testIp)
	assert.Empty(
		t,
		ret,
		"Case 1: expected empty",
	)
	assert.NoError(
		t,
		err,
		"Case 1: expected no error",
	)

	// Case 2: IP without port.
	testIp = "127.0.0.1"
	ret, err = ParseIp(testIp)
	assert.Equal(
		t,
		ret,
		testIp,
		"Case 2: expected = %v, got = %v", testIp, ret,
	)
	assert.NoError(
		t,
		err,
		"Case 2: expected no error",
	)

	// Case 3: IP with port.
	ret, err = ParseIp(testIp + ":5687")
	assert.Equal(
		t,
		ret,
		testIp,
		"Case 3: expected = %v, got = %v", testIp, ret,
	)
	assert.NoError(
		t,
		err,
		"Case 3: expected no error",
	)

	// Case 4: Invalid IP.
	expectedError := "address 127.0.0.1::::5687: too many colons in address"
	ret, err = ParseIp(testIp + "::::5687")
	assert.Equal(
		t,
		ret,
		"",
		"Case 4: expected = %v, got = %v", testIp, ret,
	)
	assert.EqualError(
		t,
		err,
		expectedError,
		"Case 4: expected error = , got = %v", expectedError, err,
	)

	// Case 5: IPv6 without port.
	testIp = "::1"
	ret, err = ParseIp(testIp)
	assert.Equal(
		t,
		ret,
		testIp,
		"Case 2: expected = %v, got = %v", testIp, ret,
	)
	assert.NoError(
		t,
		err,
		"Case 2: expected no error",
	)

	// Case 6: IPv6 without port.
	testIp = "::1"
	ret, err = ParseIp("[::1]:5687")
	assert.Equal(
		t,
		ret,
		testIp,
		"Case 2: expected = %v, got = %v", testIp, ret,
	)
	assert.NoError(
		t,
		err,
		"Case 2: expected no error",
	)
}
