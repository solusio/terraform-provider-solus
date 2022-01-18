SOLUS IO GoLang SDK
===================

[![Go Reference](https://pkg.go.dev/badge/github.com/solusio/solus-go-sdk.svg)](https://pkg.go.dev/github.com/solusio/solus-go-sdk)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/solusio/solus-go-sdk/main?label=main&logo=github)
[![Report Card](https://goreportcard.com/badge/github.com/solusio/solus-go-sdk)](https://goreportcard.com/report/github.com/solusio/solus-go-sdk)
[![Coverage Status](https://coveralls.io/repos/github/solusio/solus-go-sdk/badge.svg?branch=master)](https://coveralls.io/github/solusio/solus-go-sdk?branch=master)

solus-go-sdk is a Go client for accessing [SOLUS IO API v1](https://docs.solus.io/api/)

SOLUS IO is a virtual infrastructure management solution that facilitates
choice, simplicity, and performance for ISPs and enterprises. Offer blazing
fast, on-demand VMs, a simple API, and an easy-to-use self-service control
panel for your customers to unleash your full potential for growth.

[Official site](https://www.solus.io/)

Usage
-----

```go
client, err := solus.NewClient(baseURL, solus.EmailAndPasswordAuthenticator{
    Email: "email@example.com",
    Password: "12345678",
})
```

Or

```go
client, err := solus.NewClient(baseURL, solus.APITokenAuthenticator{Token: "api token"})
```

Development
-----------

For (re)generating code just run `go generate`
