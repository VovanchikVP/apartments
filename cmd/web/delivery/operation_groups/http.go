package operation_groups

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

type OperationGroupHandler struct {
	datastore datastore.OperationGroups
}

func New(operationGroups datastore.OperationGroups) OperationGroupHandler {
	return OperationGroupHandler{datastore: operationGroups}
}

func (a OperationGroupHandler) Handler(w http.ResponseWriter, r *http.Request) {
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

func (a OperationGroupHandler) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)

	//url := "cmd/web/tmpl/"
	//tmpl := template.Must(template.ParseFiles(url+"operation_groups.gohtml", url+"index.gohtml"))
	//_ = tmpl.ExecuteTemplate(w, "base", struct {
	//	Body []entities.OperationGroups
	//}{Body: resp})
	return
}

func (a OperationGroupHandler) create(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	for key, value := range r.Form {
		fmt.Printf("%s = %s\n", key, value)
	}
	operationGroup := entities.OperationGroups{
		Name: r.FormValue("Name"),
	}
	resp, err := a.datastore.Create(operationGroup)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при создании записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, _ := json.Marshal(resp)
	_, _ = w.Write(body)
}

func (a OperationGroupHandler) delete(w http.ResponseWriter, r *http.Request) {
	var operationGroups entities.OperationGroups
	body, _ := io.ReadAll(r.Body)
	data := strings.Split(string(body), "&")
	for i := 0; i < len(data); i++ {
		d := strings.Split(data[i], "=")
		if d[0] == "operation_groups_id" {
			id, err := strconv.Atoi(d[1])
			if err != nil {
				_, _ = w.Write([]byte("Не верный формат ID"))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			operationGroups.ID = id
		}
	}
	resp, err := a.datastore.Delete(operationGroups)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при удалении записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, _ = json.Marshal(resp)
	_, _ = w.Write(body)
}
