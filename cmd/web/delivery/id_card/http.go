package id_card

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

type IDCardHandler struct {
	datastore datastore.IDCard
}

func New(idCard datastore.IDCard) IDCardHandler {
	return IDCardHandler{datastore: idCard}
}

func (a IDCardHandler) Handler(w http.ResponseWriter, r *http.Request) {
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

func (a IDCardHandler) get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	i, err := strconv.Atoi(id)
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
	tmpl := template.Must(template.ParseFiles(url+"id_card.gohtml", url+"index.gohtml"))
	_ = tmpl.ExecuteTemplate(w, "base", struct {
		Body []entities.IDCard
	}{Body: resp})
	return
}

func (a IDCardHandler) create(w http.ResponseWriter, r *http.Request) {
	var idCard entities.IDCard
	idCard.Type = r.FormValue("Type")
	idCard.Number = r.FormValue("Number")
	idCard.Issued = r.FormValue("Issued")
	resp, err := a.datastore.Create(idCard)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при создании записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(resp)
	_, _ = w.Write(body)
}

func (a IDCardHandler) delete(w http.ResponseWriter, r *http.Request) {
	var idCard entities.IDCard
	body, _ := ioutil.ReadAll(r.Body)
	data := strings.Split(string(body), "&")
	for i := 0; i < len(data); i++ {
		d := strings.Split(data[i], "=")
		if d[0] == "id_card_id" {
			id, err := strconv.Atoi(d[1])
			if err != nil {
				_, _ = w.Write([]byte("Не верный формат ID"))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			idCard.ID = id
		}
	}
	resp, err := a.datastore.Delete(idCard)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при удалении записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, _ = json.Marshal(resp)
	_, _ = w.Write(body)
}
