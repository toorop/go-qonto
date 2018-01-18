// Copyright © 2018 Stéphane Depierrepont
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package qonto

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	getOrganizationResponse     = `{"organization":{"slug":"slug","bank_accounts":[{"slug":"bank-account-1","iban":"IBAN","bic":"BIC","currency":"EUR","balance":0.0,"balance_cents":0,"authorized_balance":0.0,"authorized_balance_cents":0}]}}`
	getTransactionsResponse     = `{"transactions":[{"transaction_id":"bank-account-1-transaction-1","amount":100000000.0,"amount_cents":10000000000,"local_amount":100000000.0,"local_amount_cents":10000000000,"side":"credit","operation_type":"income","currency":"EUR","local_currency":"EUR","label":"Present from Elon Musk","settled_at":"2018-01-18T06:45:57.000Z","emitted_at":"2018-01-18T07:45:58.000Z","status":"completed","note":null}],"meta":{"current_page":1,"next_page":null,"prev_page":null,"total_pages":1,"total_count":1,"per_page":100}}`
	testDoAndReturnBodyResponse = "yop"
)

func get401TestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))

}

// TestDo - check if auth token are correctly setted
func TestAuth(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "login:secret", r.Header.Get("Authorization"))
		fmt.Fprintln(w, "pong")
	}))
	defer ts.Close()
	Q := New("login", "secret")
	Q.endpoint = ts.URL
	req, _ := http.NewRequest("GET", ts.URL, nil)
	_, err := Q.do(req)
	assert.NoError(t, err)
}

// TestUnauthorized check that 401 is correctly handled
func TestUnauthorized(t *testing.T) {
	ts := get401TestServer()
	defer ts.Close()
	Q := New("login", "secret")
	Q.endpoint = ts.URL
	_, err := Q.GetOrganization("foo")
	assert.EqualError(t, err, "request failed - bad HTTP status returned: 401 Unauthorized")
}

// TestDoAndReturnBody test that doAndReturnBody return expected []byte
func TestDoAndReturnBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, testDoAndReturnBodyResponse)
	}))
	defer ts.Close()
	Q := New("login", "secret")
	Q.endpoint = ts.URL
	req, _ := http.NewRequest("GET", ts.URL, nil)
	body, _ := Q.doAndReturnBody(req)
	assert.Equal(t, testDoAndReturnBodyResponse+"\n", string(body))
}

func TestGETOrganization(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, getOrganizationResponse)
	}))
	defer ts.Close()
	Q := New("login", "secret")
	Q.endpoint = ts.URL
	orga, err := Q.GetOrganization("foo")
	assert.NoError(t, err)
	assert.Equal(t, "slug", orga.Slug)
	assert.Equal(t, "bank-account-1", orga.BankAccounts[0].Slug)
	assert.Equal(t, "IBAN", orga.BankAccounts[0].Iban)
	assert.Equal(t, "EUR", orga.BankAccounts[0].Currency)
	assert.Equal(t, "BIC", orga.BankAccounts[0].Bic)
	assert.Equal(t, 0.0, orga.BankAccounts[0].Balance)
	assert.Equal(t, 0, orga.BankAccounts[0].BalanceCents)
	assert.Equal(t, 0.0, orga.BankAccounts[0].AuthorizedBalance)
	assert.Equal(t, 0, orga.BankAccounts[0].AuthorizedBalanceCents)

}

func TestGetTransaction(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, getTransactionsResponse)
	}))
	defer ts.Close()
	Q := New("login", "secret")
	Q.endpoint = ts.URL
	options := GetTransactionOptions{
		Slug: "slug",
		Iban: "iban",
	}
	transactions, err := Q.GetTransactions(options)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(transactions))
	tx := transactions[0]
	assert.Equal(t, "bank-account-1-transaction-1", tx.ID)
	assert.Equal(t, 100000000.0, tx.Amount)
	assert.Equal(t, uint64(10000000000), tx.AmountCents)
	assert.Equal(t, 100000000.0, tx.LocalAmount)
	assert.Equal(t, uint64(10000000000), tx.LocalAmountCents)
	assert.Equal(t, "credit", tx.Side)
	assert.Equal(t, "income", tx.OperationType)
	assert.Equal(t, "EUR", tx.Currency)
	assert.Equal(t, "EUR", tx.LocalCurrency)
	assert.Equal(t, "Present from Elon Musk", tx.Label)
	assert.Equal(t, "Present from Elon Musk", tx.Label)
	expectedT, _ := time.Parse(ISO8601, "2018-01-18T06:45:57.000Z")
	assert.Equal(t, expectedT, tx.SettleAt.Time)
	expectedT, _ = time.Parse(ISO8601, "2018-01-18T07:45:58.000Z")
	assert.Equal(t, expectedT, tx.EmittedAt.Time)
	assert.Equal(t, "completed", tx.Status)
	assert.Equal(t, "", tx.Note)
}
