package main

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"sync/atomic"
)

var hit_counter atomic.Uint64

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        hit_count := hit_counter.Add(1)
        log.Printf("hit=%v\n", hit_count)
        writeGzCompressedResponse(w)
        log.Println("delivered")
    })
    mux.HandleFunc("/favicon.ico/", func(w http.ResponseWriter, r *http.Request) {
        // TODO: add naughty ico compression?
    })

    log.Println("serving on :8080")
    err := http.ListenAndServe("localhost:8080", mux)
    if err != nil {
        log.Fatalf("server error: %s", err)
    }
}

type LowEntropyWriter struct {
    w io.Writer
}

func (w *LowEntropyWriter) Write(n uint64) (int, error) {
    p:= []byte("<!DOCTYPE html><html><head><title>Crazy? I was crazy once. They put me in a room. A rubber room. A rubber room with rats. The rats made me crazy. Crazy? I was crazy once...</title></head><body>")
    post := []byte("<p>Crazy?</p></body></html>")
    body := make([]byte, n)
    p = append(p, body...)
    p = append(p, post...)
    return w.w.Write(p)
}

func writeGzCompressedResponse(w http.ResponseWriter) {
    w.Header().Set("Content-Type", "text/html") // TODO: vary
    w.Header().Set("Content-Encoding", "gzip")

    gw, err1 := gzip.NewWriterLevel(w, 9)
    if err1 != nil { log.Fatalln("c") }
    defer gw.Close()
    lw := LowEntropyWriter{w: gw}
    _, err := lw.Write(1_000_000_000)
    if err != nil {
        log.Fatalln("failed to write low entropy")
    }
}
