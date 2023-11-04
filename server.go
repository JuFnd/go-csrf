package main

import "net/http"

type Server struct {
	mx http.ServeMux
}
