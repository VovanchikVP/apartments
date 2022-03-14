package id_card

import (
	"apartments/cmd/web/entities"
	"database/sql"
)

type IDCardStorer struct {
	db *sql.DB
}

func New(db *sql.DB) IDCardStorer {
	return IDCardStorer{db: db}
}

func (a IDCardStorer) GetByID(id int) (idCard entities.IDCard, err error){
	var row *sql.Row

	row = a.db.QueryRow("SELECT ROWID, * FROM id_cards WHERE ROWID = ?", id)

	switch err = row.Scan(&idCard.ID, &idCard.Type, &idCard.Number, &idCard.Issued); err {
	case sql.ErrNoRows:
		return entities.IDCard{}, err
	case nil:
		return idCard, nil
	default:
		return entities.IDCard{}, err
	}
}

func (a IDCardStorer) Get(id int) ([]entities.IDCard, error) {
	var rows *sql.Rows
	var err error

	if id != 0 {
		rows, err = a.db.Query("SELECT ROWID, * FROM id_cards WHERE ROWID = ?", id)
	} else {
		rows, err = a.db.Query("SELECT ROWID, * FROM id_cards")
	}
	if rows != nil {
		defer func() {
			_ = rows.Close()
		}()
	}

	var idCard []entities.IDCard

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var a entities.IDCard
		_ = rows.Scan(&a.ID, &a.Type, &a.Number, &a.Issued)
		idCard = append(idCard, a)
	}

	return idCard, nil
}

func (a IDCardStorer) Create(idCard entities.IDCard) (entities.IDCard, error) {

	res, err := a.db.Exec("INSERT INTO id_cards(type, number, issued) VALUES (?, ?, ?)", idCard.Type, idCard.Number, idCard.Issued)

	if err != nil {
		return entities.IDCard{}, err
	}

	id, _ := res.LastInsertId()
	idCard.ID = int(id)

	return idCard, nil
}
