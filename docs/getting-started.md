Getting Started
---
The following sections will help you to get started quickly.

### Install
Use `go get` to install the latest version of this library:

```bash
$ go get -u github.com/jsattler/comdirect-golang
```

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

**Authenticate using the `Authenticator`**

```go
authentication, err := authenticator.Authenticate()
```
The `Authentication` struct holds all relevant information for subsequent requests to the API.

### First Steps with comdirect.Client

**Create a new `Client` from `AuthOptions`**

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

**Create a new `Client` from an `Authenticator`**
```go
// omitting error validation, imports and packages

options := &comdirect.AuthOptions{
    Username:     os.Getenv("COMDIRECT_USERNAME"),
    Password:     os.Getenv("COMDIRECT_PASSWORD"),
    ClientId:     os.Getenv("COMDIRECT_CLIENT_ID"),
    ClientSecret: os.Getenv("COMDIRECT_CLIENT_SECRET"),
}

authenticator := options.NewAuthenticator()

client := comdirect.NewWithAuthenticator(authenticator)
```
