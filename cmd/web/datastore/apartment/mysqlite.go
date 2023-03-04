package apartment

import (
	"apartments/cmd/web/datastore/address"
	"apartments/cmd/web/datastore/property_document"
	"apartments/cmd/web/entities"
	"database/sql"
)

type ApartmentStorer struct {
	db *sql.DB
}

func New(db *sql.DB) ApartmentStorer {
	return ApartmentStorer{db: db}
}

func (a ApartmentStorer) GetByID(id int) (apartment entities.Apartment, err error) {
	var row *sql.Row

	row = a.db.QueryRow("SELECT ROWID, * FROM apartments WHERE ROWID = ?", id)

	addressDB := address.New(a.db)
	propertyDB := property_document.New(a.db)

	switch err = row.Scan(&apartment.ID, &apartment.Address.ID, &apartment.CountRooms, &apartment.PropertyDocuments.ID, &apartment.Rent); err {
	case sql.ErrNoRows:
		return entities.Apartment{}, err
	case nil:
		apartment.Address, _ = addressDB.GetByID(apartment.Address.ID)
		apartment.PropertyDocuments, _ = propertyDB.GetByID(apartment.PropertyDocuments.ID)
		return apartment, nil
	default:
		return entities.Apartment{}, err
	}
}

func (a ApartmentStorer) Get(id int) (apartments []entities.Apartment, err error) {
	var rows *sql.Rows

	if id != 0 {
		rows, err = a.db.Query("SELECT ROWID, * FROM apartments WHERE ROWID = ?", id)
	} else {
		rows, err = a.db.Query("SELECT ROWID, * FROM apartments")
	}
	if rows != nil {
		defer func() {
			_ = rows.Close()
		}()
	}

	if err != nil {
		return nil, err
	}

	addressDB := address.New(a.db)
	propertyDB := property_document.New(a.db)

	for rows.Next() {
		var a entities.Apartment
		_ = rows.Scan(&a.ID, &a.Address.ID, &a.CountRooms, &a.PropertyDocuments.ID, &a.Rent)
		a.Address, _ = addressDB.GetByID(a.Address.ID)
		a.PropertyDocuments, _ = propertyDB.GetByID(a.PropertyDocuments.ID)
		apartments = append(apartments, a)
	}

	return apartments, nil
}

func (a ApartmentStorer) Create(apartment entities.Apartment) (entities.Apartment, error) {

	res, err := a.db.Exec(`INSERT INTO apartments(address_id, count_rooms, property_document_id, rent) 
								 VALUES (?, ?, ?, ?)`, apartment.Address.ID, apartment.CountRooms,
		apartment.PropertyDocuments.ID, apartment.Rent)

	if err != nil {
		return entities.Apartment{}, err
	}

	id, _ := res.LastInsertId()
	apartment.ID = int(id)

	return apartment, nil
}

func (a ApartmentStorer) Delete(apartment entities.Apartment) (bool, error) {
	_, err := a.db.Exec("DELETE FROM apartments WHERE ROWID = ?", apartment.ID)
	if err != nil {
		return false, err
	}
	return true, nil
}
