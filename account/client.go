package account

import "context"

// ResponseError provides details about a bad record
type ResponseError struct {
	Detail string
}

func (err ResponseError) Error() string {
	return err.Detail
}

// Account represents an individual account
type Account struct {
	AccountId int    `json:"account_id"`
	Status    string `json:"status"`
	CreatedOn string `json:"created_on"`
}

// GetAccountsRequest is the request object to get all the account data.  This
// object is empty for now, but can be expanded to handle additional input, like
// to provide page info
type GetAccountsRequest struct {
}

// GetAccountsResponse is the response object for all account data.  The
// response object has fields for "next" and "previous", but since I can't glean
// the shape of the data, nor is that info relevant for the scope of this task,
// I will leave it off for now.
type GetAccountsResponse struct {
	Results []*Account
}

// GetAccountRequest is the request object to retrieve a single account
type GetAccountRequest struct {
	AccountId string
}

// Client manages requests to the wpengine server.  An interface is used here to
// facilitate unit testing and could make it easier to add additional layers for
// things like caching
type Client interface {
	// GetAccounts retrieves all accounts on the server
	GetAccounts(context.Context, *GetAccountsRequest) (*GetAccountsResponse, error)
	// GetAccount retrieves a single account
	GetAccount(context.Context, *GetAccountRequest) (*Account, error)
}
