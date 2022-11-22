package form3

//go:generate moq -out api_mock_test.go . Api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"mkuznets.com/go/form3/models"
)

// Api a is low-level Form3 API client to access arbitrary endpoints. Most users should use specialised clients (such as AccountsClient) instead.
type Api interface {
	// Do performs a Form3 API call according to the values provided in Call.
	Do(ctx context.Context, call *Call) error
}

type api struct {
	c *Client
}

func (a *api) Do(ctx context.Context, call *Call) error {
	baseUrl, err := url.Parse(a.c.baseUrl)
	if err != nil {
		return err
	}

	request, err := call.httpRequest(ctx, baseUrl)
	if err != nil {
		return err
	}

	resp, err := a.withRetries(func() (*http.Response, error) {
		resp, errC := a.c.httpClient.Do(request)
		if errC != nil {
			return nil, errC
		}
		return resp, errorFromResponse(resp)
	})
	if err != nil {
		return err
	}

	defer drainBody(resp)

	if call.Response != nil {
		body := models.Body{Data: call.Response}
		if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
			return err
		}
	}

	return nil
}

func (a *api) withRetries(handler func() (*http.Response, error)) (*http.Response, error) {
	for backOff := a.c.backOffProvider(); ; {
		resp, err := handler()
		if err != nil && shouldRetry(err) {
			delay := backOff.NextBackOff()
			if delay >= 0 {
				time.Sleep(delay)
				continue
			} else {
				return nil, err
			}
		}
		return resp, err
	}
}

func errorFromResponse(resp *http.Response) error {
	if resp.StatusCode/100 == 2 {
		return nil
	}
	defer drainBody(resp)

	apiErr := Error{
		StatusCode: resp.StatusCode,
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	apiErr.RawBody = body
	_ = json.Unmarshal(body, &apiErr)

	return apiErr
}

func shouldRetry(err error) bool {
	switch e := err.(type) {
	case Error:
		return e.Type() == ErrorServerError || e.Type() == ErrorTooManyRequests
	default:
		return true
	}
}

func drainBody(resp *http.Response) {
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	_, _ = io.Copy(io.Discard, resp.Body)
}
