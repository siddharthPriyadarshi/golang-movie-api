package main

import (
	"fmt"
	"github.com/siddharthPriyadarshi/mongoapi/router"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Welcome to Mongo API")
	r := router.Router()
	fmt.Println("Server is getting started...")

	fmt.Println("Server started on http://localhost:3000")

	log.Fatal(http.ListenAndServe(":3000", r))

}
