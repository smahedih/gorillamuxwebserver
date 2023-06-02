package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gorilla/mux"
)

// Job Struct (Model)
type Job struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Init jobs var as a slice Job struct
var jobs []Job

// Get All Jobs
func getJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

// Get Single Job
func getJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) //get params
	//Loop through jobs and find with id
	for _, item := range jobs {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Job{})
}

// Create Job Post
func createJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var job Job
	_ = json.NewDecoder(r.Body).Decode(&job)
	job.ID = strconv.Itoa(rand.Intn(10000000)) //Mock ID - Not safe
	jobs = append(jobs, job)
	json.NewEncoder(w).Encode(&job)
}

// Update Job Post
func updateJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range jobs {
		if item.ID == params["id"] {
			jobs = append(jobs[:index], jobs[index+1:]...)
			var job Job
			_ = json.NewDecoder(r.Body).Decode(&job)
			job.ID = params["id"]
			jobs = append(jobs, job)
			json.NewEncoder(w).Encode(&job)
			return
		}
	}
	json.NewEncoder(w).Encode(&jobs)
}

// Delete Job Post
func deleteJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range jobs {
		if item.ID == params["id"] {
			jobs = append(jobs[:index], jobs[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(&jobs)
}

func main() {

	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

	//Init Router
	r := mux.NewRouter()

	//Mock Data
	jobs = append(jobs, Job{ID: "1", Title: "Job One"})
	jobs = append(jobs, Job{ID: "2", Title: "Job two"})

	// Route Handlers / Endpoints
	r.HandleFunc("/api/jobs", getJobs).Methods("GET")
	r.HandleFunc("/api/jobs/{id}", getJob).Methods("GET")
	r.HandleFunc("/api/jobs", createJob).Methods("PUT")
	r.HandleFunc("/api/jobs/{id}", updateJob).Methods("PATCH")
	r.HandleFunc("/api/jobs/{id}", deleteJob).Methods("DELETE")

	//http.ListenAndServe(":8080", r)

	go func() {

		if err := http.ListenAndServe(":8080", r); err != nil && err != http.ErrServerClosed {

			logger.Fatalf("cannot listen on defined port: %s\n", err)

		}
	}()

	logger.Println("This is an info message.")
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
}
