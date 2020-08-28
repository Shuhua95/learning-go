package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type event struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type events []event

var evts = events{
	{
		ID:          "1",
		Title:       "introduction to golang",
		Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
	},
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}

func createEvt(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	var evt event
	if err := json.Unmarshal(body, &evt); err != nil {
		fmt.Println("error:", err)
	}

	evts = append(evts, evt)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(evt)
}

func getEvt(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	for _, evt := range evts {
		if evt.ID == id {
			json.NewEncoder(w).Encode(evt)

			break
		}
	}
}

func getAllEvts(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(evts)
}

func updateEvt(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	var newEvt event
	json.Unmarshal(body, &newEvt)

	for i, evt := range evts {
		if evt.ID == id {
			evt.Title = newEvt.Title
			evt.Description = newEvt.Description
			evts = append(evts[:i], evt)

			json.NewEncoder(w).Encode(evt)

			break
		}
	}
}

func delEvt(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	for i, evt := range evts {
		if evt.ID == id {
			evts = append(evts[:i], evts[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", id)

			break
		}
	}
}

func main() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", index)

	r.HandleFunc("/event", createEvt).Methods("POST")
	r.HandleFunc("/event", getAllEvts).Methods("GET")
	r.HandleFunc("/event/{id}", getEvt).Methods("GET")
	r.HandleFunc("/event/{id}", updateEvt).Methods("PATCH")
	r.HandleFunc("/event/{id}", delEvt).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
