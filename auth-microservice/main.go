package main

import (
	"context"
	"log"
	"net/http"

	mw "github.com/cemgurhan/auth-microservice/middleware"

	usercontext "github.com/cemgurhan/auth-microservice/context"
)

func RequestID(ctx context.Context) string {
	requestID := ctx.Value(usercontext.UserContextKey)

	if requestID == nil {
		return "none"
	}

	return requestID.(string)

}

func main() { // try r.Group to share same context!!!
	// r := chi.NewRouter()
	// r.Use(mw.IsAuthorized)
	// r.Get("/", func(w http.ResponseWriter, req *http.Request) {

	// 	fmt.Println("SUCCESS")
	// 	fmt.Println(usercontext.UserFromContext(req.Context()))

	// },
	// )
	// // fmt.Println(RequestID(context.TODO()))

	// http.ListenAndServe(":9001", r)

	mux := http.NewServeMux()

	contextedMux := mw.IsAuthorized(mux)

	// mux.Handle("/", mw.IsAuthorized(http.HandlerFunc(finalHandler)))

	log.Fatal(http.ListenAndServe(":9001", contextedMux))

}

// func finalHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("handled"))
// }

// w.Header().Set("Content-Type", "application/json")

// requestURL := fmt.Sprintf("http://localhost:8081%v", req.RequestURI)
// var response *http.Response
// var err error
// if req.Method == "GET" {
// 	response, err = http.Get(requestURL)
// } else {
// 	response, err = http.Post(requestURL, "application/json", req.Body)
// }

// if err != nil {
// 	w.WriteHeader(http.StatusBadGateway)
// 	w.Write([]byte(fmt.Sprintf("%v", err)))
// }

// responseData, err := ioutil.ReadAll(response.Body)
// if err != nil {
// 	log.Fatal(err)
// }

// w.Write(responseData)
// w.WriteHeader(http.StatusOK)
// fmt.Println(usercontext.UserFromContext())
