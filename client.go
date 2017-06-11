package widevine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// HTTPClient defines an HTTP client.
type HTTPClient struct {
	*http.Client
}

func (c *HTTPClient) get(url string, i interface{}) error {

	fmt.Println("querying..." + url)
	rsp, e := c.Get(url)
	if e != nil {
		return e
	}

	defer rsp.Body.Close()

	b, e := ioutil.ReadAll(rsp.Body)
	if e != nil {
		return e
	}
	if rsp.Status[0] != '2' {
		return fmt.Errorf("expected status 2xx, got %s: %s", rsp.Status, string(b))
	}

	return json.Unmarshal(b, &i)
}

func (c *HTTPClient) post(url string, i interface{}, body interface{}) error {

	fmt.Println("querying..." + url)
	payload, _ := json.Marshal(body)
	req, e := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if e != nil {
		return e
	}

	req.Header.Add("content-type", "application/json")

	rsp, e := http.DefaultClient.Do(req)
	if e != nil {
		return e
	}

	defer rsp.Body.Close()

	b, e := ioutil.ReadAll(rsp.Body)
	if e != nil {
		return e
	}
	if rsp.Status[0] != '2' {
		return fmt.Errorf("expected status 2xx, got %s: %s", rsp.Status, string(b))
	}
	return json.Unmarshal(b, &i)
}

// NewClient creates an HTTPClient instance.
func NewClient() (*HTTPClient, error) {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	return &HTTPClient{Client: &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}}, nil
}
