package indication

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

type IndicationHandler struct {
	datastoreIndication datastore.Indication
	datastoreCounter    datastore.Counter
}

func New(indication datastore.Indication, counter datastore.Counter) IndicationHandler {
	return IndicationHandler{
		datastoreIndication: indication,
		datastoreCounter:    counter,
	}
}

func (a IndicationHandler) Handler(w http.ResponseWriter, r *http.Request) {
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

func (a IndicationHandler) get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	i, err := strconv.Atoi(id)
	if err != nil {
		_, _ = w.Write([]byte("Не верный формат ID"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := a.datastoreIndication.Get(i)
	if err != nil {
		fmt.Println(err)
		_, _ = w.Write([]byte("запись с переданным ID отсутствует в базе данных"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	respCounters, err := a.datastoreCounter.Get(0)
	url := "cmd/web/tmpl/"
	tmpl := template.Must(template.ParseFiles(url+"indication.gohtml", url+"index.gohtml"))
	_ = tmpl.ExecuteTemplate(w, "base", struct {
		Body     []entities.Indication
		Counters []entities.Counter
	}{
		Body:     resp,
		Counters: respCounters,
	})
	return
}

func (a IndicationHandler) create(w http.ResponseWriter, r *http.Request) {
	var indication entities.Indication
	indication.Date = r.FormValue("Date")
	indicationConv, err := strconv.ParseFloat(r.FormValue("Data"), 32)
	indicationCount, err := strconv.Atoi(r.FormValue("Counter"))
	indication.Data = float32(indicationConv)
	indication.Counter = entities.Counter{ID: indicationCount}
	if err != nil {
		_, _ = w.Write([]byte("Ошибка в запросе"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := a.datastoreIndication.Create(indication)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при создании записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(resp)
	_, _ = w.Write(body)
}

func (a IndicationHandler) delete(w http.ResponseWriter, r *http.Request) {
	var indication entities.Indication
	body, _ := io.ReadAll(r.Body)
	data := strings.Split(string(body), "&")
	for i := 0; i < len(data); i++ {
		d := strings.Split(data[i], "=")
		if d[0] == "indication_id" {
			id, err := strconv.Atoi(d[1])
			if err != nil {
				_, _ = w.Write([]byte("Не верный формат ID"))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			indication.ID = id
		}
	}
	resp, err := a.datastoreIndication.Delete(indication)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при удалении записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, _ = json.Marshal(resp)
	_, _ = w.Write(body)
}
