package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"wizzomafizzo/steamgrid-proxy/config"

	"golang.org/x/text/unicode/norm"
)

type ProxySearchResponse struct {
	Success bool       `json:"success"`
	Data    []GameData `json:"data"`
}

type GameData struct {
	Types       []string `json:"types"`
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Verified    bool     `json:"verified"`
	ReleaseDate string   `json:"release_date,omitempty"`
}

type ProxyGridResponse struct {
	Success bool       `json:"success"`
	Data    []GridData `json:"data"`
}

type GridData struct {
	Url   string `json:"url"`
	Thumb string `json:"thumb"`
}

type Result struct {
	GameName     string `json:"gameName"`
	ImageUrl     string `json:"imageUrl"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

const BASE_URL = "https://www.steamgriddb.com/api/v2"

const GRID_DIMENSIONS = "dimensions=600x900"

// const HGRID_DIMENSIONS = "dimensions=920x430"
// const HERO_DIMENSIONS = "dimensions=1920x620"

func callAPI(e string, t string, p string) (*http.Response, error) {
	cnf := *config.Cnf
	client := &http.Client{}

	req, err := http.NewRequest("GET", BASE_URL+e+t+"?"+p, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+cnf.ApiKey)
	req.Header.Set("User-Agent", "curl/7.79.1")

	r, err := client.Do(req)

	return r, err
}

func AutocompleteSearch(t string) (searchResponse ProxySearchResponse, err error) {
	res, err := callAPI("/search/autocomplete/", t, "")
	if err != nil {
		return ProxySearchResponse{}, err
	}

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return ProxySearchResponse{}, err
	}

	err = json.Unmarshal(buf, &searchResponse)
	if err != nil {
		return ProxySearchResponse{}, err
	}

	if len(searchResponse.Data) == 0 {
		return ProxySearchResponse{}, errors.New("no results found")
	}

	return searchResponse, nil
}

func Search(t string, s string) ([]Result, error) {
	results := []Result{}

	searchResponse, err := AutocompleteSearch(t)
	if err != nil {
		return results, err
	}

	dimensions := GRID_DIMENSIONS

	var itype string = s

	res, err := callAPI(
		fmt.Sprintf("/%s/game/", itype),
		fmt.Sprint(searchResponse.Data[0].Id),
		dimensions,
	)
	if err != nil {
		return results, err
	}

	var gridResponse ProxyGridResponse

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return results, err
	}

	err = json.Unmarshal(buf, &gridResponse)
	if err != nil {
		return results, err
	}

	if len(gridResponse.Data) == 0 {
		return results, nil
	}

	for _, v := range gridResponse.Data {
		results = append(results, Result{
			GameName:     searchResponse.Data[0].Name,
			ImageUrl:     v.Url,
			ThumbnailUrl: v.Thumb,
		})
	}

	err = CreateCache(t, s, results)
	if err != nil {
		return results, err
	}

	return results, nil
}

func CreateCache(t string, s string, results []Result) (err error) {
	t = norm.NFC.String(t)

	_, err = os.Create(filepath.Join(config.ProcessPath, "cache", s, t+".txt"))
	if err != nil {
		return err
	}

	msg, err := json.Marshal(results)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(config.ProcessPath, "cache", s, t+".txt"), []byte(msg), 0)
	if err != nil {
		return err
	}

	return nil
}
