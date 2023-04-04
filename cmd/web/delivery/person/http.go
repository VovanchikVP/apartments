package person

import (
	"apartments/cmd/web/datastore"
	"apartments/cmd/web/entities"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type PersonHandler struct {
	datastorePerson  datastore.Person
	datastoreIDCARD  datastore.IDCard
	datastoreAddress datastore.Address
}

func New(dsPerson datastore.Person, dsIDCard datastore.IDCard, dsAddress datastore.Address) PersonHandler {
	return PersonHandler{
		datastorePerson:  dsPerson,
		datastoreIDCARD:  dsIDCard,
		datastoreAddress: dsAddress,
	}
}

func (a PersonHandler) Handler(w http.ResponseWriter, r *http.Request) {
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

func (a PersonHandler) get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	i, err := strconv.Atoi(id)
	fmt.Println(i)
	if err != nil {
		_, _ = w.Write([]byte("Не верный формат ID"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := a.datastorePerson.Get(i)
	if err != nil {
		fmt.Println(err)
		_, _ = w.Write([]byte("запись с переданным ID отсутствует в базе данных"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rIDCard, err := a.datastoreIDCARD.Get(0)
	rAddress, err := a.datastoreAddress.Get(0)

	url := "cmd/web/tmpl/"
	tmpl := template.Must(template.ParseFiles(url+"person.gohtml", url+"index.gohtml"))
	_ = tmpl.ExecuteTemplate(w, "base", struct {
		Body    []entities.Person
		IDCards []entities.IDCard
		Address []entities.Address
	}{
		Body:    resp,
		IDCards: rIDCard,
		Address: rAddress,
	})
	return
}

func (a PersonHandler) create(w http.ResponseWriter, r *http.Request) {
	var person entities.Person
	person.FirstName = r.FormValue("FirstName")
	person.LastName = r.FormValue("LastName")
	person.Patronymic = r.FormValue("Patronymic")
	person.Phone = r.FormValue("Phone")
	IDCard, err := strconv.Atoi(r.FormValue("IDCard"))
	person.IDCard.ID = IDCard
	Address, err := strconv.Atoi(r.FormValue("Address"))
	person.Address.ID = Address
	if err != nil {
		_, _ = w.Write([]byte("Ошибка в запросе"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := a.datastorePerson.Create(person)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при создании записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(resp)
	_, _ = w.Write(body)
}

func (a PersonHandler) delete(w http.ResponseWriter, r *http.Request) {
	var person entities.Person
	body, _ := io.ReadAll(r.Body)
	data := strings.Split(string(body), "&")
	for i := 0; i < len(data); i++ {
		d := strings.Split(data[i], "=")
		if d[0] == "person_id" {
			id, err := strconv.Atoi(d[1])
			if err != nil {
				_, _ = w.Write([]byte("Не верный формат ID"))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			person.ID = id
		}
	}
	resp, err := a.datastorePerson.Delete(person)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при удалении записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, _ = json.Marshal(resp)
	_, _ = w.Write(body)
}
