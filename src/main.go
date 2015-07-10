package main

import (
	"fmt"
	"log"
	"encoding/json"
	"golang.org/x/net/websocket"
	"github.com/goinggo/tracelog"
)

var TAG = "main"

// Represents a transaction we will write to the graph database.
type Transaction struct {
	hash string
	time int32
	// key: address, value: value
	inputs map[string]int64
	outputs map[string]int64
}

func SatoshisToBTC(value int64) float64 {
	return float64(value) / 100000000
}

func (t Transaction) String() string {
	s := "\n\n"
	s += fmt.Sprintf("Transaction Hash: %s\n", t.hash)
	s += fmt.Sprintf("Time: %v\n", t.time)
	s += " Inputs:\n"
	for key, value := range t.inputs {
		s += fmt.Sprintf("   %s : %v BTC\n", key, SatoshisToBTC(value))
	}
	s += " Outputs:\n"
	for key, value := range t.outputs {
		s += fmt.Sprintf("   %s : %v BTC\n", key, SatoshisToBTC(value))
	}
	s += "\n"
	return s
}

// This is pretty hacky, but it works.
// Look to optimize this in the future, especially when we need more data.
func TransactionFromJSON(b []byte) Transaction {	
	//fmt.Println(string(b))
	
	var objmap map[string]*json.RawMessage
	json.Unmarshal(b, &objmap)
	
	var x map[string]*json.RawMessage
	json.Unmarshal(*objmap["x"], &x)
	
	var hash string
	json.Unmarshal(*x["hash"], &hash)
	
	var time int32
	json.Unmarshal(*x["time"], &time)
	
	var out []*json.RawMessage
	json.Unmarshal(*x["out"], &out)
	
	outputMap := make(map[string]int64)
	for i := range out {
		var outMap map[string]*json.RawMessage
		json.Unmarshal(*out[i], &outMap)
	
		var addr string
		json.Unmarshal(*outMap["addr"], &addr)
		
		var value int64
		json.Unmarshal(*outMap["value"], &value)
		
		outputMap[addr] += value
	}
	
	var inputs []*json.RawMessage
	json.Unmarshal(*x["inputs"], &inputs)
	
	inputMap := make(map[string]int64)
	for i := range inputs {
		var inMap map[string]*json.RawMessage
		json.Unmarshal(*inputs[i], &inMap)
		
		var prev_out map[string]*json.RawMessage
		json.Unmarshal(*inMap["prev_out"], &prev_out)
		
		var addr string
		json.Unmarshal(*prev_out["addr"], &addr)
		
		var value int64
		json.Unmarshal(*prev_out["value"], &value)
		
		inputMap[addr] += value
	}
	
	return Transaction{hash, time, inputMap, outputMap} //, inputs, outputs}
}

// Checks if bytes are valid json.
// Probably a better way to this.
func IsValidJson(b []byte) bool {

	var objmap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objmap)
	
	return err == nil
}

// Handles a transaction 
func HandleTransaction(b []byte) {
	tracelog.Trace(TAG, "HandleTransaction", "Converting JSON to Transaction...")
	t := TransactionFromJSON(b)
	tracelog.Info(TAG, "HandleTransaction", t.String())
}

// Starts the main server.
func RunServer() {
	origin := "http://localhost/"
	url := "wss://ws.blockchain.info:443/inv"
	
	// Connect
	tracelog.Info(TAG, "RunServer", "Connecting to %s...\n", url)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	tracelog.Info(TAG, "RunServer", "Connected!")

	// Subscribe
	subscriptionMessage := "{\"op\":\"unconfirmed_sub\"}";
	subscriptionBytes := []byte(subscriptionMessage)
	tracelog.Info(TAG, "Subscribing with %s...\n", subscriptionMessage)
	if _, err := ws.Write(subscriptionBytes); err != nil {
		log.Fatal(err)
	}
	tracelog.Info(TAG, "RunServer", "Subscribed!")
	
	jsonData := make([]byte, 0)
	// Forever
	for {
		n := -1
		buffer := make([]byte, 1024)
		// Read
		tracelog.Trace(TAG, "RunServer", "Reading from socket...")
		if n, err = ws.Read(buffer); err != nil {
			log.Fatal(err)
		}
		
		jsonData = append(jsonData, buffer[:n]...)
		
		if IsValidJson(jsonData) {
			// Process
			go HandleTransaction(jsonData)
			
			// Restart
			jsonData = make([]byte, 0)
		}
	}
}

func main() {
	tracelog.Start(tracelog.LevelInfo)
	
	RunServer()
	
	tracelog.Stop()
}