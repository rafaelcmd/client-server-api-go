package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cambio struct {
	USDBRL USDBRL `json:"USDBRL"`
}

type USDBRL struct {
	Code       string `json:"code"`
	CodeIn     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	TimeStamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type Response struct {
	Value string
}

func main() {
	http.HandleFunc("/cotacao", handler)
	port := ":8080"
	fmt.Printf("Server listening on port %s...\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr, Timeout: 10 * time.Second}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*2000))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("Error making request:", err)
	}
	defer res.Body.Close()

	_, err = io.WriteString(w, getBid(*res))
	if err != nil {
		fmt.Println("Error getting bid:", err)
	}
}

func getBid(response http.Response) string {
	decoder := json.NewDecoder(response.Body)
	var data Cambio
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding data:", err)
	}

	var responseValue Response

	db, err := sql.Open("sqlite3", "./database/cotacao.db")
	if err != nil {
		log.Fatal(err)
	}

	sqlCreateTable := `CREATE TABLE dolar(currentValue REAL);`

	_, err = db.Exec(sqlCreateTable)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlCreateTable)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*10))
	defer cancel()

	sqlInsertCurrentValue := `INSERT INTO dolar(currentValue) VALUES(?);`

	_, err = db.ExecContext(ctx, sqlInsertCurrentValue, data.USDBRL.Bid)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlInsertCurrentValue)
	}

	value, err := db.Query("SELECT * FROM dolar")
	if err != nil {
		log.Fatal(err)
	}
	for value.Next() {
		var dol string
		err = value.Scan(&dol)
		if err != nil {
			log.Fatal(err)
		}
		responseValue.Value = dol
	}

	return responseValue.Value
}
