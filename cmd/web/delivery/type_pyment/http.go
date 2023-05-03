package type_pyment

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

type TypePymentHandler struct {
	datastore datastore.TypePyment
}

func New(typePyment datastore.TypePyment) TypePymentHandler {
	return TypePymentHandler{datastore: typePyment}
}

func (a TypePymentHandler) Handler(w http.ResponseWriter, r *http.Request) {
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

func (a TypePymentHandler) get(w http.ResponseWriter, r *http.Request) {
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
	tmpl := template.Must(template.ParseFiles(url+"type_payment.gohtml", url+"index.gohtml"))
	_ = tmpl.ExecuteTemplate(w, "base", struct {
		Body []entities.TypePayment
	}{Body: resp})
	return
}

func (a TypePymentHandler) create(w http.ResponseWriter, r *http.Request) {
	typePayment := entities.TypePayment{
		Name: r.FormValue("Name"),
	}
	resp, err := a.datastore.Create(typePayment)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при создании записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, _ := json.Marshal(resp)
	_, _ = w.Write(body)
}

func (a TypePymentHandler) delete(w http.ResponseWriter, r *http.Request) {
	var typePayment entities.TypePayment
	body, _ := ioutil.ReadAll(r.Body)
	data := strings.Split(string(body), "&")
	for i := 0; i < len(data); i++ {
		d := strings.Split(data[i], "=")
		if d[0] == "type_payment_id" {
			id, err := strconv.Atoi(d[1])
			if err != nil {
				_, _ = w.Write([]byte("Не верный формат ID"))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			typePayment.ID = id
		}
	}
	resp, err := a.datastore.Delete(typePayment)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при удалении записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, _ = json.Marshal(resp)
	_, _ = w.Write(body)
}
