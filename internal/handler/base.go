package handler

import "net/http"

func Root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("root"))
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func Panic(w http.ResponseWriter, r *http.Request) {
	panic("test")
}
