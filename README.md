Go API client and CLI for [Qonto](https://qonto.eu/) banking service.

## DISCLAIMER
This package and the CLI are not provided nor supported by Qonto

## DISCLAMER 2
Work in progress, don't use it for now.

## go-qonto package quickstart

```go get github.com/toorop/go-qonto```

```go
import "github.com/toorop/go-qonto"

func main(){
    Q := qonto.New("qonto-login", "qonto-API-secret")
    organization, err := Q.GetOrganization("slug")
	if err != nil {
		fmt.Println("ERROR ! unable to get organization -", err)
		os.Exit(1)
	}
	fmt.Println(organization)
}
```

To get last transactions:

```go
import "github.com/toorop/go-qonto"

func main(){
    Q := qonto.New("qonto-login", "qonto-API-secret")
	options := qonto.GetTransactionOptions{
		Slug:   "slug",
		Iban:   "iban",
		Status: []string{"pending", "reversed", "declined", "completed"},
    }
    transactions, err := Q.GetTransactions(options)
	if err != nil {
		fmt.Println("ERR: ", err)
		os.Exit(1)
    }
    for _, transaction := range transactions {
        fmt.Println(transaction)
    }
}
```

## qonto CLI

### Installation
 qonto CLI is available as compiled binaries for Windows, Mac OS, and Linux 

### Configuration
You need to setup some configuration before running.

At least you must provide your Qonto ID and your Qonto secret API key (see "integration" on your Qonto dashboard).

You have two ways to proceed:

* Using a config file (needed if you want to use the 'watch' command with email alert).
    - Get config.sample.yaml from: https://raw.githubusercontent.com/toorop/go-qonto/master/qonto/config.sample.yaml
    - Save it as config.yaml in the same path as your qonto binarie (you can save it elsewhere but in this case you need ton provide its full path with the -c option).
    - Fill the required fields.
    - That's it. 

WARNING: if you put qonto binary on your binaries PATH, you need to use the --config option or to put your config file on your working dir.

- Using environment variables (for basic usages)
    - export QONTO_LOGIN & QONTO_SECRET

If you need new commands or features please open an issue.

### help
qonto CLI have a builtin help, to access it you just have to add a *--help* flag.

Example:

```
$ qonto organization --help

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

Usage:
  qonto organization [flags]

Flags:
  -h, --help   help for organization

Global Flags:
  -c, --config string   config file
```


### watch command

*watch* command allow you to be informed on new transations or transaction updates.

By default *watch command* display update on stdout, but it can send a email and/or send a POST request, containing the transaction update, to a predefined URL.

To keep the *watch* command running when you logout, add "&" at the and of the command ou use [tmux](https://github.com/tmux/tmux/wiki)

#### email notifications

To receive email notifications you must configure *smtp* section of the config file to provide at least *smtp.host*, *smtp.port* and *smtp.mailfrom*

Example:

To receive an email notification on each event

```
qonto watch --slug SLUG --iban IBAN -m EMAIL_ADDRESS_TO_SEND_MAIL_TO
```


#### webhook

To enable webhook notification you just have to add the *--webhook* with a valid URL

```
qonto watch --slug SLUG --ib IBAN --webhook https://qonto.toorop.fr -m EMAIL_ADDRESS_TO_SEND_MAIL_TO
```

qonto CLI will do a POST request to this URL with the JSON encoded transaction in request body for each event.

The format of the JSON object is same as the one returned by [Qonto API](https://api-doc.qonto.eu/2.0/models/transaction)


### organization command

*organization* command return details about your organization and banks accounts.

Example:
```
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

```

## Support this project
If this project is useful for you, please consider making a donation.

### Bitcoin

Address: 1JvMRNRxiTiN9H7LyZTq4yzR7ez86M7ND6

![Bitcoin QR code](https://raw.githubusercontent.com/toorop/wallets/master/btc.png)


### Ethereum

Address: 0xA84684B45969efbD54fd25A1e2eD8C7790A0C497

![ETH QR code](https://raw.githubusercontent.com/toorop/wallets/master/eth.png)

