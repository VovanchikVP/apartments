package counter

import (
	"apartments/cmd/web/datastore"
	"apartments/cmd/web/entities"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type CounterHandler struct {
	datastoreCounter   datastore.Counter
	datastoreApartment datastore.Apartment
}

func New(counter datastore.Counter, apartment datastore.Apartment) CounterHandler {
	return CounterHandler{
		datastoreCounter:   counter,
		datastoreApartment: apartment,
	}
}

func (a CounterHandler) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.get(w, r)
	case http.MethodPost:
		a.create(w, r)
	case http.MethodDelete:
		a.delete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a CounterHandler) get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		_, _ = w.Write([]byte("Не верный формат ID"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := a.datastoreCounter.Get(i)
	if err != nil {
		fmt.Println(err)
		_, _ = w.Write([]byte("запись с переданным ID отсутствует в базе данных"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	respApartments, err := a.datastoreApartment.Get(0)
	url := "cmd/web/tmpl/"
	tmpl := template.Must(template.ParseFiles(url+"counter.gohtml", url+"index.gohtml"))
	_ = tmpl.ExecuteTemplate(w, "base", struct {
		Body       []entities.Counter
		Apartments []entities.Apartment
	}{
		Body:       resp,
		Apartments: respApartments,
	})
	return
}

func (a CounterHandler) create(w http.ResponseWriter, r *http.Request) {
	var counter entities.Counter
	counter.Number = r.FormValue("Number")
	counter.Type = r.FormValue("Type")
	apartmentId, err := strconv.Atoi(r.FormValue("Apartment"))
	if err != nil {
		_, _ = w.Write([]byte("Ошибка в запросе"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	counter.Apartment.ID = apartmentId
	counter.VerificationDate = r.FormValue("VerificationDate")

	resp, err := a.datastoreCounter.Create(counter)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при создании записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(resp)
	_, _ = w.Write(body)
}

func (a CounterHandler) delete(w http.ResponseWriter, r *http.Request) {
	var counter entities.Counter
	id, err := strconv.Atoi(r.FormValue("counter_id"))
	if err != nil {
		_, _ = w.Write([]byte("Не верный формат ID"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	counter.ID = id
	resp, err := a.datastoreCounter.Delete(counter)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при удалении записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, _ := json.Marshal(resp)
	_, _ = w.Write(body)
}
