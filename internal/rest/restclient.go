package rest

import (
	"io"
	"net/http"
)

type RestClientIface interface {
	Get(url string) RestClientIface
	Do() (*RestResponse, error)
}

type RestClient struct {
	client  *http.Client
	request *http.Request
	err     error
}

type RestResponse struct {
	Body       []byte
	StatusCode int
}

func New(clientOverride *http.Client) RestClientIface {
	if clientOverride != nil {
		return &RestClient{client: clientOverride}
	}
	return &RestClient{client: http.DefaultClient}
}

func (rc *RestClient) Get(url string) RestClientIface {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		rc.err = err
	}
	rc.request = req
	return rc
}

func (rc *RestClient) Post(url string, body io.Reader) *RestClient {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		rc.err = err
	}
	rc.request = req
	return rc
}

func (rc *RestClient) Put(url string, body io.Reader) *RestClient {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		rc.err = err
	}
	rc.request = req
	return rc
}

func (rc *RestClient) Delete(url string, body io.Reader) *RestClient {
	req, err := http.NewRequest(http.MethodDelete, url, body)
	if err != nil {
		rc.err = err
	}
	rc.request = req
	return rc
}

func (rc *RestClient) Headers(headers map[string]string) *RestClient {
	for key, value := range headers {
		rc.request.Header.Add(key, value)
	}
	return rc
}

func (rc *RestClient) Do() (*RestResponse, error) {
	if rc.err != nil {
		return nil, rc.err
	}
	res, err := rc.client.Do(rc.request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return &RestResponse{Body: body, StatusCode: res.StatusCode}, nil
}
