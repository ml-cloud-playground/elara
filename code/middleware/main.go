package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var storage Storage

func main() {
	user := os.Getenv("elara_user")
	pass := os.Getenv("elara_pass")
	host := os.Getenv("elara_host")
	name := os.Getenv("elara_name")
	redisHost := os.Getenv("REDISHOST")
	redisPort := os.Getenv("REDISPORT")
	port := os.Getenv("PORT")

	fmt.Printf("Port: %s\n", port)

	if err := storage.Init(user, pass, host, name, redisHost, redisPort, true); err != nil {
		panic(err)
	}
	defer storage.sqlstorage.Close()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/healthz", healthHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/healthz", healthHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/match", matchHandler).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/score/{input}", scoreHandler).Methods(http.MethodGet, http.MethodOptions)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

// CORSRouterDecorator applies CORS headers to a mux.Router
type CORSRouterDecorator struct {
	R *mux.Router
}

// ServeHTTP wraps the HTTP server enabling CORS headers.
// For more info about CORS, visit https://www.w3.org/TR/cors/
func (c *CORSRouterDecorator) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, Content-Type, YourOwnHeader")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}

	c.R.ServeHTTP(rw, req)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
	return
}

func matchHandler(w http.ResponseWriter, r *http.Request) {
	subcategory := mux.Vars(r)["subcategory"]

	t, err := storage.Match(subcategory)
	if err != nil {

		if strings.Contains(err.Error(), "Rows are closed") {
			return
		}

		writeErrorMsg(w, err)
		return
	}

	writeJSON(w, t, http.StatusOK)
}

func scoreHandler(w http.ResponseWriter, r *http.Request) {
	subcategory := mux.Vars(r)["input"]
	body := []byte(`{
		"instances": [
		  { "name": "","description":"` + subcategory + `"}
		]
	  }`)

	bodyReader := bytes.NewReader(body)
	access_token := os.Getenv("elara_ml_access_token")
	url := os.Getenv("elara_ml_url")
	bearer := "Bearer " + access_token
	req, err := http.NewRequest("POST", url, bodyReader)
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	score, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer score.Body.Close()

	bodyBytes, err := io.ReadAll(score.Body)
	if err != nil {
		log.Fatal(err)
	}

	var scoreStat = &ScoreStat{}
	err = json.Unmarshal(bodyBytes, scoreStat)

	if err != nil {
		log.Fatal(err)
	}

	scores := scoreStat.Predictions[0].Scores

	// Initialize the maximum value to the first element of the array.
	max := scores[0]
	maxIndex := 0
	// Iterate through the array and update the maximum value if the current element is greater.
	for index, value := range scores {
		if value > max {
			max = value
			maxIndex = index
		}
	}

	subCategory := scoreStat.Predictions[0].Classes[maxIndex]

	t, err := storage.Match(subCategory)

	if err != nil {

		if strings.Contains(err.Error(), "Rows are closed") {
			return
		}

		writeErrorMsg(w, err)
		return
	}

	writeJSON(w, t, http.StatusOK)
	return
}

// JSONProducer is an interface that spits out a JSON string version of itself
type JSONProducer interface {
	JSON() (string, error)
	JSONBytes() ([]byte, error)
}

func writeJSON(w http.ResponseWriter, j JSONProducer, status int) {
	json, err := j.JSON()
	if err != nil {
		writeErrorMsg(w, err)
		return
	}
	writeResponse(w, status, json)
	return
}

func writeErrorMsg(w http.ResponseWriter, err error) {
	s := fmt.Sprintf("{\"error\":\"%s\"}", err)
	writeResponse(w, http.StatusInternalServerError, s)
	return
}

func writeResponse(w http.ResponseWriter, status int, msg string) {
	if status != http.StatusOK {
		weblog(fmt.Sprintf(msg))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	w.WriteHeader(status)
	w.Write([]byte(msg))

	return
}

func weblog(msg string) {
	log.Printf("Webserver : %s", msg)
}

// Message is a structure for communicating additional data to API consumer.
type Message struct {
	Text    string `json:"text"`
	Details string `json:"details"`
}

// JSON marshalls the content of a elara to json.
func (m Message) JSON() (string, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("could not marshal json for response: %s", err)
	}

	return string(bytes), nil
}

// JSONBytes marshalls the content of a elara to json as a byte array.
func (m Message) n() ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return []byte{}, fmt.Errorf("could not marshal json for response: %s", err)
	}

	return bytes, nil
}
