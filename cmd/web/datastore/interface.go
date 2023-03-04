package datastore

import "apartments/cmd/web/entities"

type Animal interface {
	Get(id int) ([]entities.Animal, error)
	Create(entities.Animal) (entities.Animal, error)
}

type Counter interface {
	Get(id int) ([]entities.Counter, error)
	Create(counter entities.Counter) (entities.Counter, error)
	GetByID(id int) (entities.Counter, error)
}

type Address interface {
	Get(id int) ([]entities.Address, error)
	Create(address entities.Address) (entities.Address, error)
	Delete(address entities.Address) (bool, error)
}

type PropertyDocuments interface {
	Get(id int) ([]entities.PropertyDocuments, error)
	Create(propertyDocument entities.PropertyDocuments) (entities.PropertyDocuments, error)
	GetByID(id int) (entities.PropertyDocuments, error)
	Delete(propertyDocument entities.PropertyDocuments) (bool, error)
}

type IDCard interface {
	Get(id int) ([]entities.IDCard, error)
	Create(idCard entities.IDCard) (entities.IDCard, error)
	GetByID(id int) (entities.IDCard, error)
	Delete(idCard entities.IDCard) (bool, error)
}

type TypePyment interface {
	Get(id int) ([]entities.TypePayment, error)
	Create(typePayment entities.TypePayment) (entities.TypePayment, error)
	GetByID(id int) (entities.TypePayment, error)
	Delete(typePayment entities.TypePayment) (bool, error)
}

type Person interface {
	Get(id int) ([]entities.Person, error)
	Create(person entities.Person) (entities.Person, error)
	GetByID(id int) (entities.Person, error)
}

type Apartment interface {
	Get(id int) ([]entities.Apartment, error)
	Create(apartment entities.Apartment) (entities.Apartment, error)
	GetByID(id int) (entities.Apartment, error)
	Delete(apartment entities.Apartment) (bool, error)
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

type ContractRent interface {
	Get(id int) ([]entities.ContractRent, error)
	Create(contract entities.ContractRent) (entities.ContractRent, error)
	GetByID(id int) (entities.ContractRent, error)
}

type Tenant interface {
	Get(id int) ([]entities.Tenant, error)
	Create(tenant entities.Tenant) (entities.Tenant, error)
	GetByID(id int) (entities.Tenant, error)
}

type Payment interface {
	Get(id int) ([]entities.Payment, error)
	Create(payment entities.Payment) (entities.Payment, error)
	GetByID(id int) (entities.Payment, error)
}

type OperationGroups interface {
	Get(id int) ([]entities.OperationGroups, error)
	Create(operationGroups entities.OperationGroups) (entities.OperationGroups, error)
	GetByID(id int) (entities.OperationGroups, error)
	Delete(operationGroups entities.OperationGroups) (bool, error)
}

type Operation interface {
	Get(id int) ([]entities.Operation, error)
	Create(operationGroups entities.Operation) (entities.Operation, error)
	GetByID(id int) (entities.Operation, error)
	Delete(operationGroups entities.Operation) (bool, error)
}
