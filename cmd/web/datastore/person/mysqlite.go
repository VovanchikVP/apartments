package person

import (
	"apartments/cmd/web/datastore/address"
	"apartments/cmd/web/datastore/id_card"
	"apartments/cmd/web/driver"
	"apartments/cmd/web/entities"
	"database/sql"
	"log"
)

type PersonStorer struct {
	db *sql.DB
}

func New(db *sql.DB) PersonStorer {
	return PersonStorer{db: db}
}

func (a PersonStorer) GetByID(id int) (person entities.Person, err error) {
	var row *sql.Row

	row = a.db.QueryRow("SELECT ROWID, * FROM persons WHERE ROWID = ?", id)

	db, err := driver.ConnectToMySQL()
	if err != nil {
		log.Println("could not connect to sql, err:", err)
		return entities.Person{}, err
	}

	idCardDB := id_card.New(db)
	addressDB := address.New(db)

	switch err = row.Scan(&person.ID, &person.LastName, &person.FirstName, &person.Patronymic, &person.IDCard.ID, &person.Phone, &person.Address.ID); err {
	case sql.ErrNoRows:
		return entities.Person{}, err
	case nil:
		person.IDCard, _ = idCardDB.GetByID(person.IDCard.ID)
		person.Address, _ = addressDB.GetByID(person.Address.ID)
		return person, nil
	default:
		return entities.Person{}, err
	}
}

func (a PersonStorer) Get(id int) ([]entities.Person, error) {
	var rows *sql.Rows
	var err error

	if id != 0 {
		rows, err = a.db.Query("SELECT ROWID, * FROM persons WHERE ROWID = ?", id)
	} else {
		rows, err = a.db.Query("SELECT ROWID, * FROM persons")
	}
	if rows != nil {
		defer func() {
			_ = rows.Close()
		}()
	}

	var person []entities.Person

	if err != nil {
		return nil, err
	}

	db, err := driver.ConnectToMySQL()
	if err != nil {
		log.Println("could not connect to sql, err:", err)
		return nil, err
	}

	idCardDB := id_card.New(db)
	addressDB := address.New(db)

	for rows.Next() {
		var a entities.Person
		_ = rows.Scan(&a.ID, &a.LastName, &a.FirstName, &a.Patronymic, &a.IDCard.ID, &a.Phone, &a.Address.ID)
		a.IDCard, _ = idCardDB.GetByID(a.IDCard.ID)
		a.Address, _ = addressDB.GetByID(a.Address.ID)
		person = append(person, a)
	}

	return person, nil
}

func (a PersonStorer) Create(person entities.Person) (entities.Person, error) {

	res, err := a.db.Exec("INSERT INTO persons(last_name, first_name, patronymic, id_card_id, phone, address_id) VALUES (?, ?, ?, ?, ?, ?)", person.LastName, person.FirstName, person.Patronymic, person.IDCard.ID, person.Phone, person.Address.ID)

	if err != nil {
		return entities.Person{}, err
	}

	id, _ := res.LastInsertId()
	person.ID = int(id)

	return person, nil
}

func (a PersonStorer) Delete(person entities.Person) (bool, error) {
	_, err := a.db.Exec("DELETE FROM persons WHERE ROWID = ?", person.ID)
	if err != nil {
		return false, err
	}
	return true, nil
}
