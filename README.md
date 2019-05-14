<a href="https://echo.labstack.com"><img height="80" src="https://cdn.labstack.com/images/echo-logo.svg"></a>

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/adverax/echo)
[![Go Report Card](https://goreportcard.com/badge/github.com/adverax/echo?style=flat-square)](https://goreportcard.com/report/github.com/adverax/echo)

## IMPORTANT
Declarative high performance idiomatic framework.

Documentation is under development. 
The following is the source file from labstack/adverax.

## Supported Go versions

As of version 4.0.0, Echo is available as a [Go module](https://github.com/golang/go/wiki/Modules).
Therefore a Go version capable of understanding /vN suffixed imports is required:

- 1.9.7+
- 1.10.3+
- 1.11+

Any of these versions will allow you to import Echo as `github.com/echo/echo/v5` which is the recommended
way of using Echo going forward.

For older versions, please use the latest v3 tag.

## Feature Overview

- Optimized HTTP router which smartly prioritize routes
- Build robust and scalable RESTful APIs
- Group APIs
- Extensible middleware framework
- Define middleware at root, group or route level
- Data binding for JSON, XML and form payload
- Handy functions to send variety of HTTP responses
- Centralized HTTP error handling
- Template rendering with any template engine
- Define your format for the logger
- Highly customizable
- Automatic TLS via Let’s Encrypt
- HTTP/2 support

## Benchmarks

Date: 2018/03/15<br>
Source: https://github.com/vishr/web-framework-benchmark<br>
Lower is better!

<img src="https://i.imgur.com/I32VdMJ.png">

### Example

```go
package main

import (
	"net/http"

	"github.com/adverax/echo"
	"github.com/adverax/echo/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	// Start server
	e.Logger.Error(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
```

## Contribute

**Use issues for everything**

- For a small change, just send a PR.
- For bigger changes open an issue for discussion before sending a PR.
- PR should have:
  - Test case
  - Documentation
  - Example (If it makes sense)
- You can also contribute by:
  - Reporting issues
  - Suggesting new features or enhancements
  - Improve/fix documentation

## License

[MIT](https://github.com/adverax/echo/blob/master/LICENSE)
