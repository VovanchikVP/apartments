package operation_groups

import (
	"apartments/cmd/web/entities"
	"database/sql"
)

type OperationGroupsStorer struct {
	db *sql.DB
}

func New(db *sql.DB) OperationGroupsStorer {
	return OperationGroupsStorer{db: db}
}

func (a OperationGroupsStorer) GetByID(id int) (operationGroups entities.OperationGroups, err error) {
	var row *sql.Row

	row = a.db.QueryRow("SELECT ROWID, * FROM operation_groups WHERE ROWID = ?", id)

	switch err = row.Scan(&operationGroups.ID, &operationGroups.Name); err {
	case sql.ErrNoRows:
		return entities.OperationGroups{}, err
	case nil:
		return operationGroups, nil
	default:
		return entities.OperationGroups{}, err
	}
}

func (a OperationGroupsStorer) Get(id int) (operationGroups []entities.OperationGroups, err error) {
	var rows *sql.Rows

	if id != 0 {
		rows, err = a.db.Query("SELECT ROWID, * FROM operation_groups WHERE ROWID = ?", id)
	} else {
		rows, err = a.db.Query("SELECT ROWID, * FROM operation_groups")
	}
	if rows != nil {
		defer func() {
			_ = rows.Close()
		}()
	}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var a entities.OperationGroups
		_ = rows.Scan(&a.ID, &a.Name)
		operationGroups = append(operationGroups, a)
	}

	return operationGroups, nil
}

func (a OperationGroupsStorer) Create(operationGroups entities.OperationGroups) (entities.OperationGroups, error) {

	res, err := a.db.Exec("INSERT INTO operation_groups(name) VALUES (?)", operationGroups.Name)

	if err != nil {
		return entities.OperationGroups{}, err
	}

	id, _ := res.LastInsertId()
	operationGroups.ID = int(id)

	return operationGroups, nil
}

func (a OperationGroupsStorer) Delete(operationGroups entities.OperationGroups) (bool, error) {
	_, err := a.db.Exec("DELETE FROM operation_groups WHERE ROWID = ?", operationGroups.ID)
	if err != nil {
		return false, err
	}
	return true, nil
}
