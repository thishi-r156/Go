package main

import(
		"net/http"
		"log"
		"time")

type apiHandler struct{}

func (apiHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path{
	case "/api/index": 
		http.ServeFile(w, req, "index.html")
		w.WriteHeader(http.StatusOK)
	case "/api/profile":
		http.ServeFile(w, req, "profile.html")
	default:
		w.Write([]byte("404 Page not found"))
	}
	
}
		

func main(){
	address :=  ":8080"
	mux := http.NewServeMux()
	mux.Handle("/api/", apiHandler{})

	s := &http.Server{
		Addr:           address,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Listening on port: %v", address)
	log.Fatal(s.ListenAndServe())
}
