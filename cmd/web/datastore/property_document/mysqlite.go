package property_document

import (
	"apartments/cmd/web/entities"
	"database/sql"
)

type PropertyDocumentStorer struct {
	db *sql.DB
}

func New(db *sql.DB) PropertyDocumentStorer {
	return PropertyDocumentStorer{db: db}
}

func (a PropertyDocumentStorer) GetByID(id int) (property entities.PropertyDocuments, err error){
	var row *sql.Row

	row = a.db.QueryRow("SELECT ROWID, * FROM property_documents WHERE ROWID = ?", id)

	switch err = row.Scan(&property.ID, &property.Type, &property.Number, &property.Date); err {
	case sql.ErrNoRows:
		return entities.PropertyDocuments{}, err
	case nil:
		return property, nil
	default:
		return entities.PropertyDocuments{}, err
	}
}

func (a PropertyDocumentStorer) Get(id int) ([]entities.PropertyDocuments, error) {
	var rows *sql.Rows
	var err error

	if id != 0 {
		rows, err = a.db.Query("SELECT ROWID, * FROM property_documents WHERE ROWID = ?", id)
	} else {
		rows, err = a.db.Query("SELECT ROWID, * FROM property_documents")
	}
	if rows != nil {
		defer func() {
			_ = rows.Close()
		}()
	}

	var propertyDocuments []entities.PropertyDocuments

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var a entities.PropertyDocuments
		_ = rows.Scan(&a.ID, &a.Type, &a.Number, &a.Date)
		propertyDocuments = append(propertyDocuments, a)
	}

	return propertyDocuments, nil
}

func (a PropertyDocumentStorer) Create(propertyDocuments entities.PropertyDocuments) (entities.PropertyDocuments, error) {

	res, err := a.db.Exec("INSERT INTO property_documents(type, number, date) VALUES (?, ?, ?)", propertyDocuments.Type, propertyDocuments.Number, propertyDocuments.Date)

	if err != nil {
		return entities.PropertyDocuments{}, err
	}

	id, _ := res.LastInsertId()
	propertyDocuments.ID = int(id)

	return propertyDocuments, nil
}
