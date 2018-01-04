package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DataDog/dd-trace-go/tracer"
	"github.com/DataDog/dd-trace-go/tracer/contrib/gorilla/muxtrace"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	// muxTracer := muxtrace.NewMuxTracer("goapp", tracer.DefaultTracer)
	// muxTransport := tracer.NewTransport("localhost", "8126") // Default
	// where docker.for.mac.localhost is the host's IP
	muxTransport := tracer.NewTransport("docker.for.mac.localhost", "8126")
	muxTracer := muxtrace.NewMuxTracer("goapp", tracer.NewTracerTransport(muxTransport))
	muxTracer.HandleFunc(r, "/", hello)
	r.HandleFunc("/", hello)
	http.Handle("/", r)
	fmt.Println("Starting up on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// HandleFunc for '/'
func hello(w http.ResponseWriter, req *http.Request) {
	span := tracer.SpanFromContextDefault(req.Context())
	fmt.Printf("tracing service:%s resource:%s", span.Service, span.Resource)
	fmt.Fprintln(w, "Hello world!")
	// w.Write([]byte("hello world"))
}
