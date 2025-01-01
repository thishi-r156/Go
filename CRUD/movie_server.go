package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"math/rand"
	"github.com/gorilla/mux"
	"slices"
)

type Movie struct{
	ID string `json:id`
	Title string `json:title`
}

func getMovies(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	for _, item := range movies{
		if params["id"] == item.ID{
			json.NewEncoder(w).Encode(item)
			return
		}
	}

}

func createMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(100))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movies)
}

func updateMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	params := mux.Vars(r)
	for index, item := range movies{
		if params["id"] == item.ID{
			movies[index].Title = movie.Title
			json.NewEncoder(w).Encode(movies)
			return
		}
	}

	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	params := mux.Vars(r)
	for index, item := range movies{
		if params["id"] == item.ID{
			movies = slices.Delete(movies, index, index+1)
			return 
		}
	}
	json.NewEncoder(w).Encode(movies)
}

var movies []Movie

func main(){
	movies = append(movies, Movie{ID: "123", Title: "Movie 1"})
	movies = append(movies, Movie{ID: "456", Title: "Movie 2"})
	movies = append(movies, Movie{ID: "789", Title: "Movie 3"})
	r := mux.NewRouter()
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies/", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")
	fmt.Println("Starting on server: 8000")
	log.Fatal(http.ListenAndServe(":8000", r))

}