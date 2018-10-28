package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var apiKey = "sqlrestTestKey"
var sqlRestURL = "http://localhost:5050"

var client = http.Client{}

type connectResponse struct {
	Message string `json:"message"`
}

type sqlRestQueryRequest struct {
	Query string `json:"query"`
}

type sqlRestProcRequest struct {
	Name        string            `json:"name"`
	Parameters  map[string]string `json:"parameters"`
	ExecuteOnly bool              `json:"executeOnly"`
}

type sqlRestResponse struct {
	Message string     `json:"message"`
	Error   string     `json:"error"`
	Columns []string   `json:"Columns"`
	Data    [][]string `json:"Data"`
}

func main() {
	acceptedInputs := "Accepted inputs are: 'ping', 'connect', 'query', and 'procedure'"
	f := ""

	if len(os.Args) == 2 {
		f = os.Args[1]
	}

	if f == "" {
		fmt.Println("You must provide a function to execute")
		fmt.Println(acceptedInputs)
		return
	}

	switch f {
	case "ping":
		ping()
		break
	case "connect":
		connect()
		break
	case "query":
		query()
		break
	case "procedure":
		procedure()
		break
	default:
		fmt.Println("Unknown input")
		fmt.Println(acceptedInputs)
	}
}

func ping() {
	url := sqlRestURL + "/ping"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	setAuthHeader([]byte(""), req)

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	fmt.Println(string(body))
}

func connect() {
	url := sqlRestURL + "/connect"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	setAuthHeader([]byte(""), req)

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	connectResponse := connectResponse{}
	jsonErr := json.Unmarshal(body, &connectResponse)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Println(connectResponse.Message)
}

func query() {
	url := sqlRestURL + "/v1/query"

	q := sqlRestQueryRequest{"SELECT TOP 3 * FROM Flights.dbo.Airlines"}
	jsonQuery, _ := json.Marshal(q)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonQuery))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	setAuthHeader(jsonQuery, req)

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	queryResponse := sqlRestResponse{}
	jsonErr := json.Unmarshal(body, &queryResponse)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Println(queryResponse)
}

func procedure() {
	url := sqlRestURL + "/v1/procedure"

	p := sqlRestProcRequest{Name: "Flights.dbo.AirportsByAirline",
		Parameters:  map[string]string{"airlineId": "109"},
		ExecuteOnly: false}

	procRequest, _ := json.Marshal(p)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(procRequest))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	setAuthHeader(procRequest, req)

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	queryResponse := sqlRestResponse{}
	jsonErr := json.Unmarshal(body, &queryResponse)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Println(queryResponse)
}

func setAuthHeader(msgBody []byte, request *http.Request) {
	realm := "testing-func"

	sqlRestHmac := ""

	if len(msgBody) > 0 {
		mac := hmac.New(sha256.New, []byte(apiKey))
		mac.Write(msgBody)
		sqlRestHmac = hex.EncodeToString(mac.Sum(nil))
	}

	nonce := "randomness"
	currentTime := time.Now().UnixNano() / 1000000
	timestamp := strconv.FormatInt(currentTime, 10)

	request.Header.Set("Authorization", strings.Join([]string{realm, sqlRestHmac, nonce, timestamp}, ":"))
}
