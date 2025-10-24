# Getting Started with Go Core Runtime
[![Go Reference](https://pkg.go.dev/badge/github.com/apimatic/go-core-runtime.svg)](https://pkg.go.dev/github.com/apimatic/go-core-runtime)
[![GitHub release](https://img.shields.io/github/v/release/apimatic/go-core-runtime)](https://pkg.go.dev/github.com/apimatic/go-core-runtime?tab=versions)
[![Test Coverage][coverage-badge]][coverage-url]
[![Maintainability Rating][maintainability-badge]][maintainability-url]
[![Vulnerabilities][vulnerabilities-badge]][vulnerabilities-url]
[![Licence][license-badge]][license-url]
[![Tests][test-badge]][test-url]


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
Following is a list of all packages. Here each package contains its own test files.

### HTTPS
The https package provides logic related to HTTP requests, including building and making the request call. It offers features such as handling form data, file parameters, headers, and interceptors.

| File Name                                            | Description                                                                                     |
|------------------------------------------------------|-------------------------------------------------------------------------------------------------|
| [`API Response`](https/apiResponse.go)               | Provides a struct around the HTTP response and the data.                                        |
| [`API Error`](https/apiError.go)                     | Provides a structure to represent error responses from API calls.                               |
| [`Array Serialization`](https/arraySerialization.go) | Provides constants for array serialization like Indexed, Plain, etc.                            |
| [`Auth Group`](https/authenticationGroup.go)         | Provides capability to group authentication keys into groups.                                   |
| [`Auth Interface`](https/authenticationInterface.go) | Defines the functionality of an authentication manager.                                         |
| [`Call Builder`](https/callBuilder.go)               | Provides the logic related to the HTTPs request. Includes building and making the request call. |
| [`File Wrapper`](https/fileWrapper.go)               | Provides a wrapper for file parameters to use in the HTTPs calls.                               |
| [`Form Data`](https/formData.go)                     | Provides handling of form parameters in the request.                                            |
| [`HTTP Client`](https/httpClient.go)                 | Provides an interface for the HTTP Client to use for making the calls.                          |
| [`HTTP Configuration`](https/httpConfiguration.go)   | Provides configurations for the HTTP calls.                                                     |
| [`HTTP Context`](https/httpContext.go)               | Provides a struct that holds request and corresponding response instances.                      |
| [`HTTP Headers`](https/httpHeaders.go)               | Provides handling for headers to send with the request.                                         |
| [`Internal Error`](https/internalError.go)           | Provides handling for internal errors that may occur during the API calls.                      |
| [`Interceptors`](https/interceptors.go)              | Provides handling to intercept requests.                                                        |
| [`Retryer`](https/retryer.go)                        | Provides handling to automatically retry for failed requests.                                   |

### Logger
The logger package provides logic related to logging. It offers the Facade Design Pattern for configuring the Logger and SDK Logger. Additionally, it provides the LoggerConfiguration to customize logging behavior.

| File Name                                                                | Description                                                                                                                   |
|--------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------|
| [`Console Logger`](logger/defaultLogger.go)                              | Provides default implementation of [`Logger Interface`](logger/defaultLogger.go) to log messages.                             |
| [`Level`](logger/level.go)                                               | Provides constants for log level like Level_ERROR, Level_INFO, etc.                                                           |
| [`Logger Configuration`](logger/loggerConfiguration.go)                  | Provides logging configurations for the Sdk Logger.                                                                           |
| [`Response Logger Configuration`](logger/responseLoggerConfiguration.go) | Provides response logging configurations for the Sdk Logger.                                                                  |
| [`Request Logger Configuration`](logger/requestLoggerConfiguration.go)   | Provides request logging configurations for the Sdk Logger.                                                                   |
| [`Sdk Logger`](logger/sdkLogger.go)                                      | Provides default and null implementation of [` Sdk Logger Interface`](logger/sdkLogger.go) to log API requests and responses. |

### Security
The security package provides logic related to request authentication and verification. It includes digest encoding/decoding, signature verification, and HMAC-based request validation.

| File Name                                                                                     | Description                                                                                                                                                                                         |
|-----------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [`DigestCodec`](security/digest_codec.go)                                                     | Provides digest encoding/decoding interfaces and implementations for Hex, Base64, and Base64URL formats.                                                                                            |
| [`SignatureVerifier`](security/hmac_signature_verifier.go)                                    | Provides interfaces for verifying signatures for HTTP requests.                                                                                                                             |
| [`HmacSignatureVerifier`](security/hmac_signature_verifier.go)                                | Provides an implementation of [`SignatureVerifier`](security/hmac_signature_verifier.go) for HMAC-based signature verification for HTTP requests, supporting configurable algorithms and templates. |
| [`VerificationResult`](security/verification_result.go)                                       | Defines the structure for representing signature verification outcomes with success/failure states and error details.                                                                               |


### Test Helper
Package testHelper provides helper functions for testing purposes.

| File Name                                                | Description                                                               |
|----------------------------------------------------------|---------------------------------------------------------------------------|
| [`BodyMatchers`](testHelper/bodyMatchers.go)             | Provides functions to match JSON response bodies with expected bodies.    |
| [`HeadersMatchers`](testHelper/headersMatchers.go)       | Provides functions to match HTTP headers with expected headers.           |
| [`StatusCodeMatchers`](testHelper/statusCodeMatchers.go) | Provides functions to match HTTP status codes with expected status codes. |

### Types
Package types provides utility types and functions.

| File Name                       | Description                                                  |
|---------------------------------|--------------------------------------------------------------|
| [`Optional`](types/optional.go) | Provides a wrapper to use any type as optional and nullable. |

### Utilities
The utilities package provides utility functions for making HTTP calls.

| File Name                                               | Description                                                                                     |
|---------------------------------------------------------|-------------------------------------------------------------------------------------------------|
| [`API Helper`](utilities/apiHelper.go)                  | Provides helper functions for making the HTTP calls.                                            |
| [`Marshal Error`](utilities/marshal_error.go)           | Defines a structure for error that will be returned if marshalling/unmarshalling fail.          |
| [`Union Type Helper`](utilities/unionTypeHelper.go)     | Provides helper functions for unmarshalling containers of multiple types.                       |
| [`Additional Properties`](utilities/additionalProperties.go) | Provides helper functions for handling additional properties. |

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
[coverage-badge]: https://sonarcloud.io/api/project_badges/measure?project=apimatic_go-core-runtime&metric=coverage
[coverage-url]: https://sonarcloud.io/summary/new_code?id=apimatic_go-core-runtime
[maintainability-badge]: https://sonarcloud.io/api/project_badges/measure?project=apimatic_go-core-runtime&metric=sqale_rating
[maintainability-url]: https://sonarcloud.io/summary/new_code?id=apimatic_go-core-runtime
[vulnerabilities-badge]: https://sonarcloud.io/api/project_badges/measure?project=apimatic_go-core-runtime&metric=vulnerabilities
[vulnerabilities-url]: https://sonarcloud.io/summary/new_code?id=apimatic_go-core-runtime
