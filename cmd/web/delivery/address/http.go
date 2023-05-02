package address

import (
	"apartments/cmd/web/datastore"
	"apartments/cmd/web/entities"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type AddressHandler struct {
	datastore datastore.Address
}

func New(address datastore.Address) AddressHandler {
	return AddressHandler{datastore: address}
}

func (a AddressHandler) Handler(w http.ResponseWriter, r *http.Request) {
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

func (a AddressHandler) get(w http.ResponseWriter, r *http.Request) {
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

	jsonResponse := r.URL.Query().Get("json")
	if jsonResponse == "1" {
		body, _ := json.Marshal(resp)
		_, _ = w.Write(body)
		return
	}

	url := "cmd/web/tmpl/"
	tmpl := template.Must(template.ParseFiles(url+"address.gohtml", url+"index.gohtml"))
	_ = tmpl.ExecuteTemplate(w, "base", struct {
		Body []entities.Address
	}{Body: resp})
	return
}

func (a AddressHandler) create(w http.ResponseWriter, r *http.Request) {
	var address entities.Address
	index, err := strconv.Atoi(r.FormValue("Index"))
	if err != nil {
		fmt.Println("Неверный формат индекса.")
		index = 0
	}

	address.Index = index
	address.City = r.FormValue("City")
	address.Street = r.FormValue("Street")
	address.House = r.FormValue("House")
	address.Apartment = r.FormValue("Apartment")
	resp, err := a.datastore.Create(address)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при создании записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(resp)
	_, _ = w.Write(body)
}

func (a AddressHandler) delete(w http.ResponseWriter, r *http.Request) {
	var address entities.Address
	body, _ := ioutil.ReadAll(r.Body)
	data := strings.Split(string(body), "&")
	for i := 0; i < len(data); i++ {
		d := strings.Split(data[i], "=")
		if d[0] == "address_id" {
			id, err := strconv.Atoi(d[1])
			if err != nil {
				_, _ = w.Write([]byte("Не верный формат ID"))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			address.ID = id
		}
	}
	resp, err := a.datastore.Delete(address)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при удалении записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, _ = json.Marshal(resp)
	_, _ = w.Write(body)
}
