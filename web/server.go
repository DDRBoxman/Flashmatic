package web

import (
	"github.com/gorilla/mux"
	"net/http"
		"log"
	"github.com/DDRBoxman/Flashmatic/display"
	"encoding/json"
)

type keyEvent struct {
	KeyID int `json:"key_id"`
}

func StartServer(display *display.Display, keychan chan int) {
	router := mux.NewRouter()

	router.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("/key", func(w http.ResponseWriter, r *http.Request) {
		keyEvent := keyEvent{}
		err := json.NewDecoder(r.Body).Decode(&keyEvent)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		keychan <- keyEvent.KeyID
	}).Methods("POST")

	router.HandleFunc("/display", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(display.Icons)
	}).Methods("GET")

	router.PathPrefix("/icon/").Handler(http.StripPrefix("/icon/", http.FileServer(http.Dir(display.IconDir))))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Println(err)
	}
}