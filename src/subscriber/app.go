package subscriber

import (
	//_ "google.golang.org/cloud"
	//_ "google.golang.org/cloud/pubsub"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/", handler)
}

type Todo struct {
}

func handler(w http.ResponseWriter, r *http.Request) {

	var todo Todo

	// read body into []byte
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// unmarshal
	if err := json.Unmarshal(body, &todo); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
