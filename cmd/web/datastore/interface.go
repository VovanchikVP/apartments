package datastore

import "apartments/cmd/web/entities"

type Animal interface {
	Get(id int) ([]entities.Animal, error)
	Create(entities.Animal) (entities.Animal, error)
}

type Counter interface {
	Get(id int) ([]entities.Counter, error)
	Create(counter entities.Counter) (entities.Counter, error)
	GetByID(id int) (entities.Apartment, error)
}

type Address interface {
	Get(id int) ([]entities.Address, error)
	Create(address entities.Address) (entities.Address, error)
}

type PropertyDocuments interface {
	Get(id int) ([]entities.PropertyDocuments, error)
	Create(propertyDocument entities.PropertyDocuments) (entities.PropertyDocuments, error)
	GetByID(id int) (entities.PropertyDocuments, error)
}

type IDCard interface {
	Get(id int) ([]entities.IDCard, error)
	Create(idCard entities.IDCard) (entities.IDCard, error)
	GetByID(id int) (entities.IDCard, error)
}

type TypePyment interface {
	Get(id int) ([]entities.TypePayment, error)
	Create(typePyment entities.TypePayment) (entities.TypePayment, error)
}

type Person interface {
	Get(id int) ([]entities.Person, error)
	Create(person entities.Person) (entities.Person, error)
}

type Apartment interface {
	Get(id int) ([]entities.Apartment, error)
	Create(apartment entities.Apartment) (entities.Apartment, error)
	GetByID(id int) (entities.Apartment, error)
}

type Indication interface {
	Get(id int) ([]entities.Indication, error)
	Create(indication entities.Indication) (entities.Indication, error)
	GetByID(id int) (entities.Indication, error)
}

type Tariff interface {
	Get(id int) ([]entities.Tariff, error)
	Create(tariff entities.Tariff) (entities.Tariff, error)
	GetByID(id int) (entities.Tariff, error)
}
