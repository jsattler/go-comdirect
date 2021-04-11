`comdirect-golang`
===
`comdirect-golang` is a client library to interact with
the [comdirect REST API](https://www.comdirect.de/cms/kontakt-zugaenge-api.html).

> **Additional Notes**
> * The library is currently unstable and will change frequently until version 1.0.0 is released
> * Please read the comdirect API documentation prior to using this software
> * Use of this software is at your own risk

Documentation
---
You can find more detailed documentation
* on the [website](https://jsattler.github.io/comdirect-golang/#/)
* in the [`docs/`](docs/getting-started.md) folder
* or in the [`examples/`](examples) older

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
    * [ ] All Depots
    * [ ] Positions by Depot ID
    * [ ] Position by Depot ID and Position ID
    * [ ] Transactions by Depot ID
* [ ] **Instrument**
* [ ] **Order**
* [ ] Quote
* [ ] **Documents**
* [ ] Reports