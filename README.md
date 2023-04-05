# Chi-Zerologformater

This package provides integrate [zerolog](https://github.com/rs/zerolog/) to [go-chi](https://github.com/go-chi/chi)
## Installation

You can install Chi-Zerologformater using
```bash
go get -u github.com/lastbyte32/chi-zerologformater
```

### Simple Example
```go
package main

import (
    "fmt"
    "net/http"
    "os"

    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    zeroformater "github.com/lastbyte32/chi-zerologformater"
)

func main() {
    // Create a new logger using Zerolog
    logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

    // Create a new HTTP request logger using ZeroFormater
    httpZeroLog := zeroformater.New(logger)

    // Create a new Chi router
    router := chi.NewRouter()

    // Add the ZeroFormater middleware to the router
    router.Use(middleware.RequestLogger(httpZeroLog))

    // Define a simple route
    router.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })

    // Start the HTTP server
	logger.Info().Msg("Listening on :3000...")
    http.ListenAndServe(":3000", router)
}

```
