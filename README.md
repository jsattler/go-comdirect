go-comdirect
===
![version](https://img.shields.io/github/v/release/jsattler/go-comdirect?include_prereleases)
[![Apache License v2](https://img.shields.io/github/license/jsattler/go-comdirect)](http://www.apache.org/licenses/)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/jsattler/go-comdirect)](https://github.com/jsattler/go-comdirect)

`go-comdirect` is both a client library and [CLI tool](comdirect) to interact with
the [comdirect REST API](https://www.comdirect.de/cms/kontakt-zugaenge-api.html).

> **Additional Notes**
> * The library is currently unstable and will change frequently until version 1.0.0 is released
> * Please read the comdirect API documentation prior to using this software
> * Use of this software is at your own risk
> * 10 requests per second are allowed by comdirect
> * 3 invalid TAN validation attempts will cancel the online access 

Features
---
* **Auth:** Authenticate and authorize with the comdirect API.
* **Account:** Access your account data like balances or transactions.
* **Depot:** Access your depot data like balances, positions or transactions.
* **Instrument:** Access instrument data by providing a WKN, ISIN or symbol.
* **Order:** create, modify and delete orders.
In addition, you can query the order book and the status of individual orders, as well as view the display the cost statement for an order. (open #8)
* **Quote:** Do live trading and prepare the query of a quote or execute it (open #9).
* **Documents:** Access and download Postbox-documents.
* **Reports:** Access aggregated reports for multiple of your comdirect products.

Install
---
Use `go get` to install the latest version of this library:
```bash
$ go get -u github.com/jsattler/go-comdirect
```

Use `go install` to install the `comdirect` CLI tool:
```shell
go install github.com/jsattler/go-comdirect/comdirect@main
```

Quick Start
---
```go
// omitting error validation, imports and packages

options := &comdirect.AuthOptions{
    Username:     os.Getenv("COMDIRECT_USERNAME"),
    Password:     os.Getenv("COMDIRECT_PASSWORD"),
    ClientId:     os.Getenv("COMDIRECT_CLIENT_ID"),
    ClientSecret: os.Getenv("COMDIRECT_CLIENT_SECRET"),
}

client := comdirect.NewWithAuthOptions(options)
```

Documentation
---
You can find detailed documentation
* [on how to install and use the comdirect CLI tool](comdirect/README.md)
* on this [website](https://jsattler.github.io/go-comdirect/#/)
* in the [`docs/`](docs/getting-started.md) folder
* or in the tests of [`pkg/comdirect`](pkg/comdirect)


## Contributing
Your contributions are appreciated! Please refer to [CONTRIBUTING.md](CONTRIBUTING.md) for further information.

## License
Please refer to [LICENSE](LICENSE) for further information.