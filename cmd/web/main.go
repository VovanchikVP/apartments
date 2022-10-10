package main

import (
	"apartments/cmd/web/datastore/address"
	"apartments/cmd/web/datastore/animal"
	"apartments/cmd/web/datastore/apartment"
	"apartments/cmd/web/datastore/contract_rent"
	"apartments/cmd/web/datastore/counter"
	"apartments/cmd/web/datastore/id_card"
	"apartments/cmd/web/datastore/indication"
	"apartments/cmd/web/datastore/payment"
	"apartments/cmd/web/datastore/person"
	"apartments/cmd/web/datastore/property_document"
	"apartments/cmd/web/datastore/tariff"
	"apartments/cmd/web/datastore/tenant"
	"apartments/cmd/web/datastore/type_pyment"
	handlerAddress "apartments/cmd/web/delivery/address"
	handlerAnimal "apartments/cmd/web/delivery/animal"
	handlerApartment "apartments/cmd/web/delivery/apartment"
	handlerContractRent "apartments/cmd/web/delivery/contract_rent"
	handlerCounter "apartments/cmd/web/delivery/counter"
	handlerIDCard "apartments/cmd/web/delivery/id_card"
	handlerIndication "apartments/cmd/web/delivery/indication"
	handlerPayment "apartments/cmd/web/delivery/payment"
	handlerPerson "apartments/cmd/web/delivery/person"
	handlerProperty "apartments/cmd/web/delivery/property_document"
	handlerTariff "apartments/cmd/web/delivery/tariff"
	handlerTenant "apartments/cmd/web/delivery/tenant"
	handlerTypePyment "apartments/cmd/web/delivery/type_pyment"
	"apartments/cmd/web/driver"
	"fmt"
	"log"
	"net/http"
)

func main() {
	var err error

	db, err := driver.ConnectToMySQL()
	if err != nil {
		log.Println("could not connect to sql, err:", err)
		return
	}

	datastore := animal.New(db)
	handler := handlerAnimal.New(datastore)
	http.HandleFunc("/animal", handler.Handler)

	datastoreC := counter.New(db)
	handlerC := handlerCounter.New(datastoreC)
	http.HandleFunc("/counter", handlerC.Handler)

	addressDB := address.New(db)
	handlerA := handlerAddress.New(addressDB)
	http.HandleFunc("/address", handlerA.Handler)

	propertyDB := property_document.New(db)
	handlerPD := handlerProperty.New(propertyDB)
	http.HandleFunc("/property_document", handlerPD.Handler)

	idCardDB := id_card.New(db)
	handlerIC := handlerIDCard.New(idCardDB)
	http.HandleFunc("/id_cards", handlerIC.Handler)

	typePymentDB := type_pyment.New(db)
	handlerTP := handlerTypePyment.New(typePymentDB)
	http.HandleFunc("/type_payments", handlerTP.Handler)

	personDB := person.New(db)
	handlerP := handlerPerson.New(personDB)
	http.HandleFunc("/person", handlerP.Handler)

	apartmentDB := apartment.New(db)
	handlerAp := handlerApartment.New(apartmentDB)
	http.HandleFunc("/apartment", handlerAp.Handler)

	indicationDB := indication.New(db)
	handlerI := handlerIndication.New(indicationDB)
	http.HandleFunc("/indication", handlerI.Handler)

	tariffDB := tariff.New(db)
	handlerTr := handlerTariff.New(tariffDB)
	http.HandleFunc("/tariff", handlerTr.Handler)

	contractDB := contract_rent.New(db)
	handlerCR := handlerContractRent.New(contractDB)
	http.HandleFunc("/contract_rent", handlerCR.Handler)

	tenantDB := tenant.New(db)
	handlerTT := handlerTenant.New(tenantDB)
	http.HandleFunc("/tenant", handlerTT.Handler)

	paymentDB := payment.New(db)
	handlerPm := handlerPayment.New(paymentDB)
	http.HandleFunc("/payment", handlerPm.Handler)

	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/", indexHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("cmd/web/tmpl"))))
	fmt.Println(http.ListenAndServe(":9000", nil))
}
