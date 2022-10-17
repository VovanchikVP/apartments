package type_pyment

import (
	"apartments/cmd/web/entities"
	"database/sql"
	"fmt"
)

type TypePymentStorer struct {
	db *sql.DB
}

func New(db *sql.DB) TypePymentStorer {
	return TypePymentStorer{db: db}
}

func (a TypePymentStorer) GetByID(id int) (payment entities.TypePayment, err error) {
	var row *sql.Row

	row = a.db.QueryRow("SELECT ROWID, * FROM type_pyments WHERE ROWID = ?", id)

	switch err = row.Scan(&payment.ID, &payment.Name); err {
	case sql.ErrNoRows:
		return entities.TypePayment{}, err
	case nil:
		return payment, nil
	default:
		return entities.TypePayment{}, err
	}
}

func (a TypePymentStorer) Get(id int) ([]entities.TypePayment, error) {
	var rows *sql.Rows
	var err error

	if id != 0 {
		rows, err = a.db.Query("SELECT ROWID, * FROM type_pyments WHERE ROWID = ?", id)
	} else {
		rows, err = a.db.Query("SELECT ROWID, * FROM type_pyments")
	}
	if rows != nil {
		defer func() {
			_ = rows.Close()
		}()
	}

	var typePyment []entities.TypePayment

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var a entities.TypePayment
		_ = rows.Scan(&a.ID, &a.Name)
		typePyment = append(typePyment, a)
	}

	return typePyment, nil
}

func (a TypePymentStorer) Create(typePyment entities.TypePayment) (entities.TypePayment, error) {

	fmt.Print("1222")
	res, err := a.db.Exec("INSERT INTO type_pyments(name) VALUES (?)", typePyment.Name)
	fmt.Print(res)
	if err != nil {
		return entities.TypePayment{}, err
	}

	id, _ := res.LastInsertId()
	typePyment.ID = int(id)

	return typePyment, nil
}

func (a TypePymentStorer) Delete(typePayment entities.TypePayment) (bool, error) {
	_, err := a.db.Exec("DELETE FROM type_pyments WHERE ROWID = ?", typePayment.ID)
	if err != nil {
		return false, err
	}
	return true, nil
}
