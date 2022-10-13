# Getting Started with Go Core Runtime
[![Go Reference](https://pkg.go.dev/badge/github.com/apimatic/go-core-runtime.svg)](https://pkg.go.dev/github.com/apimatic/go-core-runtime)
[![GitHub release](https://img.shields.io/github/v/release/apimatic/go-core-runtime)](https://pkg.go.dev/github.com/apimatic/go-core-runtime?tab=versions)
[![Licence][license-badge]][license-url]
[![Tests Passing](https://github.com/apimatic/go-core-runtime/actions/workflows/test.yaml/badge.svg)](https://github.com/apimatic/go-core-runtime/actions/workflows/test.yaml)

## Introduction

Core library for Apimatic's Go SDK hosted on [github.com/apimatic/go-core-runtime](https://pkg.go.dev/github.com/apimatic/go-core-runtime)

## Requirement

Go v1.18

## Install the Package

Run the following command to install the package and automatically add the package to your module or .mod file:

```go
go get "github.com/apimatic/go-core-runtime"
```

Or add it to the go.mod file manually as given below:

```go
require "github.com/apimatic/go-core-runtime" v0.0.x
```
And run the following command to install the package automatically:

```go
go get ./...
```

## Package Details 
### HTTPS

| File Name                                                                        | Description                                                           |
|-----------------------------------------------------------------------------|-----------------------------------------------------------------------|
| [`Call Builder`](https/callBuilder.go)   | Provides the logic related to the HTTPs request. Includes building and making the request call.                        |
| [`File Wrapper`](https/fileWrapper.go) | Provides a wrapper for the file parameter to use in the HTTPs calls.                    |
| [`Form Data`](https/formData.go) | Provides handling of form parameters in the request.                    |
| [`HTTP Client`](https/httpClient.go) | Provides an interface for the HTTP Client to use for making the calls.                    |
| [`HTTP Configuration`](https/httpConfiguration.go) | Provides configurations for the HTTP calls.                    |
| [`HTTP Context`](https/httpContext.go) | Provides a struct that holds request and corresponding response instances.                    |
| [`HTTP Headers`](https/httpHeaders.go) | Provides handling for headers to send with the request.                    |
| [`Interceptors`](https/interceptors.go) | Provides handling to intercept requests.                    |


### API Error

| File Name                                                                        | Description                                                           |
|-----------------------------------------------------------------------------|-----------------------------------------------------------------------|
| [`API Error`](apiError/apiError.go)   | Provides the error struct that is used in the endpoint calls.                        |


### Utilities

| File Name                                                                        | Description                                                           |
|-----------------------------------------------------------------------------|-----------------------------------------------------------------------|
| [`API Helper`](utilities/apiHelper.go)   | Provides helper methods for making the HTTP calls.                        |


Each package contains its test files and code coverage reports as well.



[license-badge]: https://img.shields.io/badge/licence-APIMATIC-blue
[license-url]: LICENSE
