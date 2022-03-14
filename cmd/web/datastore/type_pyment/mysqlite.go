package type_pyment

import (
	"apartments/cmd/web/entities"
	"database/sql"
)

type TypePymentStorer struct {
	db *sql.DB
}

func New(db *sql.DB) TypePymentStorer {
	return TypePymentStorer{db: db}
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

	res, err := a.db.Exec("INSERT INTO type_pyments(name) VALUES (?)", typePyment.Name)

	if err != nil {
		return entities.TypePayment{}, err
	}

	id, _ := res.LastInsertId()
	typePyment.ID = int(id)

	return typePyment, nil
}
