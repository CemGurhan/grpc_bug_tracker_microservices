package main

import (
	"log"
	"net/http"
)

// func Success(w http.ResponseWriter, r *http.Request) {

// 	fmt.Print("SUCCESS")

// }

func main() {

	// http.Handle("/", mw.IsAuthorized(Success))

	log.Fatal(http.ListenAndServe(":9001", nil))

}
