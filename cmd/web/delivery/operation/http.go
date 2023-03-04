package operation

import (
	"apartments/cmd/web/datastore"
	"apartments/cmd/web/entities"
	"encoding/json"
	"fmt"
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
	return OperationHandler{datastoreOperation: operation,
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

func (a OperationHandler) options(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
	return
}

type FormOperation struct {
	Date           string
	OperationGroup string
	Type           string
	Value          string
}

func (a OperationHandler) create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	var form FormOperation
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		fmt.Println(err)
		fmt.Println(1)
		_, _ = w.Write([]byte("Ошибка в запросе"))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("test")
		return
	}
	fmt.Println(form)
	value, err := strconv.ParseFloat(form.Value, 32)
	if err != nil {
		fmt.Println(err)
		fmt.Println(2)
		_, _ = w.Write([]byte("Ошибка в запросе"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	operationGroupId, err := strconv.Atoi(form.OperationGroup)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка в запросе"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	operation := entities.Operation{
		Date:  form.Date,
		Type:  form.Type,
		Proof: true,
		Group: entities.OperationGroups{ID: operationGroupId},
		Value: value,
	}

	resp, err := a.datastoreOperation.Create(operation)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при создании записи."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res, err := a.datastoreOperation.Get(resp.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	jData, err := json.Marshal(res)
	_, _ = w.Write(jData)
	return
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

func (a OperationHandler) patch(w http.ResponseWriter, r *http.Request) {
	fmt.Println("patch")
}
