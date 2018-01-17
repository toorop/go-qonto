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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	qonto "github.com/toorop/go-qonto"
)

// organizationCmd represents the organization command
var organizationCmd = &cobra.Command{
	Use:   "organization",
	Short: "Display your Qonto organizations",
	Long: `
Display your Qonto organizations

Example:

$ qonto organization
my-orga-42
                Slug: my-orga-42-bank-account-1
                IBAN: FR76XXXXXXXXX
                BIC: XXXXXXXX
                Currency: EUR
                Balance (€): 100000000.000000
                Balance (cents): 10000000000
                Authorized Balance (€): 0.000000
                Auhorized Balance (cents): 0
	`,
	Run: func(cmd *cobra.Command, args []string) {
		Q := qonto.New(viper.GetString("login"), viper.GetString("secret"))
		organization, err := Q.GetOrganization(viper.GetString("login"))
		if err != nil {
			fmt.Println("ERROR ! unable to get organization -", err)
			os.Exit(1)
		}
		fmt.Println(organization)

	},
}

func init() {
	rootCmd.AddCommand(organizationCmd)
}
