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
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: "0.0.1-alpha",
	Use:     "qonto",
	Short:   "qonto is a CLI tool to interact with Qonto banking services (https://qonto.eu/)",
	Long: `
qonto is a CLI (command line interface) to interact with Qonto banking services (https://qonto.eu/)

qonto cli needs to setup some configuration before running.
At least you must provide your Qonto ID and your Qonto secret API key (see "integration" on your Qonto dashboard).

You have two ways to proceed:

1 - Using a config file (needed if you want to use the 'watch' command with email alert).
- Get config.sample.yaml from: https://raw.githubusercontent.com/toorop/go-qonto/master/qonto/config.sample.yaml
- Save it as config.yaml in the same path as your qonto binarie (you can save it elsewhere but in this case you need ton provide its full path with the -c option).
- Fill the required fields.
- That's it. 

WARNING: if you put qonto binary on your binaries PATH, you need to use the --config option or to put your config file on your working dir.

2 - Using environment variables (for basic usages)
- export QONTO_LOGIN & QONTO_SECRET

If you need new commands or features please open an issue on the github repository: https://github.com/toorop/go-qonto

`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("unable to read config:", err)
		os.Exit(1)
	}

	// from ENV
	viper.SetEnvPrefix("qonto")
	viper.AutomaticEnv()
	// check config
	if viper.GetString("login") == "" {
		fmt.Println("config: 'login' is missing !")
		os.Exit(1)
	}
	if viper.GetString("secret") == "" {
		fmt.Println("config: 'secret' is missing !")
		os.Exit(1)
	}
}
