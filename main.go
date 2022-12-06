package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/ncostamagna/gocourse_product/internal/product"
	"github.com/ncostamagna/gocourse_product/pkg/bootstrap"
)

func main() {

	router := mux.NewRouter()
	_ = godotenv.Load()
	l := bootstrap.InitLogger()

	db, err := bootstrap.DBConnection()
	if err != nil {
		l.Fatal(err)
	}

	productRepo := product.NewRepo(l, db)
	productSrv := product.NewService(l, productRepo)
	productEnd := product.MakeEndpoints(productSrv)

	router.HandleFunc("/products", productEnd.Create).Methods("POST")
	router.HandleFunc("/products/{id}", productEnd.Get).Methods("GET")
	router.HandleFunc("/products", productEnd.GetAll).Methods("GET")
	router.HandleFunc("/products/{id}", productEnd.Update).Methods("PATCH")
	router.HandleFunc("/products/{id}", productEnd.Delete).Methods("DELETE")

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
