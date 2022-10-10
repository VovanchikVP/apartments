package property_document

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

type PropertyHandler struct {
	datastore datastore.PropertyDocuments
}

func New(propertyDocument datastore.PropertyDocuments) PropertyHandler {
	return PropertyHandler{datastore: propertyDocument}
}

func (a PropertyHandler) Handler(w http.ResponseWriter, r *http.Request) {
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

func (a PropertyHandler) get(w http.ResponseWriter, r *http.Request) {
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

	url := "cmd/web/tmpl/"
	tmpl := template.Must(template.ParseFiles(url+"property_document.gohtml", url+"index.gohtml"))
	_ = tmpl.ExecuteTemplate(w, "base", struct {
		Body []entities.PropertyDocuments
	}{Body: resp})
	return
}

func (a PropertyHandler) create(w http.ResponseWriter, r *http.Request) {
	var propertyDocument entities.PropertyDocuments

	propertyDocument.Type = r.FormValue("Type")
	propertyDocument.Number = r.FormValue("Number")
	propertyDocument.Date = r.FormValue("Date")

	resp, err := a.datastore.Create(propertyDocument)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при создании записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(resp)
	_, _ = w.Write(body)
}

func (a PropertyHandler) delete(w http.ResponseWriter, r *http.Request) {
	var propertyDocument entities.PropertyDocuments
	body, _ := ioutil.ReadAll(r.Body)
	data := strings.Split(string(body), "&")
	for i:=0; i<len(data); i++ {
		d := strings.Split(data[i], "=")
		if d[0] == "property_document_id"{
			id, err := strconv.Atoi(d[1])
			if err != nil {
				_, _ = w.Write([]byte("Не верный формат ID"))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			propertyDocument.ID = id
		}
	}
	resp, err := a.datastore.Delete(propertyDocument)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при удалении записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, _ = json.Marshal(resp)
	_, _ = w.Write(body)
}
