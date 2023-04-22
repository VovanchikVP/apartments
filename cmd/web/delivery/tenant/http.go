package tenant

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

type TenantHandler struct {
	datastoreTenant       datastore.Tenant
	datastoreContractRent datastore.ContractRent
	datastorePerson       datastore.Person
}

func New(tenant datastore.Tenant, contractRent datastore.ContractRent, person datastore.Person) TenantHandler {
	return TenantHandler{
		datastoreTenant:       tenant,
		datastoreContractRent: contractRent,
		datastorePerson:       person,
	}
}

func (a TenantHandler) Handler(w http.ResponseWriter, r *http.Request) {
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

func (a TenantHandler) get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	i, err := strconv.Atoi(id)

	if err != nil {
		_, _ = w.Write([]byte("Не верный формат ID"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := a.datastoreTenant.Get(i)
	if err != nil {
		fmt.Println(err)
		_, _ = w.Write([]byte("запись с переданным ID отсутствует в базе данных"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	respContractRent, err := a.datastoreContractRent.Get(0)
	respPerson, err := a.datastorePerson.Get(0)
	url := "cmd/web/tmpl/"
	tmpl := template.Must(template.ParseFiles(url+"tenant.gohtml", url+"index.gohtml"))
	_ = tmpl.ExecuteTemplate(w, "base", struct {
		Body         []entities.Tenant
		ContractRent []entities.ContractRent
		Persons      []entities.Person
	}{
		Body:         resp,
		ContractRent: respContractRent,
		Persons:      respPerson,
	})
	return
}

func (a TenantHandler) create(w http.ResponseWriter, r *http.Request) {
	var tenant entities.Tenant

	contractRentId, err := strconv.Atoi(r.FormValue("ContractRent"))
	personId, err := strconv.Atoi(r.FormValue("Person"))
	if err != nil {
		_, _ = w.Write([]byte("Ошибка в запросе"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tenant.ContractRent.ID = contractRentId
	tenant.Person.ID = personId

	resp, err := a.datastoreTenant.Create(tenant)

	if err != nil {
		_, _ = w.Write([]byte("Ошибка при создании записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(resp)
	_, _ = w.Write(body)
}

func (a TenantHandler) delete(w http.ResponseWriter, r *http.Request) {
	var tenant entities.Tenant
	body, _ := io.ReadAll(r.Body)
	data := strings.Split(string(body), "&")
	for i := 0; i < len(data); i++ {
		d := strings.Split(data[i], "=")
		if d[0] == "tenant_id" {
			id, err := strconv.Atoi(d[1])
			if err != nil {
				_, _ = w.Write([]byte("Не верный формат ID"))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			tenant.ID = id
		}
	}
	resp, err := a.datastoreTenant.Delete(tenant)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при удалении записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, _ = json.Marshal(resp)
	_, _ = w.Write(body)
}
