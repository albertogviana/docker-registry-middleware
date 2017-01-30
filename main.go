package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/heroku/docker-registry-client/registry"
	"log"
	"net/http"
	"sort"
)

type Config struct {
	Url      string
	Username string
	Password string
}

const CONTENT_TYPE = "application/json"

func (config *Config) Auth() (*registry.Registry, error) {
	hub, err := registry.New(config.Url, config.Username, config.Password)
	if err != nil {
		log.Panic(err)
	}

	return hub, nil
}

func (config *Config) GetRepositories(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s request to %s\n", r.Method, r.RequestURI)

	hub, err := config.Auth()
	if err != nil {
		log.Panic(err)
	}

	repositories, err := hub.Repositories()
	if err != nil {
		log.Panic(err)
	}

	response, err := json.Marshal(repositories)

	w.Header().Set("Content-Type", CONTENT_TYPE)
	w.Write(response)
}

func (config *Config) GetTags(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s request to %s \n", r.Method, r.RequestURI)

	vars := mux.Vars(r)
	name := vars["name"]

	hub, err := config.Auth()
	if err != nil {
		panic(err)
	}

	tags, err := hub.Tags(name)
	if err != nil {
		panic(err)
	}

	sort.Strings(tags)
	response, err := json.Marshal(tags)

	w.Header().Set("Content-Type", CONTENT_TYPE)
	w.Write(response)
}

func main() {
	config := &Config{
		Url:      "",
		Username: "",
		Password: "",
	}

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/catalog", config.GetRepositories)
	router.HandleFunc("/tags/{name}", config.GetTags).Methods("GET")
	log.Fatal("ListenAndServe: ", http.ListenAndServe(":8080", router))
}
