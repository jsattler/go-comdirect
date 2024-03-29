`comdirect` CLI
===
`comdirect` CLI tool lets you interact with your comdirect account, depot, documents and other resources.

> So far only tested on macOS 12.1 and Fedora 33

Install
---
```shell
go install github.com/jsattler/go-comdirect/comdirect@main
```

Getting Started
---
>All commands and subcommands use the singular form, so instead of `accounts` it's `account`.
>This is a convention to make it easier for users to remember the commands.

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
  --id=<clientID> \
  --secret=<clientSecret> \
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

List basic account information

```shell
comdirect account
```

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
Some notes on the current behavior:
* the tool does not check if a file already exists. If it does, it will download and truncate the existing file
* You need to specify the `--download` flag to download the files

List all documents from the postbox
```shell
comdirect document
```

List a specific document
```shell
comdirect document <documentID>
```

Download first 20 documents
```shell
comdirect document --count=20 --download
```

Download document by ID
```shell
comdirect document --download <documentID>
```
