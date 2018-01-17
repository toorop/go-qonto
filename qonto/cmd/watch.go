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

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	qonto "github.com/toorop/go-qonto"
)

const (
	// email subject
	emailSubject = "[QONTO WATCHER] update for transaction %s"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch an account and display logs or/and send an email or/and call a webhook on each mouvment",
	Long: `
Watch an account and for each new transation (or transaction update):

- display logs on stdout
- send an email (optional)
- call a webhook (optional)

If you want to recieve email notifications, you have ton setup "smtp" section on the config file.

Examples:

1 - Receive notification by email on each update:
qonto watch --slug account-slug --iban IBAN --email toorop@gmail.com

2 - Call a webhook on each update
qonto watch --slug account-slug --iban IBAN --webhook https://qonto.toorop.fr/

3 - Receive email notification and call a webhook (Amazing !)
qonto watch --slug account-slug --iban IBAN --email toorop@gmail.com --webhook https://qonto.toorop.fr/




`,
	Run: watch,
}

func init() {
	rootCmd.AddCommand(watchCmd)

	// slug (required)
	watchCmd.Flags().StringP("slug", "s", "", "slug of the account to be watched (required)")
	viper.BindPFlag("slug", watchCmd.Flags().Lookup("slug"))

	// iban (required)
	watchCmd.Flags().StringP("iban", "i", "", "IBAN of the account to be watched (required)")
	viper.BindPFlag("iban", watchCmd.Flags().Lookup("iban"))

	// statuses
	watchCmd.Flags().StringSlice("statuses", []string{"pending", "reversed", "declined", "completed"}, "statuses you want to bind. Example: --status \"reversed declined\" . By default all type of transactions are watched.")
	viper.BindPFlag("statuses", watchCmd.Flags().Lookup("statuses"))

	// email
	watchCmd.Flags().StringP("email", "m", "", "email addresse where watcher will send notification on change")
	viper.BindPFlag("send-email-to", watchCmd.Flags().Lookup("email"))

	// webhook
	watchCmd.Flags().StringP("webhook", "w", "", "Webhook URL")
	viper.BindPFlag("webhook", watchCmd.Flags().Lookup("webhook"))
}

func watch(cmd *cobra.Command, args []string) {
	// check flags

	// slug
	if viper.GetString("slug") == "" {
		fmt.Println("--slug option is required. qonto watch --help for more details.")
		os.Exit(1)
	}

	// iban
	if viper.GetString("iban") == "" {
		fmt.Println("--iban option is required. qonto watch --help for more details.")
		os.Exit(1)
	}

	// if email, check that at least smtp.host and smtp.port are set
	if viper.GetString("send-email-to") != "" {
		if viper.GetString("smtp.host") == "" {
			fmt.Println("config smtp.host is missing. You need to set smtp options in the config file if you want to receive email notification. qonto watch --help for more details.")
			os.Exit(1)
		}
		if viper.GetString("smtp.port") == "" {
			fmt.Println("config smtp.port is missing. You need to set smtp options in the config file if you want to receive email notification. qonto watch --help for more details.")
			os.Exit(1)
		}
		if viper.GetString("smtp.mailfrom") == "" {
			fmt.Println("config smtp.mailfrom is missing. You need to set smtp options in the config file if you want to receive email notification. qonto watch --help for more details.")
			os.Exit(1)
		}
	}

	// if webhook check url
	if viper.GetString("webhook") != "" {
		if !govalidator.IsURL(viper.GetString("webhook")) {
			fmt.Println("webhook url seems invalid. qonto watch --help for more details.")
			os.Exit(1)
		}
	}

	Q := qonto.New(viper.GetString("login"), viper.GetString("secret"))
	options := qonto.GetTransactionOptions{
		Slug:   viper.GetString("slug"),
		Iban:   viper.GetString("iban"),
		Status: viper.GetStringSlice("statuses"),
	}
	// Warning this basic algo will fail if you have more than 100 transacs per minute
	// if it's the case contact me, i will solve your problem for less than one minute.
	tac := time.Now()
	var tic time.Time
	for {
		// let's start by a little snap
		time.Sleep(60 * time.Second)
		// get last transactions

		transactions, err := Q.GetTransactions(options)
		if err != nil {
			log.Println("ERR: ", err)
			continue
		}
		tic = tac
		// WARNING there is a black hole here !!!
		tac = time.Now()
		for _, transaction := range transactions {
			if transaction.EmittedAt.After(tic) || transaction.SettleAt.After(tic) {
				go handleNewTransaction(transaction)
			}
		}
	}
}

// display logs, send email, callwebhook
func handleNewTransaction(transaction qonto.Transaction) {
	// Log
	log.Println(transaction.DisplayInline())

	// send email
	var auth smtp.Auth
	if viper.GetString("send-email-to") != "" {
		// Auth ?
		if viper.GetString("smtp.user") != "" && viper.GetString("smtp.password") != "" {
			auth = smtp.PlainAuth("", viper.GetString("smtp.user"), viper.GetString("smtp.password"), viper.GetString("smtp.host"))
		}
		// let's go
		msg := []byte(transaction.String())
		if err := smtp.SendMail(fmt.Sprintf("%s:%s", viper.GetString("smtp.host"), viper.GetString("smtp.port")), auth, viper.GetString("smtp.mailfrom"), []string{viper.GetString("send-email-to")}, msg); err != nil {
			log.Println("ERR: unable to send mail - ", err)
		}
	}

	// webhook
	if viper.GetString("webhook") != "" {
		payload, err := json.Marshal(transaction)
		if err != nil {
			log.Println("ERR: ", err)
			return
		}
		resp, err := http.Post(viper.GetString("webhook"), "application/json", bytes.NewBuffer(payload))
		if err != nil {
			log.Println("ERR: ", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Println(fmt.Sprintf("ERR: webhook call has failed - %s ", resp.Status))
		}
	}
}
