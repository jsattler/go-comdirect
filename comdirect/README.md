`comdirect` CLI
===
`comdirect` CLI tool lets you interact with your comdirect account, depot and documents. 

Install
---
```shell
go install github.com/jsattler/go-comdirect
```

Getting Started
---

```text
Usage:
    comdirect [COMMAND] [SUBCOMMAND] [OPTIONS] [ID]
```

### Authentication
To log in you can specify the credentials through the options or when prompted.
> The login command will try to store the credentials in one of the following locations depending on your OS
> * OS X KeyChain
> * Secret Service dbus interface (GNOME Keyring)
> * Windows Credentials Manager


```shell
comdirect login \
  --clientID=<clientID> \
  --clientSecret=<clientSecret> \
  --username=<username> \
  --password=<password>
```
or 
```shell
comdirect login
```

The logout command will remove all stored credentials, access and refresh tokens from the mentioned credential providers.

```shell
comdirect logout 
```

### Account

List all account information and balances (giro Konto, tagesgeldplus etc.)

```shell
comdirect account balance
```

Retrieve account information and balances for a specific account

```shell
comdirect account balance <accountID>
```

Retrieve account transactions for a specific account
```shell
comdirect account transaction <accountID>
```

### Depot

Retrieve *depot* information 

```shell
comdirect depot
```

Retrieve *depot* positions for a specific depot

```shell
comdirect depot position <depotID>
```

Retrieve a specific depot position for a specific depot

```shell
comdirect depot position --position=<positionID> <depotID>
```

Retrieve all transactions for a specific depot

```shell
comdirect depot transaction <depotID>
```

### Document 

Retrieve all documents from the postbox
```shell
comdirect document
```

Retrieve a specific document
```shell
comdirect document <documentID>
```

Security
---
