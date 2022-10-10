package contract_rent

import (
	"apartments/cmd/web/datastore/apartment"
	"apartments/cmd/web/datastore/person"
	"apartments/cmd/web/driver"
	"apartments/cmd/web/entities"
	"database/sql"
	"log"
)

type ContractRentStorer struct {
	db *sql.DB
}

func New(db *sql.DB) ContractRentStorer {
	return ContractRentStorer{db: db}
}

func (a ContractRentStorer) GetByID(id int) (contract entities.ContractRent, err error){
	var row *sql.Row

	row = a.db.QueryRow("SELECT ROWID, * FROM contracts_rent WHERE ROWID = ?", id)

	db, err := driver.ConnectToMySQL()
	if err != nil {
		log.Println("could not connect to sql, err:", err)
		return entities.ContractRent{}, err
	}

	apartmentDB := apartment.New(db)
	personDB := person.New(db)

	switch err = row.Scan(&contract.ID, &contract.Number, &contract.Date, &contract.Employer.ID, &contract.Landlord.ID, &contract.Apartment.ID, &contract.DateStartRent, &contract.DateEndRent, &contract.DateApartmentTransfer, &contract.Rental, &contract.DateRental, &contract.Deposit, &contract.TransferredAmount, &contract.PaymentsCommunal, &contract.PaymentsNetwork, &contract.PaymentsElectric, &contract.PaymentsHeating, &contract.PaymentsColdWater, &contract.PaymentsHotWater, &contract.AdditionalTerms, &contract.FileContract); err {
	case sql.ErrNoRows:
		return entities.ContractRent{}, err
	case nil:
		contract.Apartment, _ = apartmentDB.GetByID(contract.Apartment.ID)
		contract.Employer, _ = personDB.GetByID(contract.Employer.ID)
		contract.Landlord, _ = personDB.GetByID(contract.Landlord.ID)
		return contract, nil
	default:
		return entities.ContractRent{}, err
	}
}

func (a ContractRentStorer) Get(id int) (contract []entities.ContractRent, err error) {
	var rows *sql.Rows

	if id != 0 {
		rows, err = a.db.Query("SELECT * FROM contracts_rent WHERE ROWID = ?", id)
	} else {
		rows, err = a.db.Query("SELECT * FROM contracts_rent")
	}

	if rows != nil {
		defer func() {
			_ = rows.Close()
		}()
	}

	if err != nil {
		return nil, err
	}

	db, err := driver.ConnectToMySQL()
	if err != nil {
		log.Println("could not connect to sql, err:", err)
		return nil, err
	}

	apartmentDB := apartment.New(db)
	personDB := person.New(db)

	for rows.Next() {
		var a entities.ContractRent
		_ = rows.Scan(&a.ID, &a.Number, &a.Date, &a.Employer.ID, &a.Landlord.ID, &a.Apartment.ID, &a.DateStartRent, &a.DateEndRent, &a.DateApartmentTransfer, &a.Rental, &a.DateRental, &a.Deposit, &a.TransferredAmount, &a.PaymentsCommunal, &a.PaymentsNetwork, &a.PaymentsElectric, &a.PaymentsHeating, &a.PaymentsColdWater, &a.PaymentsHotWater, &a.AdditionalTerms, &a.FileContract)
		a.Apartment, _ = apartmentDB.GetByID(a.Apartment.ID)
		a.Employer, _ = personDB.GetByID(a.Employer.ID)
		a.Landlord, _ = personDB.GetByID(a.Landlord.ID)
		contract = append(contract, a)
	}

	return contract, nil
}

func (a ContractRentStorer) Create(contract entities.ContractRent) (entities.ContractRent, error) {
	res, err := a.db.Exec("INSERT INTO contracts_rent(number, date, employer_id, landlord_id, apartment_id, date_start_rent, date_end_rent, date_apartment_transfer, rental, date_rental, deposit, transferred_amount, payments_communal, payments_network, payments_electric, payments_heating, payments_cold_water, payments_hot_water, additional_terms, file_contract) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", contract.Number, contract.Date, contract.Employer.ID, contract.Landlord.ID, contract.Apartment.ID, contract.DateStartRent, contract.DateEndRent, contract.DateApartmentTransfer, contract.Rental, contract.DateRental, contract.Deposit, contract.TransferredAmount, contract.PaymentsCommunal, contract.PaymentsNetwork, contract.PaymentsElectric, contract.PaymentsHeating, contract.PaymentsColdWater, contract.PaymentsHotWater, contract.AdditionalTerms, contract.FileContract)

	if err != nil {
		return entities.ContractRent{}, err
	}

	id, _ := res.LastInsertId()
	contract.ID = int(id)

	return contract, nil
}
