package main

import (
    "github.com/joho/godotenv"
    "fmt"
    "time"
    "net/http"
    "encoding/json"
    "os"
)

type VersionResponse struct {
    Version string `json:"version"`
}

type AddressResponse struct {
    Address string `json:"address"`
}

func main() {
    if _, err := os.Stat(".env"); err == nil {
        err := godotenv.Load()
        must(err)
    }

    port, set := os.LookupEnv("IFTTT_PORT")
    if !set {
        port = "3693"
    }
    host, set := os.LookupEnv("IFTTT_HOST")
    if !set {
        host = "localhost"
    }
    fmt.Println("-------------------")
    fmt.Println("If Tangle Then That")
    fmt.Println("-------------------")
    fmt.Printf("Listening on http://%s:%s\n", host, port)
    http.HandleFunc("/", home)
    http.HandleFunc("/address", getAddress)
    http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil)
}

func home(w http.ResponseWriter, r *http.Request) {
    data := VersionResponse{"0.1"}
    js, err := json.Marshal(data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    sendJsonResponse(w, r, js)
}

func getAddress(w http.ResponseWriter, r *http.Request) {
    // Placeholder, should generate a new address for us!
    data := AddressResponse{"ABCDE9999999999999999999999999999999999999999999"}
    js, err := json.Marshal(data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    sendJsonResponse(w, r, js)
}

func sendJsonResponse(w http.ResponseWriter, r *http.Request, js []byte) {
    t := time.Now()
    fmt.Println(t.Format("[2006/01/02 15:04:05]"), r.RemoteAddr, "-",  r.Method, r.URL)
    w.Header().Set("Server", "If Tangle Then That - 0.1")
    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
}

func must(err error) {
    if err != nil {
        panic(err)
    }
}
