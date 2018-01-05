package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

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
	// r.HandleFunc("/", hello) // not required since we use muxTracer
	http.Handle("/", r)
	fmt.Println("Starting up on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// HandleFunc for '/'
func hello(w http.ResponseWriter, req *http.Request) {
	span := tracer.SpanFromContextDefault(req.Context())
	fmt.Printf("tracing service:%s resource:%s\n", span.Service, span.Resource)
	fmt.Fprintln(w, "Hello world!")
	i := rand.Intn(10)
	nestedHelloWithRandomSleep(w, req, i)
}

func nestedHelloWithRandomSleep(w http.ResponseWriter, req *http.Request, level int) int {
	if level <= 0 {
		return 0
	}
	span := tracer.NewChildSpanFromContext("mux.request.nested-hello", req.Context())
	r := rand.Intn(10) * 10
	time.Sleep(time.Duration(r) * time.Microsecond)
	fmt.Fprintln(w, "Nested Hello at level:", level, ". Slept for:", r)
	defer span.Finish()
	return nestedHelloWithRandomSleep(w, req, level-1)
}
