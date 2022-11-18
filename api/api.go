package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mkuznets.com/go/form3/models"
	"net/http"
	"net/url"
	"time"
)

// Api manages requests and responses from the Form3 API endpoints.
type Api struct {
	// BaseUrl is the parsed base URL of the Form3 API.
	BaseUrl *url.URL
	// HttpClient is an instance of http.Client used for API requests.
	HttpClient *http.Client
	// BackOffProvider returns a fresh BackOff instance that governs API endpoint retry policy.
	BackOffProvider func() BackOff
	// UuidProvider returns unique UUIDv4 identifiers used as ID of new Form3 API resources.
	UuidProvider func() string
	// OrganisationId is the organisation ID used in the Form3 API requests.
	OrganisationId string
}

// New create a new Form3 api with the given base URL.
func New(baseUrl string, organisationId string) (*Api, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	return &Api{
		BaseUrl:         u,
		HttpClient:      &http.Client{},
		BackOffProvider: DefaultBackOffProvider,
		UuidProvider: func() string {
			return uuid.New().String()
		},
		OrganisationId: organisationId,
	}, nil
}

func (s *Api) Do(ctx context.Context, call *Call) error {
	request, err := call.httpRequest(ctx, s.BaseUrl)
	if err != nil {
		return err
	}

	resp, err := s.withRetries(func() (*http.Response, error) {
		resp, errC := s.HttpClient.Do(request)
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

func (s *Api) withRetries(handler func() (*http.Response, error)) (*http.Response, error) {
	backOff := s.BackOffProvider()
	if backOff == nil {
		return nil, fmt.Errorf("BackOffProvider returned nil")
	}

	for {
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
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
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
	if resp == nil {
		return
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
}
