package main

import (
	"encoding/json"
	"github.com/heroku/docker-registry-client/registry"
	"log"
	"net/http"
)

type Config struct {
	Url      string
	Username string
	Password string
}

const CONTENT_TYPE  = "application/json"

func (config *Config) Auth() (*registry.Registry, error) {
	hub, err := registry.New(config.Url, config.Username, config.Password)
	if err != nil {
		panic(err.Error())
	}

	return hub, nil
}

func (config *Config) GetRepositories(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s request to %s\n", r.Method, r.RequestURI)

	hub, err := config.Auth()
	if err != nil {
		panic(err)
	}

	repositories, err := hub.Repositories()
	if err != nil {
		panic(err)
	}

	response, err := json.Marshal(repositories)

	w.Header().Set("Content-Type", CONTENT_TYPE)
	w.Write(response)
}

func main() {
	config := &Config{
		Url:      "",
		Username: "",
		Password: "",
	}

	http.HandleFunc("/catalog", config.GetRepositories)
	log.Fatal("ListenAndServe: ", http.ListenAndServe(":8080", nil))
}
