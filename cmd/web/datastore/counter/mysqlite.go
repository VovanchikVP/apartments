package counter

import (
	"apartments/cmd/web/datastore/apartment"
	"apartments/cmd/web/driver"
	"apartments/cmd/web/entities"
	"database/sql"
	"log"
)

type CounterStorer struct {
	db *sql.DB
}

func New(db *sql.DB) CounterStorer {
	return CounterStorer{db: db}
}

func (a CounterStorer) GetByID(id int) (counter entities.Counter, err error) {
	var row *sql.Row

	row = a.db.QueryRow(`SELECT c.ROWID, c.type, c.number, c.verification_date, c.apartment_id 
							    FROM counters c WHERE ROWID = ?`, id)

	db, err := driver.ConnectToMySQL()
	if err != nil {
		log.Println("could not connect to sql, err:", err)
		return entities.Counter{}, err
	}

	apartmentDB := apartment.New(db)

	switch err = row.Scan(&counter.ID, &counter.Type, &counter.Number, &counter.VerificationDate, &counter.Apartment.ID); err {
	case sql.ErrNoRows:
		return entities.Counter{}, err
	case nil:
		counter.Apartment, _ = apartmentDB.GetByID(counter.Apartment.ID)
		return counter, nil
	default:
		return entities.Counter{}, err
	}
}

func (a CounterStorer) Get(id int) (counters []entities.Counter, err error) {
	var rows *sql.Rows

	if id != 0 {
		rows, err = a.db.Query(`SELECT c.ROWID, c.type, c.number, c.verification_date, c.apartment_id 
									  FROM counters c WHERE ROWID = ?`, id)
	} else {
		rows, err = a.db.Query(`SELECT c.ROWID, c.type, c.number, c.verification_date, c.apartment_id  
    								  FROM counters c`)
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

	for rows.Next() {
		var a entities.Counter
		_ = rows.Scan(&a.ID, &a.Type, &a.Number, &a.VerificationDate, &a.Apartment.ID)
		a.Apartment, _ = apartmentDB.GetByID(a.Apartment.ID)
		counters = append(counters, a)
	}
	return counters, nil
}

func (a CounterStorer) Create(counter entities.Counter) (entities.Counter, error) {
	res, err := a.db.Exec("INSERT INTO counters(type, number, verification_date, apartment_id) VALUES (?, ?, ?, ?)", counter.Type, counter.Number, counter.VerificationDate, counter.Apartment.ID)

	if err != nil {
		return entities.Counter{}, err
	}

	id, _ := res.LastInsertId()
	counter.ID = int(id)

	return counter, nil
}

func (a CounterStorer) Delete(counter entities.Counter) (bool, error) {
	_, err := a.db.Exec("DELETE FROM counters WHERE ROWID = ?", counter.ID)
	if err != nil {
		return false, err
	}
	return true, nil
}
