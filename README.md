# go-middlewares

Some standard Middlewares for Go projects

## Logging middleware

```go
package main

import (
	"fmt"
	"html"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/bigkevmcd/go-middlewares/logging"
	"github.com/go-logr/logr"
)

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func main() {
	logger := logr.FromSlogHandler(slog.NewTextHandler(os.Stdout, nil))
	handler := logging.RequestLogger(logger)(http.HandlerFunc(helloWorldHandler))

	http.Handle("/", handler)

	log.Println("listening on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

Running this and using curl to hit "http://localhost:8080/testing"

```
time=2025-04-06T18:03:56.475+01:00 level=INFO msg=OK method=GET host=localhost:8080 path=/testing remote_addr=[::1]:51360 user_agent=curl/8.7.1 protocol="" duration=0ms
```
