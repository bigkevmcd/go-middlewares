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
