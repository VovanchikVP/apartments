package address

import (
	"apartments/cmd/web/entities"
	"database/sql"
)

type AddressStorer struct {
	db *sql.DB
}

func New(db *sql.DB) AddressStorer {
	return AddressStorer{db: db}
}

func (a AddressStorer) GetByID(id int) (address entities.Address, err error){
	var row *sql.Row

	row = a.db.QueryRow("SELECT ROWID, * FROM address WHERE ROWID = ?", id)

	switch err = row.Scan(&address.ID, &address.Index, &address.City, &address.Street, &address.House, &address.Apartment); err {
	case sql.ErrNoRows:
		return entities.Address{}, err
	case nil:
		return address, nil
	default:
		return entities.Address{}, err
	}
}

func (a AddressStorer) Get(id int) ([]entities.Address, error) {
	var rows *sql.Rows
	var err error

	if id != 0 {
		rows, err = a.db.Query("SELECT ROWID, * FROM address WHERE ROWID = ?", id)
	} else {
		rows, err = a.db.Query("SELECT ROWID, * FROM address")
	}
	if rows != nil {
		defer func() {
			_ = rows.Close()
		}()
	}

	var address []entities.Address

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var a entities.Address
		_ = rows.Scan(&a.ID, &a.Index, &a.City, &a.Street, &a.House, &a.Apartment)
		address = append(address, a)
	}

	return address, nil
}

func (a AddressStorer) Create(address entities.Address) (entities.Address, error) {

	res, err := a.db.Exec("INSERT INTO address(post_index, city, street, house, apartment) VALUES (?, ?, ?, ?, ?)", address.Index, address.City, address.Street, address.House, address.Apartment)

	if err != nil {
		return entities.Address{}, err
	}

	id, _ := res.LastInsertId()
	address.ID = int(id)

	return address, nil
}

func (a AddressStorer) Delete(address entities.Address) (bool, error)  {
	_, err := a.db.Exec("DELETE FROM address WHERE ROWID = ?", address.ID)
	if err != nil {
		return false, err
	}
	return true, nil
}

