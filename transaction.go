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
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	// ISO8601 date format
	ISO8601 = "2006-01-02T15:04:05-0700"
)

// Transaction represents a qonto transaction model
type Transaction struct {
	ID               string  `json:"transaction_id"`
	Amount           float64 `json:"amount"`
	AmountCents      uint64  `json:"amount_cents"`
	LocalAmount      float64 `json:"local_amount"`
	LocalAmountCents uint64  `json:"local_amount_cents"`
	Side             string  `json:"side"`           // credit | debit
	OperationType    string  `json:"operation_type"` // transfert | card | direct_debit | income | qonto_fee
	Currency         string  `json:"currency"`       // ISO 4217
	LocalCurrency    string  `json:"local_currency"` // ISO 4217
	SettleAt         Qtime   `json:"settle_at"`
	EmittedAt        Qtime   `json:"emitted_at"`
	Status           string  `json:"status"`
	Note             string  `json:"note"`
	Label            string  `json:"label"`
}

// Qtime is time formated as returned by Qonto API
// ISO8601 yyyy-MM-dd'T'HH:mm:ss.SSSZ
type Qtime struct {
	time.Time
}

// UnmarshalJSON is the Qtime unmarshaler
func (t *Qtime) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	t.Time, err = time.Parse(ISO8601, string(b))
	return err
}

func (t *Transaction) String() string {
	return fmt.Sprintf(`
		ID: %s
		Amount: %f
		Amount (cts): %d
		Local amount: %f
		Local amount (cts): %d
		Side: %s
		Operation type: %s 
		Currency: %s
		Local currency: %s
		Emitted at: %s
		Settle at: %s
		Status: %s
		Note: %s
		Label: %s
		`, t.ID, t.Amount, t.AmountCents, t.LocalAmount, t.LocalAmountCents, t.Side, t.OperationType, t.Currency, t.LocalCurrency, t.EmittedAt, t.SettleAt, t.Status, t.Note, t.Label)
}

// DisplayInline return transaction as one line sting
func (t *Transaction) DisplayInline() string {
	return fmt.Sprintf("%s - - Operation: %s - Status: %s -  Side: %s - Amount(cts): %d", t.ID, t.OperationType, t.Status, t.Side, t.AmountCents)
}

////
// struct and tools for HTTP request & response

// GetTransactionOptions -> options for GetTransactions
type GetTransactionOptions struct {
	Slug        string
	Iban        string
	Status      []string
	CurrentPage uint16
	PerPage     uint16
}

func (o *GetTransactionOptions) isValid() (bool, error) {
	// required
	o.Slug = strings.TrimSpace(o.Slug)
	if o.Slug == "" {
		return false, errors.New("parameter Slug is required")
	}
	o.Iban = strings.TrimSpace(o.Iban)
	if o.Iban == "" {
		return false, errors.New("parameter Iban is required")
	}
	// TODO (eventualy) check []status in enum (pending, reversed, declined, completed)
	return true, nil

}

// response to GET /transactions
type getTransactionResponse struct {
	Transactions []Transaction `json:"transactions"`
	Meta         struct {
		CurrentPage uint16 `json:"current_page"`
		NextPage    uint16 `json:"next_page"`
		PrevPage    uint16 `json:"prev_page"`
		TotalPage   uint16 `json:"total_page"`
		TotalCount  uint32 `json:"total_count"`
		PerPage     uint16 `json:"per_page"`
	} `json:"meta"`
}
