`go-comdirect`
===
![version](https://img.shields.io/github/v/release/jsattler/go-comdirect?include_prereleases)
[![Apache License v2](https://img.shields.io/github/license/jsattler/go-comdirect)](http://www.apache.org/licenses/)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/jsattler/go-comdirect)](https://github.com/jsattler/go-comdirect)

`go-comdirect` is both a client library and CLI tool to interact with
the [comdirect REST API](https://www.comdirect.de/cms/kontakt-zugaenge-api.html).

> **Additional Notes**
> * The library is currently unstable and will change frequently until version 1.0.0 is released
> * Please read the comdirect API documentation prior to using this software
> * Use of this software is at your own risk

Install
---
Use `go get` to install the latest version of this library:
```bash
$ go get -u github.com/jsattler/go-comdirect
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
* on how to use the comdirect CLI tool
* on this [website](https://jsattler.github.io/go-comdirect/#/)
* in the [`docs/`](docs/getting-started.md) folder
* or in the tests of [`pkg/comdirect`](pkg/comdirect)

Roadmap / To-Do
---

* [x] Auth
  * [x] P_TAN_PUSH
  * [ ] P_TAN_PHOTO (currently out of scope)
  * [ ] P_TAN_APP (currently out of scope)
  * [x] Refresh Token
  * [x] Revoke Token
* [x] Account
  * [x] All Balances
  * [x] Balance by Account ID
  * [x] Transactions
* [x] Depot
  * [x] All Depots
  * [x] Positions by Depot ID
  * [x] Position by Depot ID and Position ID
  * [x] Transactions by Depot ID
* [x] Instrument
  * [x] Instrument by Instrument ID
* [ ] Order
  * [x] Dimensions
  * [ ] Orders by Depot ID
  * [ ] Order by Order ID
  * [ ] Order Pre-Validation
  * [ ] Order Validation
  * [ ] Generate Order Cost Indication Ex-Ante
* [ ] Quote
  * [ ] Quote-Request Validation and TAN Challenge Request
  * [ ] Quote-Request Validation and TAN validation
  * [ ] Create Quote-Request Request
* [x] Documents
  * [x] Documents (Postbox)
  * [x] Document by ID
  * [ ] Pre-Document (currently out of scope, since I have no document to test this with)
* [x] Reports
  * [x] Balances of all comdirect products
