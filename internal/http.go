package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type HTTPError struct {
	Body   []byte
	Status int
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("request failed with status: %d", e.Status)
}

type APIClient struct {
	L *log.Logger
	C *http.Client
}

func (c *APIClient) NewRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Request, error) {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func (c *APIClient) Do(req *http.Request, v interface{}) error {
	now := time.Now()
	resp, err := c.C.Do(req)
	if err != nil {
		return err
	}
	c.L.Printf("[%d] %s %s %v", resp.StatusCode, req.Method, req.URL, time.Since(now))
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return HTTPError{
			Body:   body,
			Status: resp.StatusCode,
		}
	}
	if v == nil {
		return nil
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}
	return nil
}
