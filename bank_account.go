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

import "fmt"

// BankAccount represent qonto "Bank Account" model
type BankAccount struct {
	Slug                   string  `json:"slug"`
	Iban                   string  `json:"iban"`
	Bic                    string  `json:"bic"`
	Currency               string  `json:"currency"`
	Balance                float64 `json:"balance"`
	BalanceCents           int     `json:"balance_cents"`
	AuthorizedBalance      float64 `json:"authorized_balance"`
	AuthorizedBalanceCents int     `json:"authorized_balance_cents"`
}

// String is a stringer for bankAccount stuct
func (b *BankAccount) String() string {
	return fmt.Sprintf(`
		Slug: %s
		IBAN: %s
		BIC: %s
		Current: %s
		Balance (€): %f
		Balance (cents): %d
		Authorized Balance (€): %f
		Auhorized Balance (cents): %d
		`, b.Slug, b.Iban, b.Bic, b.Currency, b.Balance, b.BalanceCents, b.AuthorizedBalance, b.AuthorizedBalanceCents)
}
