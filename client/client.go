package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 10 * time.Second}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*300))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("Error making request:", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	cotacao := []byte("Dólar: " + string(body[:]))

	err = os.WriteFile("../cotacao", cotacao, 0666)
	if err != nil {
		panic(err)
	}
}
