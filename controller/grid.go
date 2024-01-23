package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"wizzomafizzo/steamgrid-proxy/config"
	"wizzomafizzo/steamgrid-proxy/proxy"

	"github.com/gorilla/mux"
	"golang.org/x/text/unicode/norm"
)

func getFromCache(g string, s string) ([]proxy.Result, error) {
	results := []proxy.Result{}
	g = norm.NFC.String(g)

	data, err := os.ReadFile(filepath.Join(config.ProcessPath, "cache", s, g+".txt"))
	if err != nil {
		return results, nil
	}

	err = json.Unmarshal(data, &results)
	if err != nil {
		return results, nil
	}

	// link := string(data)

	// client := &http.Client{}
	// req, err := http.NewRequest("GET", link, nil)
	// if err != nil {
	// 	return "", err
	// }

	// r, err := client.Do(req)
	// if err != nil || r.StatusCode != 200 {
	// 	return "", err
	// }

	return results, nil
}

var gridUrlRegex = regexp.MustCompile("^cdn[0-9]{1,2}.steamgriddb.com/grid/")

// https://cdn2.steamgriddb.com/grid/5fc4a6bba793371c716812a0505c72e1.png
func ImageProxy(w http.ResponseWriter, r *http.Request) {
	url := mux.Vars(r)["url"]

	if url == "" {
		w.WriteHeader(400)
		w.Write([]byte("url is empty"))
		return
	}

	if !gridUrlRegex.MatchString(url) {
		w.WriteHeader(403)
		w.Write([]byte("url is invalid"))
		return
	}

	res, err := http.Get("https://" + url)
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}

	if res.StatusCode != 200 {
		w.WriteHeader(res.StatusCode)
		w.Write([]byte("error while retrieving image"))
		return
	}

	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.WriteHeader(200)
	io.Copy(w, res.Body)
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

	cached, err := getFromCache(searchTerm, searchType)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("error while retrieving from cache"))
		fmt.Println(err)
		return
	}

	if len(cached) > 0 {
		jsonRes, err := json.Marshal(cached)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("error while parsing cached data"))
			fmt.Println(err)
			return
		} else {
			w.WriteHeader(200)
			w.Write(jsonRes)
			return
		}
	}

	res, err := proxy.Search(searchTerm, searchType)
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	} else if len(res) == 0 {
		w.WriteHeader(404)
		w.Write([]byte("no results found"))
		return
	}

	jsonRes, _ := json.Marshal(res)
	w.WriteHeader(200)
	w.Write(jsonRes)
}
