package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
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
	http.HandleFunc("/cotacao/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(req.Context(), time.Duration(time.Millisecond*200))
	defer cancel()
	req = req.WithContext(ctx)
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	_, err = io.WriteString(w, getBid(*res))
	if err != nil {
		log.Fatal(err)
	}
}

func getBid(response http.Response) string {
	decoder := json.NewDecoder(response.Body)
	var data Cambio
	decoder.Decode(&data)

	var responseValue Response

	os.Remove("../database/cotacao.db")

	db, err := sql.Open("sqlite3", "../database/cotacao.db")
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
