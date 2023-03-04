package apartment

import (
	"apartments/cmd/web/datastore"
	"apartments/cmd/web/entities"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type ApartmentHandler struct {
	datastore                  datastore.Apartment
	datastoreAddress           datastore.Address
	datastorePropertyDocuments datastore.PropertyDocuments
}

func New(
	apartment datastore.Apartment,
	address datastore.Address,
	propertyDocuments datastore.PropertyDocuments,
) ApartmentHandler {
	return ApartmentHandler{
		datastore:                  apartment,
		datastoreAddress:           address,
		datastorePropertyDocuments: propertyDocuments,
	}
}

func (a ApartmentHandler) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.get(w, r)
	case http.MethodPost:
		a.create(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a ApartmentHandler) get(w http.ResponseWriter, r *http.Request) {
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

	respAddress, _ := a.datastoreAddress.Get(0)
	respPropertyDocuments, _ := a.datastorePropertyDocuments.Get(0)

	url := "cmd/web/tmpl/"
	tmpl := template.Must(template.ParseFiles(url+"apartment.gohtml", url+"index.gohtml"))
	_ = tmpl.ExecuteTemplate(w, "base", struct {
		Body              []entities.Apartment
		Address           []entities.Address
		PropertyDocuments []entities.PropertyDocuments
	}{
		Body:              resp,
		Address:           respAddress,
		PropertyDocuments: respPropertyDocuments,
	})
	return
}

func (a ApartmentHandler) create(w http.ResponseWriter, r *http.Request) {
	address, err := strconv.Atoi(r.FormValue("Address"))
	countRooms, err := strconv.Atoi(r.FormValue("CountRooms"))
	propertyDocument, err := strconv.Atoi(r.FormValue("PropertyDocument"))
	if err != nil {
		fmt.Println(err)
		_, _ = w.Write([]byte("Ошибка в запросе"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rent := false
	if r.FormValue("Rent") == "on" {
		rent = true
	}

	apartment := entities.Apartment{
		Address:           entities.Address{ID: address},
		CountRooms:        countRooms,
		PropertyDocuments: entities.PropertyDocuments{ID: propertyDocument},
		Rent:              rent,
	}

	resp, err := a.datastore.Create(apartment)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при создании записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(resp)
	_, _ = w.Write(body)
}
