package payment

import (
	"apartments/cmd/web/datastore/apartment"
	"apartments/cmd/web/datastore/type_pyment"
	"apartments/cmd/web/driver"
	"apartments/cmd/web/entities"
	"database/sql"
	"log"
)

type PaymentStorer struct {
	db *sql.DB
}

func New(db *sql.DB) PaymentStorer {
	return PaymentStorer{db: db}
}

func (a PaymentStorer) GetByID(id int) (payment entities.Payment, err error){
	var row *sql.Row

	row = a.db.QueryRow("SELECT ROWID, * FROM payments WHERE ROWID = ?", id)

	db, err := driver.ConnectToMySQL()
	if err != nil {
		log.Println("could not connect to sql, err:", err)
		return entities.Payment{}, err
	}

	apartmentDB := apartment.New(db)
	typePaymentDB := type_pyment.New(db)

	switch err = row.Scan(&payment.ID, &payment.Apartment.ID, &payment.Cost, &payment.Admission, &payment.Type.ID, &payment.Date); err {
	case sql.ErrNoRows:
		return entities.Payment{}, err
	case nil:
		payment.Apartment, _ = apartmentDB.GetByID(payment.Apartment.ID)
		payment.Type, _ = typePaymentDB.GetByID(payment.Type.ID)
		return payment, nil
	default:
		return entities.Payment{}, err
	}
}

func (a PaymentStorer) Get(id int) (payment []entities.Payment, err error) {
	var rows *sql.Rows

	if id != 0 {
		rows, err = a.db.Query("SELECT * FROM payments WHERE ROWID = ?", id)
	} else {
		rows, err = a.db.Query("SELECT * FROM payments")
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
	typePaymentDB := type_pyment.New(db)

	for rows.Next() {
		var a entities.Payment
		_ = rows.Scan(&a.ID, &a.Apartment.ID, &a.Cost, &a.Admission, &a.Type.ID, &a.Date)
		a.Apartment, _ = apartmentDB.GetByID(a.Apartment.ID)
		a.Type, _ = typePaymentDB.GetByID(a.Type.ID)
		payment = append(payment, a)
	}

	return payment, nil
}

func (a PaymentStorer) Create(payment entities.Payment) (entities.Payment, error) {
	res, err := a.db.Exec("INSERT INTO payments(apartment_id, cost, admission, type_payment_id, date) VALUES (?, ?, ?, ?, ?)", payment.Apartment.ID, payment.Cost, payment.Admission, payment.Type.ID, payment.Date)

	if err != nil {
		return entities.Payment{}, err
	}

	id, _ := res.LastInsertId()
	payment.ID = int(id)

	return payment, nil
}
