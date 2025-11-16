package router

import (
	"Infocenter/internal/api/handlers"
	"net/http"
)

func MainRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /infocenter/{topic}", handlers.InfocenterGETHandler)
	mux.HandleFunc("POST /infocenter/{topic}", handlers.InfocenterPOSTHandler)

	return mux
}
