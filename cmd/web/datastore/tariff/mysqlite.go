package tariff

import (
	"apartments/cmd/web/datastore/counter"
	"apartments/cmd/web/driver"
	"apartments/cmd/web/entities"
	"database/sql"
	"log"
)

type TariffStorer struct {
	db *sql.DB
}

func New(db *sql.DB) TariffStorer {
	return TariffStorer{db: db}
}

func (a TariffStorer) GetByID(id int) (tariff entities.Tariff, err error){
	var row *sql.Row

	row = a.db.QueryRow("SELECT ROWID, * FROM tariffs WHERE ROWID = ?", id)

	db, err := driver.ConnectToMySQL()
	if err != nil {
		log.Println("could not connect to sql, err:", err)
		return entities.Tariff{}, err
	}

	counterDB := counter.New(db)

	switch err = row.Scan(&tariff.ID, &tariff.Counter.ID, &tariff.SetDate, &tariff.Cost); err {
	case sql.ErrNoRows:
		return entities.Tariff{}, err
	case nil:
		tariff.Counter, _ = counterDB.GetByID(tariff.Counter.ID)
		return tariff, nil
	default:
		return entities.Tariff{}, err
	}
}

func (a TariffStorer) Get(id int) (tariff []entities.Tariff, err error) {
	var rows *sql.Rows

	if id != 0 {
		rows, err = a.db.Query("SELECT ROWID, * FROM tariffs WHERE ROWID = ?", id)
	} else {
		rows, err = a.db.Query("SELECT ROWID, * FROM tariffs")
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
		var a entities.Tariff
		_ = rows.Scan(&a.ID, &a.Counter.ID, &a.SetDate, &a.Cost)
		a.Counter, _ = counterDB.GetByID(a.Counter.ID)
		tariff = append(tariff, a)
	}

	return tariff, nil
}

func (a TariffStorer) Create(tariff entities.Tariff) (entities.Tariff, error) {

	res, err := a.db.Exec("INSERT INTO tariffs(counter_id, set_date, cost) VALUES (?, ?, ?)", tariff.Counter.ID, tariff.SetDate, tariff.Cost)

	if err != nil {
		return entities.Tariff{}, err
	}

	id, _ := res.LastInsertId()
	tariff.ID = int(id)

	return tariff, nil
}
