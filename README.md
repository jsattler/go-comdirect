`comdirect-golang`
===
`comdirect-golang` is a client library to interact with
the [comdirect REST API](https://www.comdirect.de/cms/kontakt-zugaenge-api.html).

> **Additional Notes**
> * The library is currently unstable and will change frequently until version 1.0.0 is released
> * Please read the comdirect API documentation prior to using this software
> * Use of this software is at your own risk

Install
---
Use `go get` to install the latest version of this library:

```bash
$ go get -u github.com/jsattler/comdirect-golang
```

Examples
---
In the following examples we are reading the comdirect credentials from predefined environment variables.

### Authentication

**Creating a new Authenticator from AuthOptions**:

```go
// omitting error validation, imports and packages

options := &comdirect.AuthOptions{
    Username:     os.Getenv("COMDIRECT_USERNAME"),
    Password:     os.Getenv("COMDIRECT_PASSWORD"),
    ClientId:     os.Getenv("COMDIRECT_CLIENT_ID"),
    ClientSecret: os.Getenv("COMDIRECT_CLIENT_SECRET"),
}
authenticator := options.NewAuthenticator()
```

**Creating a new Authenticator with AuthOptions**:

```go
// omitting error validation, imports and packages

options := &comdirect.AuthOptions{
    Username:     os.Getenv("COMDIRECT_USERNAME"),
    Password:     os.Getenv("COMDIRECT_PASSWORD"),
    ClientId:     os.Getenv("COMDIRECT_CLIENT_ID"),
    ClientSecret: os.Getenv("COMDIRECT_CLIENT_SECRET"),
}

authenticator := comdirect.NewAuthenticator(options)
```

### Client Creation

**Create a new Client from AuthOptions**

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

Roadmap / To-Do
---
> Bold items have priority

**Functional**
* [x] **Auth**
    * [x] P_TAN_PUSH
    * [ ] P_TAN_PHOTO (currently out of scope, since the package is not intended for use in front end apps)
    * [ ] P_TAN_APP (currently out of scope, since I have no chance to test this)
* [x] Refresh Token Flow
* [x] Revoke Token
* [x] **Account**
    * [x] All Balances
    * [x] Balance by Account ID
    * [x] Transactions
* [ ] **Depot**
* [ ] **Instrument**
* [ ] **Order**
* [ ] Quote
* [ ] **Documents**
* [ ] Reports