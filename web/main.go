package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var sensorFolder *string

func main() {
	sensorFolder = flag.String("sensorfolder", "sensordata", "Folder to read data from")
	flag.Parse()
	r := mux.NewRouter()

	// Add the static folder (containing the transaction sending UI for both admin and user)
	r.HandleFunc("/sensors", Sensors).Methods("GET")
	r.HandleFunc("/readings/{sensorId}", Readings).Methods("GET")

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static"))))

	srv := &http.Server{
		Handler: cors.Default().Handler(r),
		Addr:    ":8001",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	panic(srv.ListenAndServe())
}

type Sensor struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func Sensors(w http.ResponseWriter, r *http.Request) {
	arr := make([]Sensor, 0)

	files, err := ioutil.ReadDir(*sensorFolder)
	if err != nil {
		writeError(w, err)
		return
	}

	for _, f := range files {
		if f.IsDir() {
			s := Sensor{}
			mdFile := filepath.Join(*sensorFolder, f.Name(), "metadata.json")
			if _, err := os.Stat(mdFile); !os.IsNotExist(err) {
				mdJson, err := ioutil.ReadFile(mdFile)
				if err == nil {
					err = json.Unmarshal(mdJson, &s)
					if err != nil {
						fmt.Println("Metadata file invalid: %s", err.Error())
					}
				}
			}
			s.Id = f.Name()
			arr = append(arr, s)
		}
	}

	writeJson(w, arr)
}

type Reading struct {
	Statement       string
	Base64Proof     string
	SensorTimestamp int64
}

func Readings(w http.ResponseWriter, r *http.Request) {
	arr := make([]Reading, 0)

	vars := mux.Vars(r)

	sensorPath := filepath.Join(*sensorFolder, vars["sensorId"])
	absSensorPath, err := filepath.Abs(sensorPath)
	if err != nil {
		writeError(w, fmt.Errorf("Sensor not found"))
		return
	}
	absSensorFolder, err := filepath.Abs(*sensorFolder)
	if err != nil {
		writeError(w, fmt.Errorf("Sensor not found"))
		return
	}

	if !strings.HasPrefix(absSensorPath, absSensorFolder) {
		writeError(w, fmt.Errorf("Sensor not found"))
		return
	}

	if _, err := os.Stat(absSensorPath); os.IsNotExist(err) {
		writeError(w, fmt.Errorf("Sensor not found"))
		return
	}

	files, err := ioutil.ReadDir(absSensorPath)
	if err != nil {
		writeError(w, err)
		return
	}

	for _, f := range files {
		if !f.IsDir() && f.Name() != "metadata.json" {
			fsFile := filepath.Join(sensorPath, f.Name())
			fsBytes, err := ioutil.ReadFile(fsFile)
			if err == nil {
				timestamp, _ := strconv.ParseInt(strings.Split(f.Name(), "-")[0], 10, 64)

				fs := ForeignStatementFromBytes(fsBytes)
				r := Reading{
					Base64Proof:     base64.StdEncoding.EncodeToString(fsBytes),
					Statement:       fs.StatementPreimage,
					SensorTimestamp: timestamp,
				}
				arr = append(arr, r)
			}
		}
	}
	writeJson(w, arr)
}

// writeError returns an error to the caller and sets the HTTP status
// to 500 (internal server error)
func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	fmt.Fprintf(w, "%s", err.Error())
}

func writeJson(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(v)
}
