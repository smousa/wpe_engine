package account

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

var httpClient = http.DefaultClient

const (
	AccountsEndpoint = "/v1/accounts"
)

// WPClient is the connector to the wpengine server that implements client
type WPClient struct {
	url *url.URL
}

var _ Client = &WPClient{}

// NewWPClient instantiates a new client
func NewWPClient(addr string) *WPClient {
	url, err := url.Parse(addr)
	if err != nil {
		panic("invalid address")
	}
	return &WPClient{url}
}

func (c *WPClient) GetAccounts(ctx context.Context, req *GetAccountsRequest) (*GetAccountsResponse, error) {
	// make a request to the server
	url, err := c.url.Parse(AccountsEndpoint)
	if err != nil {
		// This should not happen
		panic("could not load endpoint")
	}
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create request")
	}
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to the server")
	}
	defer httpResp.Body.Close()
	// decode the body as json
	dec := json.NewDecoder(httpResp.Body)

	// handle non-200 response code
	if httpResp.StatusCode != http.StatusOK {
		var respErr ResponseError
		if err := dec.Decode(&respErr); err != nil {
			return nil, errors.Wrap(err, "could not decode response from the server")
		}
		return nil, errors.Wrap(err, "could not look up accounts")
	}
	var resp GetAccountsResponse
	if err := dec.Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "could not decode response from the server")
	}
	return &resp, nil
}

func (c *WPClient) GetAccount(ctx context.Context, req *GetAccountRequest) (*Account, error) {
	// make a request to the server
	url, err := c.url.Parse(path.Join(AccountsEndpoint, req.AccountId))
	if err != nil {
		// This should not happen
		panic("could not load endpoint")
	}
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create request")
	}

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to the server")
	}
	defer httpResp.Body.Close()
	// decode body as json
	dec := json.NewDecoder(httpResp.Body)

	// handle non-200 response code
	if httpResp.StatusCode != http.StatusOK {
		var respErr ResponseError
		if err := dec.Decode(&respErr); err != nil {
			return nil, errors.Wrap(err, "could not decode response from the server")
		}
		return nil, errors.Wrap(respErr, "could not look up account")
	}
	var resp Account
	if err := dec.Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "could not decode response from the server")
	}
	return &resp, nil
}
