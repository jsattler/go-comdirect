`comdirect-golang`
===
`comdirect-golang` is a client library to interact with the [comdirect REST API](https://www.comdirect.de/cms/kontakt-zugaenge-api.html).

> The library is unstable at the moment and will change frequently until version 1.0.0 is released.

Install
---
To use the client within your application simply run

```bash
$ go get github.com/j-sattler/comdirect-golang
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
    AutoRefresh:  true,
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
    AutoRefresh:  true,
}

authenticator := comdirect.NewAuthenticator(options)
```


TO DO
---

* [ ] Auth
    * [ ] P_TAN_PUSH
    * [ ] P_TAN_PHOTO
    * [ ] P_TAN_APP
* [ ] Refresh Token Flow
* [ ] Revoke Token
* [ ] Account
* [ ] Depot
* [ ] Instrument
* [ ] Order
* [ ] Quote
* [ ] Documents
* [ ] Reports

Third Party Libraries
---
