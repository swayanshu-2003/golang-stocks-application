package router

import (
	"stocks-api/middlewares"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/stock/{id}", middlewares.GetStock).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/stocks", middlewares.GetAllStocks).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/stock/create", middlewares.CreateStock).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/stock/{id}", middlewares.UpdateStock).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/stock/delete/{id}", middlewares.DeleteStock).Methods("DELETE", "OPTIONS")
	return router
}
