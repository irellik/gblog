package helpers

import (
	"io/ioutil"
	"log"
	"net/http"
)

func HttpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return "error", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}
