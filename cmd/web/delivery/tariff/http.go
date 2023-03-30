package tariff

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

type TariffHandler struct {
	datastoreTariff  datastore.Tariff
	datastoreCounter datastore.Counter
}

func New(tariff datastore.Tariff, counter datastore.Counter) TariffHandler {
	return TariffHandler{
		datastoreTariff:  tariff,
		datastoreCounter: counter,
	}
}

func (a TariffHandler) Handler(w http.ResponseWriter, r *http.Request) {
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

func (a TariffHandler) get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	i, err := strconv.Atoi(id)
	fmt.Println(i)
	if err != nil {
		_, _ = w.Write([]byte("Не верный формат ID"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := a.datastoreTariff.Get(i)
	if err != nil {
		fmt.Println(err)
		_, _ = w.Write([]byte("запись с переданным ID отсутствует в базе данных"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	respCounters, err := a.datastoreCounter.Get(0)
	url := "cmd/web/tmpl/"
	tmpl := template.Must(template.ParseFiles(url+"tariff.gohtml", url+"index.gohtml"))
	_ = tmpl.ExecuteTemplate(w, "base", struct {
		Body     []entities.Tariff
		Counters []entities.Counter
	}{
		Body:     resp,
		Counters: respCounters,
	})
	return
}

func (a TariffHandler) create(w http.ResponseWriter, r *http.Request) {
	var tariff entities.Tariff
	tariff.SetDate = r.FormValue("SetDate")
	tariffCost, err := strconv.ParseFloat(r.FormValue("Cost"), 32)
	tariffCount, err := strconv.Atoi(r.FormValue("Counter"))
	tariff.Cost = float32(tariffCost)
	tariff.Counter = entities.Counter{ID: tariffCount}
	if err != nil {
		_, _ = w.Write([]byte("Ошибка в запросе"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := a.datastoreTariff.Create(tariff)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при создании записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(resp)
	_, _ = w.Write(body)
}

func (a TariffHandler) delete(w http.ResponseWriter, r *http.Request) {
	var tariff entities.Tariff
	body, _ := io.ReadAll(r.Body)
	data := strings.Split(string(body), "&")
	for i := 0; i < len(data); i++ {
		d := strings.Split(data[i], "=")
		if d[0] == "tariff_id" {
			id, err := strconv.Atoi(d[1])
			if err != nil {
				_, _ = w.Write([]byte("Не верный формат ID"))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			tariff.ID = id
		}
	}
	resp, err := a.datastoreTariff.Delete(tariff)
	if err != nil {
		_, _ = w.Write([]byte("Ошибка при удалении записи."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, _ = json.Marshal(resp)
	_, _ = w.Write(body)
}
