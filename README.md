# Getting Started with Go Core Runtime
[![Go Reference](https://pkg.go.dev/badge/github.com/apimatic/go-core-runtime.svg)](https://pkg.go.dev/github.com/apimatic/go-core-runtime)
[![GitHub release](https://img.shields.io/github/v/release/apimatic/go-core-runtime)](https://pkg.go.dev/github.com/apimatic/go-core-runtime?tab=versions)
[![Licence][license-badge]][license-url]
[![Tests][test-badge]][test-url]
[![Test Coverage](https://api.codeclimate.com/v1/badges/2c5a5f8dca8e970ac36e/test_coverage)](https://codeclimate.com/github/apimatic/go-core-runtime/test_coverage)
[![Maintainability](https://api.codeclimate.com/v1/badges/2c5a5f8dca8e970ac36e/maintainability)](https://codeclimate.com/github/apimatic/go-core-runtime/maintainability)

## Introduction

The `go-core-runtime` is a core library for Apimatic's Go SDK, providing essential utilities and structures for handling API requests and responses using the HTTP protocol. For detailed API documentation and usage examples, visit the [GoDoc documentation](https://pkg.go.dev/github.com/apimatic/go-core-runtime).

## Requirement

This package requires Go v1.18 or higher.

## Installation

To install the package, you can use `go get`:

```bash
go get github.com/apimatic/go-core-runtime
```

Alternatively, you can add the package manually to your go.mod file:


```go
require "github.com/apimatic/go-core-runtime" v0.0.x
```
Then, run the following command to install the package automatically:

```go
go get ./...
```

## Package Details 

### API Error
The apiError package provides a structure to represent error responses from API calls.

| File Name                                                                        | Description                                                           |
|-----------------------------------------------------------------------------|-----------------------------------------------------------------------|
| [`API Error`](apiError/apiError.go)   | Provides a structure to represent error responses from API calls.                        |

### HTTPS
The https package provides logic related to HTTP requests, including building and making the request call. It offers features such as handling form data, file parameters, headers, and interceptors.

| File Name                                           | Description                                                                                      |
| --------------------------------------------------- | ------------------------------------------------------------------------------------------------ |
| [`API Response`](https/apiResponse.go)              | Provides a struct around the HTTP response and the data.                                                         |
| [`Call Builder`](https/callBuilder.go)              | Provides the logic related to the HTTPs request. Includes building and making the request call. |
| [`File Wrapper`](https/fileWrapper.go)              | Provides a wrapper for file parameters to use in the HTTPs calls.                               |
| [`Form Data`](https/formData.go)                    | Provides handling of form parameters in the request.                                             |
| [`HTTP Client`](https/httpClient.go)                | Provides an interface for the HTTP Client to use for making the calls.                           |
| [`HTTP Configuration`](https/httpConfiguration.go)  | Provides configurations for the HTTP calls.                                                     |
| [`HTTP Context`](https/httpContext.go)              | Provides a struct that holds request and corresponding response instances.                      |
| [`HTTP Headers`](https/httpHeaders.go)              | Provides handling for headers to send with the request.                                         |
| [`Internal Error`](https/internalError.go)          | Provides handling for internal errors that may occur during the API calls.                                              |
| [`Interceptors`](https/interceptors.go)             | Provides handling to intercept requests.                                                        |
| [`Retryer`](https/retryer.go)                       | Provides handling to automatically retry for failed requests.                                                  |

### Test Helper
Package testHelper provides helper functions for testing purposes.
| File Name                                   | Description                                                      |
|---------------------------------------------|------------------------------------------------------------------|
| [`BodyMatchers`](testHelper/bodyMatchers.go)           | Provides functions to match JSON response bodies with expected bodies.      |
| [`HeadersMatchers`](testHelper/headersMatchers.go)     | Provides functions to match HTTP headers with expected headers.            |
| [`StatusCodeMatchers`](testHelper/statusCodeMatchers.go) | Provides functions to match HTTP status codes with expected status codes.  |

### Types
Package types provides utility types and functions.

| File Name                                                                        | Description                                                           |
|-----------------------------------------------------------------------------|-----------------------------------------------------------------------|
| [`Optional`](types/optional.go)   | Provides a wrapper to use any type as optional and nullable.   


### Utilities
The utilities package provides utility functions for making HTTP calls.

| File Name                                                                        | Description                                                           |
|-----------------------------------------------------------------------------|-----------------------------------------------------------------------|
| [`API Helper`](utilities/apiHelper.go)   | Provides helper methods for making the HTTP calls.                        |


Each package contains its test files.


## Contributing
Contributions are welcome! If you encounter any issues or have suggestions for improvement, please open an issue.

## License
This project is licensed under the [MIT License](LICENSE).


## Contact
For any questions or support, please feel free to contact us at support@apimatic.io.


[license-badge]: https://img.shields.io/badge/licence-MIT-blue
[license-url]: LICENSE
[test-badge]: https://github.com/apimatic/go-core-runtime/actions/workflows/test.yaml/badge.svg
[test-url]: https://github.com/apimatic/go-core-runtime/actions/workflows/test.yaml
