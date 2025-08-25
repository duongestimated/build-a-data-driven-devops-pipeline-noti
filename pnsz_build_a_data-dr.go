package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type notification struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type config struct {
	WebhookURL string `json:"webhook_url"`
}

var cfg config

func init() {
	loadConfig()
}

func loadConfig() {
	jsonFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(jsonFile, &cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func sendNotification(message, status string) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", cfg.WebhookURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	noti := notification{Message: message, Status: status}
	jsonData, err := json.Marshal(noti)
	if err != nil {
		log.Fatal(err)
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(jsonData))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func handleBuildEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	buildID := vars["build_id"]
	buildStatus := vars["build_status"]

	sendNotification(fmt.Sprintf("Build %s has finished with status %s", buildID, buildStatus), buildStatus)

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Notification sent successfully"))
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/build/{build_id}/{build_status}", handleBuildEvent).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}