package main

import (
    "github.com/joho/godotenv"
    "github.com/iotaledger/iota.go/api"
    "github.com/iotaledger/iota.go/converter"
    "github.com/iotaledger/iota.go/bundle"
    "github.com/iotaledger/iota.go/account"
    "github.com/iotaledger/iota.go/account/deposit"
    "github.com/iotaledger/iota.go/account/event"
    "github.com/iotaledger/iota.go/account/event/listener"
    "github.com/iotaledger/iota.go/account/timesrc"
    "github.com/iotaledger/iota.go/account/builder"
    "github.com/iotaledger/iota.go/account/store/badger"
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
    Cda string `json:"cda"`
}

var acc account.Account

func main() {
    if _, err := os.Stat(".env"); err == nil {
        err := godotenv.Load()
        must(err)
    }

    apiSettings := api.HTTPClientSettings{URI: os.Getenv("IFTTT_NODE_URI")}
    iotaAPI, err := api.ComposeAPI(apiSettings)
    must(err)
    seed := os.Getenv("IFTTT_SEED")
    // Init the badger store for the account module
    store, err := badger.NewBadgerStore("db")
    defer store.Close()
    must(err)

	em := event.NewEventMachine()

    timesource := timesrc.NewNTPTimeSource("time.google.com")
	
	acc, err = builder.NewBuilder().
		WithAPI(iotaAPI).
		WithStore(store).
		WithSeed(seed).
		WithMWM(9).
		WithTimeSource(timesource).
		WithEvents(em).
		WithDefaultPlugins().
		Build()
	must(err)

	defer acc.Shutdown()
    
	must(acc.Start())

	lis := listener.NewCallbackEventListener(em)
	lis.RegReceivedDeposits(func(bun bundle.Bundle) {
		fmt.Println("Receiving Deposit!")
		for _, tx := range bun {
			msg, err := converter.TrytesToASCII(removeSuffixNine(tx.SignatureMessageFragment))
			if err == nil {
				fmt.Println("Message: ", msg)
			}
		}
	})
	
	balance, err := acc.AvailableBalance()
	must(err)
	fmt.Println("Balance:", balance, "i")

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
    sendJsonResponse(w, r, data)
}

func getAddress(w http.ResponseWriter, r *http.Request) {
    // Placeholder, should generate a new address for us!
    timesource := timesrc.NewNTPTimeSource("time.google.com")
	now, err := timesource.Time()
	must(err)
	fmt.Println(acc)
	now = now.Add(time.Duration(24) * time.Hour)
	conditions := &deposit.Conditions{TimeoutAt: &now, MultiUse: false}
	cda, err := acc.AllocateDepositAddress(conditions)
	must(err)
	link, err := cda.AsMagnetLink()
	must(err)
    data := AddressResponse{cda.Address, link}
    sendJsonResponse(w, r, data)
}

func sendJsonResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
    js, err := json.Marshal(data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
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

func removeSuffixNine(frag string) string {
    fraglen := len(frag)
    var firstNonNineAt int
    for i := fraglen - 1; i > 0; i-- {
         if frag[i] != '9' {
             firstNonNineAt = i
             break;
        }
    }
    return frag[:firstNonNineAt+1]
}
