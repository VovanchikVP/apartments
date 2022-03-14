package indication

import (
	"apartments/cmd/web/datastore/counter"
	"apartments/cmd/web/driver"
	"apartments/cmd/web/entities"
	"database/sql"
	"log"
)

type IndicationStorer struct {
	db *sql.DB
}

func New(db *sql.DB) IndicationStorer {
	return IndicationStorer{db: db}
}

func (a IndicationStorer) GetByID(id int) (indication entities.Indication, err error){
	var row *sql.Row

	row = a.db.QueryRow("SELECT ROWID, * FROM indications WHERE ROWID = ?", id)

	db, err := driver.ConnectToMySQL()
	if err != nil {
		log.Println("could not connect to sql, err:", err)
		return entities.Indication{}, err
	}

	counterDB := counter.New(db)

	switch err = row.Scan(&indication.ID, &indication.Counter.ID, &indication.Date, &indication.Data); err {
	case sql.ErrNoRows:
		return entities.Indication{}, err
	case nil:
		indication.Counter, _ = counterDB.GetByID(indication.Counter.ID)
		return indication, nil
	default:
		return entities.Indication{}, err
	}
}

func (a IndicationStorer) Get(id int) (indication []entities.Indication, err error) {
	var rows *sql.Rows

	if id != 0 {
		rows, err = a.db.Query("SELECT ROWID, * FROM indications WHERE ROWID = ?", id)
	} else {
		rows, err = a.db.Query("SELECT ROWID, * FROM indications")
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

	counterDB := counter.New(db)

	for rows.Next() {
		var a entities.Indication
		_ = rows.Scan(&a.ID, &a.Counter.ID, &a.Date, &a.Data)
		a.Counter, _ = counterDB.GetByID(a.Counter.ID)
		indication = append(indication, a)
	}

	return indication, nil
}

func (a IndicationStorer) Create(indication entities.Indication) (entities.Indication, error) {

	res, err := a.db.Exec("INSERT INTO indications(counter_id, date, data) VALUES (?, ?, ?)", indication.Counter.ID, indication.Date, indication.Data)

	if err != nil {
		return entities.Indication{}, err
	}

	id, _ := res.LastInsertId()
	indication.ID = int(id)

	return indication, nil
}
