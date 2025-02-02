package main 

import ("net/http"
		"log"
		"github.com/gorilla/mux"
		)


func welcomePage(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Welcome😊"))
}

func profilePage(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("These are the profile details"))
}

func main(){
	r := mux.NewRouter()
	r.HandleFunc("/", welcomePage)
	r. HandleFunc("/profile", profilePage)
	log.Fatal(http.ListenAndServe(":8080", r))
	
}