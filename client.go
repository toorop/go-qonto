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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	// Qonto API endpoint
	endpoint = "https://thirdparty.qonto.eu/v2"
	// client timeout in seconds
	clientTimeout = 15
)

// Client is the client to interact with Qonto REST API
type Client struct {
	client   *http.Client
	login    string
	secret   string
	endpoint string
}

// New returns a Qonto client
func New(login, secret string) (qonto Client) {
	qonto.client = new(http.Client)
	qonto.client.Timeout = clientTimeout * time.Second
	qonto.login = login
	qonto.secret = secret
	qonto.endpoint = endpoint
	return
}

// do is a wrapper for http.client.Do wich add authentification
func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", c.login+":"+c.secret)
	return c.client.Do(req)
}

// doAndReturnBody is a http.Client.Do wrapper with auth which returns response body as bytes slice
func (c *Client) doAndReturnBody(req *http.Request) (body []byte, err error) {
	response, err := c.do(req)
	if err != nil {
		return body, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return body, fmt.Errorf("request failed - bad HTTP status returned: %s", response.Status)
	}
	return ioutil.ReadAll(response.Body)
}

// GetOrganization is a wrapper that handle GET /organizations/{login} call
func (c *Client) GetOrganization(organizationName string) (organization Organization, err error) {
	req, err := http.NewRequest("GET", c.endpoint+"/organizations/"+url.QueryEscape(organizationName), nil)
	if err != nil {
		return
	}
	resp, err := c.doAndReturnBody(req)
	if err != nil {
		return
	}
	response := new(GetOrganizationResponse)
	err = json.Unmarshal(resp, response)
	return response.Organization, err
}

// GetTransactions is a wrapper that handle GET /transaction call
func (c *Client) GetTransactions(options GetTransactionOptions) (transactions []Transaction, err error) {
	// validate options
	if valid, err := options.isValid(); !valid {
		return transactions, err
	}

	payload := bytes.NewBuffer([]byte("{pending, reversed, declined, completed}"))

	req, err := http.NewRequest("GET", c.endpoint+"/transactions?slug="+url.QueryEscape(options.Slug)+"&iban="+url.QueryEscape(options.Iban), payload)
	if err != nil {
		return
	}

	resp, err := c.doAndReturnBody(req)
	if err != nil {
		return transactions, err
	}
	response := new(getTransactionResponse)
	err = json.Unmarshal(resp, response)
	return response.Transactions, err
}
