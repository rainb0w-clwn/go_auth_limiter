package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var httpBaseURL = "http://app:" + os.Getenv("HTTP_PORT")

func httpClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}

func doJSONRequest(
	t *testing.T,
	method string,
	url string,
	reqBody any,
) (*http.Response, []byte) {
	t.Helper()

	var body io.Reader
	if reqBody != nil {
		b, err := json.Marshal(reqBody)
		require.NoError(t, err)
		body = bytes.NewBuffer(b)
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, method, httpBaseURL+url, body)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient().Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, respBody
}

type IPNetRequest struct {
	IPNet string `json:"ipNet"`
}

type BucketResetRequestDTO struct {
	Login string `json:"login"`
	IP    string `json:"ip"`
}

type LimitCheckRequestDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	IP       string `json:"ip"`
}

type LimitCheckResponseDTO struct {
	Allowed bool `json:"allowed"`
}

func TestHTTP_WhiteListAdd_OK(t *testing.T) {
	resp, _ := doJSONRequest(
		t,
		http.MethodPost,
		"/whitelist",
		IPNetRequest{IPNet: "10.10.10.0/24"},
	)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHTTP_WhiteListAdd_InvalidArgument(t *testing.T) {
	resp, _ := doJSONRequest(
		t,
		http.MethodPost,
		"/whitelist",
		IPNetRequest{IPNet: ""},
	)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHTTP_WhiteListDelete_NotFound(t *testing.T) {
	resp, _ := doJSONRequest(
		t,
		http.MethodDelete,
		"/whitelist",
		IPNetRequest{IPNet: "192.168.100.0/24"},
	)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHTTP_WhiteListDelete_InvalidArgument(t *testing.T) {
	resp, _ := doJSONRequest(
		t,
		http.MethodDelete,
		"/whitelist",
		IPNetRequest{IPNet: ""},
	)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHTTP_BlackListAdd_OK(t *testing.T) {
	resp, _ := doJSONRequest(
		t,
		http.MethodPost,
		"/blacklist",
		IPNetRequest{IPNet: "172.16.0.0/16"},
	)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHTTP_BlackListAdd_InvalidArgument(t *testing.T) {
	resp, _ := doJSONRequest(
		t,
		http.MethodPost,
		"/blacklist",
		IPNetRequest{IPNet: "invalid"},
	)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHTTP_BlackListDelete_NotFound(t *testing.T) {
	resp, _ := doJSONRequest(
		t,
		http.MethodDelete,
		"/blacklist",
		IPNetRequest{IPNet: "8.8.8.0/24"},
	)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHTTP_BlackListDelete_InvalidArgument(t *testing.T) {
	resp, _ := doJSONRequest(
		t,
		http.MethodDelete,
		"/blacklist",
		IPNetRequest{IPNet: ""},
	)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHTTP_BucketReset_OK(t *testing.T) {
	resp, body := doJSONRequest(
		t,
		http.MethodPost,
		"/check",
		LimitCheckRequestDTO{
			Login:    "user1",
			Password: "secret",
			IP:       "127.0.0.1",
		},
	)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var parsed LimitCheckResponseDTO
	require.NoError(t, json.Unmarshal(body, &parsed))

	resp, _ = doJSONRequest(
		t,
		http.MethodPost,
		"/reset",
		BucketResetRequestDTO{
			Login: "user1",
			IP:    "127.0.0.1",
		},
	)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHTTP_BucketReset_InvalidArgument(t *testing.T) {
	resp, _ := doJSONRequest(
		t,
		http.MethodPost,
		"/reset",
		BucketResetRequestDTO{
			Login: "",
			IP:    "",
		},
	)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHTTP_LimitCheck_OK(t *testing.T) {
	resp, body := doJSONRequest(
		t,
		http.MethodPost,
		"/check",
		LimitCheckRequestDTO{
			Login:    "user1",
			Password: "secret",
			IP:       "127.0.0.1",
		},
	)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var parsed LimitCheckResponseDTO
	require.NoError(t, json.Unmarshal(body, &parsed))
}

func TestHTTP_LimitCheck_InvalidArgument(t *testing.T) {
	resp, _ := doJSONRequest(
		t,
		http.MethodPost,
		"/check",
		LimitCheckRequestDTO{
			Login:    "",
			Password: "",
			IP:       "",
		},
	)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
