package payment

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

type PaymentHandler struct {
	datastorePayment    datastore.Payment
	datastoreApartment  datastore.Apartment
	datastoreTypePyment datastore.TypePyment
}

func New(payment datastore.Payment, apartment datastore.Apartment, typePayment datastore.TypePyment) PaymentHandler {
	return PaymentHandler{
		datastorePayment:    payment,
		datastoreApartment:  apartment,
		datastoreTypePyment: typePayment,
	}
}

func (a PaymentHandler) Handler(w http.ResponseWriter, r *http.Request) {
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

func (a PaymentHandler) get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	i, err := strconv.Atoi(id)
	fmt.Println(i)
	if err != nil {
		_, _ = w.Write([]byte("Не верный формат ID"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := a.datastorePayment.Get(i)
	if err != nil {
		fmt.Println(err)
		_, _ = w.Write([]byte("запись с переданным ID отсутствует в базе данных"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	respApartments, err := a.datastoreApartment.Get(0)
	respTypePayment, err := a.datastoreTypePyment.Get(0)

	url := "cmd/web/tmpl/"
	tmpl := template.Must(template.ParseFiles(url+"payment.gohtml", url+"index.gohtml"))
	_ = tmpl.ExecuteTemplate(w, "base", struct {
		Body        []entities.Payment
		Apartments  []entities.Apartment
		TypePayment []entities.TypePayment
	}{
		Body:        resp,
		Apartments:  respApartments,
		TypePayment: respTypePayment,
	})
	return
}

func (a PaymentHandler) create(w http.ResponseWriter, r *http.Request) {
	var payment entities.Payment

	cost, err := strconv.ParseFloat(r.FormValue("Cost"), 32)
	payment.Date = r.FormValue("Type")
	apartmentId, err := strconv.Atoi(r.FormValue("Apartment"))
	typePaymentId, err := strconv.Atoi(r.FormValue("TypePayment"))
	if err != nil {
		_, _ = w.Write([]byte("Ошибка в запросе"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	payment.Cost = float32(cost)
	payment.Apartment.ID = apartmentId
	payment.Type.ID = typePaymentId
	payment.Admission, err = strconv.ParseBool(r.FormValue("Admission"))
	if err != nil {
		payment.Admission = false
	}

	resp, err := a.datastorePayment.Create(payment)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при создании записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(resp)
	_, _ = w.Write(body)
}

func (a PaymentHandler) delete(w http.ResponseWriter, r *http.Request) {
	var payment entities.Payment
	body, _ := io.ReadAll(r.Body)
	data := strings.Split(string(body), "&")
	for i := 0; i < len(data); i++ {
		d := strings.Split(data[i], "=")
		if d[0] == "payment_id" {
			id, err := strconv.Atoi(d[1])
			if err != nil {
				_, _ = w.Write([]byte("Не верный формат ID"))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			payment.ID = id
		}
	}
	resp, err := a.datastorePayment.Delete(payment)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при удалении записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, _ = json.Marshal(resp)
	_, _ = w.Write(body)
}
