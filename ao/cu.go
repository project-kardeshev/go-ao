package ao

import (
	"fmt"

	resty "github.com/go-resty/resty/v2"
)

type CuClient struct {
	client *resty.Client
	CuUrl string
}

type CuAPIs interface {
	DryRun(input DryRunInput, processId string) (result *Result, err error)
	Result(id string, processId string) (result *Result, err error)
	State(processId string) (state []byte, err error)
}

func (cuClient *CuClient) DryRun(input DryRunInput, processId string) (result *Result, err error) {
	response, err := cuClient.client.R().
		SetBody(input).
		SetQueryParam("process-id", processId).
		Post("/dry-run")

	if err != nil {
		return nil, err
	}

	return response.Result().(*Result), err
}

func (cuClient *CuClient) Result(id string, processId string) (result *Result, err error) {
	var cuResult Result
	
	response, err := cuClient.client.R().
		SetQueryParam("process-id", processId).
		SetResult(&cuResult).  // Let Resty handle the unmarshaling
		Get("/result/" + id)

	if err != nil {
		return nil, err
	}

	// Check status code first
	if response.StatusCode() < 200 || response.StatusCode() >= 300 {
		return nil, fmt.Errorf("HTTP %d: %s", response.StatusCode(), response.String())
	}

	// Handle nil result (processing not complete yet)
	resultInterface := response.Result()
	if resultInterface == nil {
		return nil, fmt.Errorf("result not ready yet for message %s", id)
	}

	// Safe type assertion
	result, ok := resultInterface.(*Result)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", resultInterface)
	}

	return result, nil
}

func (cuClient *CuClient) State(processId string) (state []byte, err error) {
	response, err := cuClient.client.R().
		Get("/state/" + processId)

	if err != nil {
		return nil, err
	}

	// Check status code
	if response.StatusCode() < 200 || response.StatusCode() >= 300 {
		return nil, fmt.Errorf("HTTP %d: %s", response.StatusCode(), response.String())
	}

	return response.Body(), nil  // Return the raw bytes
}

func NewCuClient(cuUrl string) CuAPIs {
	return &CuClient{
		client: resty.New().
			SetBaseURL(cuUrl).
			SetRedirectPolicy(resty.FlexibleRedirectPolicy(10)).
			SetError(&Result{}),
		CuUrl: cuUrl,
	}
}