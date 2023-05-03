package operation

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

type OperationHandler struct {
	datastoreOperation       datastore.Operation
	datastoreOperationGroups datastore.OperationGroups
}

func New(operation datastore.Operation, operationGroups datastore.OperationGroups) OperationHandler {
	return OperationHandler{
		datastoreOperation:       operation,
		datastoreOperationGroups: operationGroups}
}

func (a OperationHandler) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		a.options(w, r)
	case http.MethodGet:
		a.get(w, r)
	case http.MethodPost:
		a.create(w, r)
	case http.MethodDelete:
		a.delete(w, r)
	case http.MethodPatch:
		a.patch(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a OperationHandler) options(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	w.WriteHeader(http.StatusNoContent)
	return
}

func (a OperationHandler) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	id := r.URL.Query().Get("id")

	i, err := strconv.Atoi(id)
	if err != nil {
		_, _ = w.Write([]byte("Не верный формат ID"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := a.datastoreOperation.Get(i)
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

	respOperationGroups, err := a.datastoreOperationGroups.Get(0)

	url := "cmd/web/tmpl/"
	tmpl := template.Must(template.ParseFiles(url+"operation.gohtml", url+"index.gohtml"))
	_ = tmpl.ExecuteTemplate(w, "base", struct {
		Body            []entities.Operation
		OperationGroups []entities.OperationGroups
	}{
		Body:            resp,
		OperationGroups: respOperationGroups,
	})
	return
}

func (a OperationHandler) create(w http.ResponseWriter, r *http.Request) {
	proof := false
	if r.FormValue("Proof") == "on" {
		proof = true
	}
	value, err := strconv.ParseFloat(r.FormValue("Value"), 32)
	if err != nil {
		fmt.Println(err)
	}
	operationGroupId, err := strconv.Atoi(r.FormValue("OperationGroups"))
	if err != nil {
		fmt.Println(err)
		_, _ = w.Write([]byte("Ошибка в запросе"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	operation := entities.Operation{
		Date:         r.FormValue("Date"),
		Type:         r.FormValue("Type"),
		Proof:        proof,
		Group:        entities.OperationGroups{ID: operationGroupId},
		Value:        value,
		Descriptions: r.FormValue("Descriptions"),
	}

	resp, err := a.datastoreOperation.Create(operation)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при создании записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(resp)
	_, _ = w.Write(body)
}

func (a OperationHandler) delete(w http.ResponseWriter, r *http.Request) {
	var operation entities.Operation
	body, _ := io.ReadAll(r.Body)
	data := strings.Split(string(body), "&")
	for i := 0; i < len(data); i++ {
		d := strings.Split(data[i], "=")
		if d[0] == "operation_id" {
			id, err := strconv.Atoi(d[1])
			if err != nil {
				_, _ = w.Write([]byte("Не верный формат ID"))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			operation.ID = id
		}
	}
	resp, err := a.datastoreOperation.Delete(operation)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при удалении записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, _ = json.Marshal(resp)
	_, _ = w.Write(body)
}

func (a OperationHandler) patch(_ http.ResponseWriter, _ *http.Request) {
	fmt.Println("patch")
}
