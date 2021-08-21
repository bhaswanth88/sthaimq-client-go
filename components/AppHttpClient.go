package components

import (
	"github.com/gojek/heimdall"
	"github.com/gojek/heimdall/v7/httpclient"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type AppHttpClient struct {
	client *httpclient.Client
}

var singleton *AppHttpClient
var instanceOnce sync.Once

func GetAppHttpClient() *AppHttpClient {
	instanceOnce.Do(func() {
		backoffInterval := 2 * time.Millisecond
		// Define a maximum jitter interval. It must be more than 1*time.Millisecond
		maximumJitterInterval := 5 * time.Millisecond

		backoff := heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval)

		// Create a new retry mechanism with the backoff
		retrier := heimdall.NewRetrier(backoff)

		timeout := 1000 * time.Millisecond
		// Create a new client, sets the retry mechanism, and the number of times you would like to retry
		appHttpClient := httpclient.NewClient(
			httpclient.WithHTTPTimeout(timeout),
			httpclient.WithRetrier(retrier),
			httpclient.WithRetryCount(3),
		)

		singleton = &AppHttpClient{
			client: appHttpClient,
		}
	})
	return singleton
}

func (a *AppHttpClient) GetRequest(url url.URL, headers map[string][]string) (string, error) {
	return a.doRequest(url, "", headers, "GET")
}
func (a *AppHttpClient) PutRequest(url url.URL, headers map[string][]string, body string) (string, error) {
	return a.doRequest(url, body, headers, "PUT")
}
func (a *AppHttpClient) PostRequest(url url.URL, headers map[string][]string, body string) (string, error) {
	return a.doRequest(url, body, headers, "POST")
}

func (a *AppHttpClient) DeleteRequest(url url.URL, headers map[string][]string) (string, error) {
	return a.doRequest(url, "", headers, "DELETE")
}
func (a *AppHttpClient) doRequest(url url.URL, body string, headers map[string][]string, method string) (string, error) {

	reqBody := ioutil.NopCloser(strings.NewReader(body))

	req := &http.Request{
		Method: method,
		URL:    &url,
		Header: headers,
		Body:   reqBody,
	}
	resp, err := a.client.Do(req)
	// check for response error
	if err != nil {
		log.Fatal("Error:", err)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return string(data), err
}
