package counter

import (
	"apartments/cmd/web/datastore"
	"apartments/cmd/web/entities"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type CounterHandler struct {
	datastore datastore.Counter
}

func New(counter datastore.Counter) CounterHandler {
	return CounterHandler{datastore: counter}
}

func (a CounterHandler) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.get(w, r)
	case http.MethodPost:
		a.create(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a CounterHandler) get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	i, err := strconv.Atoi(id)
	fmt.Println(i)
	if err != nil {
		_, _ = w.Write([]byte("Не верный формат ID"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := a.datastore.Get(i)
	if err != nil {
		fmt.Println(err)
		_, _ = w.Write([]byte("запись с переданным ID отсутствует в базе данных"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, _ := json.Marshal(resp)
	_, _ = w.Write(body)
}

func (a CounterHandler) create(w http.ResponseWriter, r *http.Request) {
	var counter entities.Counter

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &counter)
	if err != nil {
		fmt.Println(string(body))
		_, _ = w.Write([]byte("Ошибка в запросе"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := a.datastore.Create(counter)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при создании записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, _ = json.Marshal(resp)
	_, _ = w.Write(body)
}
