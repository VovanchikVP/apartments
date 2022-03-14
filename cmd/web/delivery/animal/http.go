package animal

import (
	"apartments/cmd/web/datastore"
	"apartments/cmd/web/entities"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type AnimalHandler struct {
	datastore datastore.Animal
}

func New(animal datastore.Animal) AnimalHandler {
	return AnimalHandler{datastore: animal}
}

func (a AnimalHandler) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.get(w, r)
	case http.MethodPost:
		a.create(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a AnimalHandler) get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	i, err := strconv.Atoi(id)
	if err != nil {
		_, _ = w.Write([]byte("invalid parameter id"))
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	resp, err := a.datastore.Get(i)
	if err != nil {
		_, _ = w.Write([]byte("could not retrieve animal"))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	body, _ := json.Marshal(resp)
	_, _ = w.Write(body)
}

func (a AnimalHandler) create(w http.ResponseWriter, r *http.Request) {
	var animal entities.Animal

	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, &animal)
	if err != nil {
		fmt.Println(err)
		_, _ = w.Write([]byte("invalid body"))
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	resp, err := a.datastore.Create(animal)
	if err != nil {
		_, _ = w.Write([]byte("could not create animal"))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	body, _ = json.Marshal(resp)
	_, _ = w.Write(body)
}
