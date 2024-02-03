package utils

import (
	"io"
	"net/http"
	"strings"
	"time"
)

func MkHttpClient() *http.Client {
	client := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	return &client
}

type TestRoundTripper struct {
	Responses *Stack[string]
}

func NewTestHttpClient(responses *Stack[string]) *http.Client {
	return &http.Client{
		Transport: &TestRoundTripper{Responses: responses},
	}
}

func (r *TestRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	data := r.Responses.Pop()
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(*data)),
	}, nil
}

func MkJSONRequest(method string, url string) (*http.Request, error) {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "codemetagenerator")
	request.Header.Set("Accept", "application/json")

	return request, nil
}

func DoRequest(client *http.Client, request *http.Request) (*[]byte, error) {
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return &bytes, nil
}
