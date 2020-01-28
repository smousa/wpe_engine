package account_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"

	. "github.com/wpe_merge/wpe_merge/account"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var emulator *WPEmulator

func TestAccount(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Account Suite")
}

var _ = BeforeSuite(func() {
	emulator = NewWPEmulator()
})

var _ = AfterSuite(func() {
	emulator.Close()
})

// WPEmulator is a test server so I can focus on functionality of my code rather
// than my network connectivity
type WPEmulator struct {
	svr    *httptest.Server
	router *mux.Router
	data   map[string]*Account
}

func NewWPEmulator() *WPEmulator {
	wpe := &WPEmulator{}

	router := mux.NewRouter()
	router.HandleFunc(AccountsEndpoint, wpe.getAccounts)
	router.HandleFunc(path.Join(AccountsEndpoint, "{id:[0-9]+}"), wpe.getAccount)
	svr := httptest.NewServer(router)

	wpe.svr = svr
	wpe.router = router
	wpe.data = make(map[string]*Account)

	return wpe
}

func (wpe *WPEmulator) Close() {
	wpe.svr.Close()
}

func (wpe *WPEmulator) URL() string {
	return wpe.svr.URL
}

func (wpe *WPEmulator) LoadData(data ...*Account) {
	for i, record := range data {
		key := fmt.Sprintf("%d", record.AccountId)
		wpe.data[key] = data[i]
	}
}

func (wpe *WPEmulator) ResetData() {
	wpe.data = make(map[string]*Account)
}

func (wpe *WPEmulator) getAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var resp GetAccountsResponse
	resp.Results = make([]*Account, 0, len(wpe.data))
	for key := range wpe.data {
		resp.Results = append(resp.Results, wpe.data[key])
	}
	enc := json.NewEncoder(w)
	err := enc.Encode(&resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (wpe *WPEmulator) getAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	accountID := mux.Vars(r)["id"]
	account, ok := wpe.data[accountID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		enc := json.NewEncoder(w)
		err := enc.Encode(ResponseError{Detail: "Not found."})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	enc := json.NewEncoder(w)
	err := enc.Encode(account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
