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
