package main

import (
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	response, err := http.Get("http://localhost:8080/cotacao")

	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	cotacao := []byte("Dólar: " + string(body[:]))

	err = os.WriteFile("../cotacao", cotacao, 0644)
	if err != nil {
		panic(err)
	}
}
