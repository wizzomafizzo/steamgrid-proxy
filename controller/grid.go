package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"wizzomafizzo/steamgrid-proxy/config"
	"wizzomafizzo/steamgrid-proxy/proxy"

	"github.com/gorilla/mux"
	"golang.org/x/text/unicode/norm"
)

func getFromCache(g string, s string) (string, error) {
	g = norm.NFC.String(g)

	data, err := os.ReadFile(filepath.Join(config.ProcessPath, "cache", s, g+".txt"))
	if err != nil {
		return "", err
	}

	link := string(data)

	client := &http.Client{}
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return "", err
	}

	r, err := client.Do(req)
	if err != nil || r.StatusCode != 200 {
		return "", err
	}

	return link, nil
}

func Search(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	searchTerm := vars["gameName"]
	searchType := vars["type"]

	if searchTerm == "" {
		w.WriteHeader(400)
		w.Write([]byte("game name is empty"))
		return
	}

	if searchType == "" {
		searchType = "grids"
	}

	if !config.IsValidImageType(searchType) {
		w.WriteHeader(400)
		w.Write([]byte("search type invalid, must be: grids, hgrids, heroes, logos, icons"))
		return
	}

	link, err := getFromCache(searchTerm, searchType)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("error while retrieving from cache"))
		return
	}

	if link != "" {
		jsonRes, err := json.Marshal(link)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("error while parsing cached data"))
			return
		} else {
			w.WriteHeader(200)
			w.Write(jsonRes)
			return
		}
	}

	res, err := proxy.Search(searchTerm, searchType)
	if len(res) == 0 {
		w.WriteHeader(404)
		w.Write([]byte("no results found"))
		return
	} else if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}

	jsonRes, _ := json.Marshal(res)
	w.WriteHeader(200)
	w.Write(jsonRes)
}
