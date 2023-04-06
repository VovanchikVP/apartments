package contract_rent

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

type ContractHandler struct {
	datastoreContract  datastore.ContractRent
	datastorePerson    datastore.Person
	datastoreApartment datastore.Apartment
}

func New(contract datastore.ContractRent, person datastore.Person, apartment datastore.Apartment) ContractHandler {
	return ContractHandler{
		datastoreContract:  contract,
		datastorePerson:    person,
		datastoreApartment: apartment,
	}
}

func (a ContractHandler) Handler(w http.ResponseWriter, r *http.Request) {
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

func (a ContractHandler) get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	i, err := strconv.Atoi(id)
	fmt.Println(i)
	if err != nil {
		_, _ = w.Write([]byte("Не верный формат ID"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := a.datastoreContract.Get(i)
	if err != nil {
		fmt.Println(err)
		_, _ = w.Write([]byte("запись с переданным ID отсутствует в базе данных"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rApartment, err := a.datastoreApartment.Get(0)
	rPerson, err := a.datastorePerson.Get(0)

	url := "cmd/web/tmpl/"
	tmpl := template.Must(template.ParseFiles(url+"contract_rent.gohtml", url+"index.gohtml"))
	_ = tmpl.ExecuteTemplate(w, "base", struct {
		Body       []entities.ContractRent
		Apartments []entities.Apartment
		Persons    []entities.Person
	}{
		Body:       resp,
		Apartments: rApartment,
		Persons:    rPerson,
	})
	return
}

func (a ContractHandler) create(w http.ResponseWriter, r *http.Request) {
	var contract entities.ContractRent
	var err error
	contract.Number = r.FormValue("Number")
	contract.Date = r.FormValue("Date")
	contract.Employer.ID, err = strconv.Atoi(r.FormValue("EmployerID"))
	contract.Landlord.ID, err = strconv.Atoi(r.FormValue("LandlordID"))
	contract.Apartment.ID, err = strconv.Atoi(r.FormValue("ApartmentID"))
	contract.DateStartRent = r.FormValue("DateStartRent")
	contract.DateEndRent = r.FormValue("DateEndRent")
	contract.DateApartmentTransfer = r.FormValue("DateApartmentTransfer")
	rental, err := strconv.ParseFloat(r.FormValue("Rental"), 32)
	contract.Rental = float32(rental)
	contract.DateRental = r.FormValue("DateRental")
	deposit, err := strconv.ParseFloat(r.FormValue("Deposit"), 32)
	contract.Deposit = float32(deposit)
	transferred_amount, err := strconv.ParseFloat(r.FormValue("Rental"), 32)
	contract.Rental = float32(transferred_amount)
	contract.PaymentsCommunal, err = strconv.ParseBool(r.FormValue("PaymentsCommunal"))
	contract.PaymentsNetwork, err = strconv.ParseBool(r.FormValue("PaymentsNetwork"))
	contract.PaymentsElectric, err = strconv.ParseBool(r.FormValue("PaymentsElectric"))
	contract.PaymentsHeating, err = strconv.ParseBool(r.FormValue("PaymentsHeating"))
	contract.PaymentsColdWater, err = strconv.ParseBool(r.FormValue("PaymentsColdWater"))
	contract.PaymentsHotWater, err = strconv.ParseBool(r.FormValue("PaymentsHotWater"))
	contract.AdditionalTerms = r.FormValue("AdditionalTerms")
	contract.FileContract = ""

	if err != nil {
		_, _ = w.Write([]byte("Ошибка в запросе"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := a.datastoreContract.Create(contract)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при создании записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(resp)
	_, _ = w.Write(body)
}

func (a ContractHandler) delete(w http.ResponseWriter, r *http.Request) {
	var contract entities.ContractRent
	body, _ := io.ReadAll(r.Body)
	data := strings.Split(string(body), "&")
	for i := 0; i < len(data); i++ {
		d := strings.Split(data[i], "=")
		if d[0] == "contract_id" {
			id, err := strconv.Atoi(d[1])
			if err != nil {
				_, _ = w.Write([]byte("Не верный формат ID"))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			contract.ID = id
		}
	}
	resp, err := a.datastoreContract.Delete(contract)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при удалении записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, _ = json.Marshal(resp)
	_, _ = w.Write(body)
}
