# Package `errs`

Package `errs` provides error codes and error handling functionalities.

## Overview

The `errs` package defines common error codes and provides a way to handle and generate error responses. It includes functionalities for creating error objects with error codes, messages, timestamps, and additional information. The package also maps the error codes to their corresponding HTTP status codes for convenient handling in web applications.

## Installation

To use the `errs` package, you can import it in your Go project:

```go
import "github.com/thirathawat/errs"
```

Make sure to run `go get` to retrieve the package:

```shell
go get github.com/thirathawat/errs
```

## Usage

### Error Codes

The package defines the following error codes:

- `CodeBadRequest`: Represents a bad request error.
- `CodeUnauthorized`: Represents an unauthorized error.
- `CodeForbidden`: Represents a forbidden error.
- `CodeNotFound`: Represents a not found error.
- `CodeGone`: Represents a gone error.
- `CodeTooManyRequests`: Represents a too many requests error.
- `CodeInternalServerError`: Represents an internal server error.
- `CodeNotImplemented`: Represents a not implemented error.
- `CodeServiceUnavailable`: Represents a service unavailable error.

### Creating Errors

To create a new error, use the `New` function provided by the package:

```go
err := errs.New(errs.CodeBadRequest, "Invalid request")
```

You can also provide additional options when creating an error. For example, you can include additional information or log the error:

```go
err := errs.New(errs.CodeInternalServerError, "Internal server error",
    errs.WithInfo(map[string]interface{}{"requestID": "abc123"}),
    errs.WithLogErr(innerError),
)
```

### Handling Errors

The package provides a convenient function `ResponseError` to handle errors in a Gin HTTP handler:

```go
func MyHandler(c *gin.Context) {
    err := // Some operation that may return an error
    if err != nil {
        errs.ResponseError(c, err)
        return
    }

    // Handle successful case
    ...
}
```

The `ResponseError` function checks if the provided error is an `errs.Error` object. If it is, it returns a JSON response with the error information, including the error code and message. Otherwise, it returns a generic internal server error response.

### Validation Errors

The package includes functionality to handle validation errors. If you have a validation error returned by a validation library, you can convert it to an `errs.Error` object using the `InvalidStructError` function:

```go
validationErr := // Some validation error
err := errs.InvalidStructError(validationErr)
```

This function converts the validation error into a structured error with the appropriate error code, message, and additional information about the validation errors.

## License

This package is licensed under the MIT License. See the LICENSE file for more information.
