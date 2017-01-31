package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/heroku/docker-registry-client/registry"
	"log"
	"net/http"
	"sort"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
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
		log.Panic(err)
	}

	tags, err := hub.Tags(name)
	if err != nil {
		log.Panic(err)
	}

	sort.Strings(tags)
	response, err := json.Marshal(tags)

	w.Header().Set("Content-Type", CONTENT_TYPE)
	w.Write(response)
}

func Load(path string) (config *Config, err error) {
	var configFile Config

	f, err := os.Open(os.ExpandEnv(path))
	if err != nil {
		return
	}
	defer f.Close()

	d, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(d, &configFile)
	if err != nil {
		log.Panic(err.Error())
	}

	config = &configFile

	return config, nil
}

func main() {
	config, err := Load("config.yml")
	if err != nil {
		log.Panic(err)
	}

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/catalog", config.GetRepositories).Methods("GET")
	router.HandleFunc("/tags/{name}", config.GetTags).Methods("GET")
	log.Fatal("ListenAndServe: ", http.ListenAndServe(":8080", router))
}
