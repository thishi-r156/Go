package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct{
	Id       string `json:"id" bson:"_id"`
	TaskName string `json:"task_name" bson:"task_name"`
	TaskDate string `json:"date" bson:"date"`
} 

var todoList = db().Database("GoProjects").Collection("TODOList")

func getTasks(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var tasks []Task
	cursor, err := todoList.Find(context.TODO(), bson.M{})
	if err != nil{
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return 
	}

	for cursor.Next(context.TODO()){
		var task Task
		err := cursor.Decode(&task)
		if err != nil {
			http.Error(w, "Decoding error", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	if len(tasks) == 0{
		w.Write([]byte("No tasks available"))
		return
	}else{
		json.NewEncoder(w).Encode(tasks)
	}
}

func getTask(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var task Task
	params := mux.Vars(r)
	objectID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	filter := bson.M{"_id": objectID}
	err = todoList.FindOne(context.TODO(), filter).Decode(&task)
	if err != nil{
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}else{
		json.NewEncoder(w).Encode(task)
	}
}

func createTask(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var task Task
	if task.Id == "" {
		task.Id = primitive.NewObjectID().Hex()
	}
	json.NewDecoder(r.Body).Decode(&task)
	_, err := todoList.InsertOne(context.TODO(), task)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}else{
		log.Println("Task added successfully!!")
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var task Task
	json.NewDecoder(r.Body).Decode(&task)
	params := mux.Vars(r)
	objectId, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id":objectId}
	_, err := todoList.DeleteOne(context.TODO(), filter)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateTask(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var task Task
	params := mux.Vars(r)
	json.NewDecoder(r.Body).Decode(&task)

	filter := bson.M{"_id":params["id"]}
	update := bson.M{"$set": bson.M{
		"task_name" :task.TaskName, 
		"date" : task.TaskDate,
	}}
	_, err := todoList.UpdateOne(context.TODO(), filter, update)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}else{
		log.Println("Update sucessful")
	}
}

 
func main(){
	r := mux.NewRouter()
	r.HandleFunc("/tasks", getTasks).Methods("GET")
	r.HandleFunc("/task/{id}", getTask).Methods("GET")
	r.HandleFunc("/task", createTask).Methods("POST")
	r.HandleFunc("/task/{id}", updateTask).Methods("PUT")
	r.HandleFunc("/task/{id}", deleteTask).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", r))
}