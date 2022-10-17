package operation

import (
	"apartments/cmd/web/entities"
	"database/sql"
)

type OperationStorer struct {
	db *sql.DB
}

func New(db *sql.DB) OperationStorer {
	return OperationStorer{db: db}
}

func (a OperationStorer) GetByID(id int) (operation entities.Operation, err error) {
	var row *sql.Row

	row = a.db.QueryRow(`SELECT o.ROWID, o.date, o.type, o.operation_groups_id, og.name, o.value, o.proof, 
       								  o.descriptions
					           FROM operation o 
					           LEFT JOIN operation_groups og on o.operation_groups_id = og.ROWID
					           WHERE o.ROWID = ?`, id)

	switch err = row.Scan(&operation.ID, &operation.Date, &operation.Type, &operation.Group.ID,
		&operation.Group.Name, &operation.Value, &operation.Descriptions, &operation.Proof); err {
	case sql.ErrNoRows:
		return entities.Operation{}, err
	case nil:
		return operation, nil
	default:
		return entities.Operation{}, err
	}
}

func (a OperationStorer) Get(id int) (operation []entities.Operation, err error) {
	var rows *sql.Rows

	if id != 0 {
		rows, err = a.db.Query(`SELECT o.ROWID, o.date, o.type, o.operation_groups_id, og.name, o.value, o.proof,
                                             o.descriptions
					           FROM operation o 
					           LEFT JOIN operation_groups og on o.operation_groups_id = og.ROWID
					           WHERE o.ROWID = ?`, id)
	} else {
		rows, err = a.db.Query(`SELECT o.ROWID, o.date, o.type, o.operation_groups_id, og.name, o.value, o.proof, 
                                      o.descriptions
					           FROM operation o 
					           LEFT JOIN operation_groups og on o.operation_groups_id = og.ROWID`)
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
		var a entities.Operation
		_ = rows.Scan(&a.ID, &a.Date, &a.Type, &a.Group.ID, &a.Group.Name, &a.Value, &a.Proof, &a.Descriptions)
		operation = append(operation, a)
	}

	return operation, nil
}

func (a OperationStorer) Create(operation entities.Operation) (entities.Operation, error) {

	res, err := a.db.Exec(`INSERT INTO operation(date, type, operation_groups_id, value, proof, descriptions) 
						         VALUES (?, ?, ?, ?, ?, ?)`, operation.Date, operation.Type, operation.Group.ID,
		operation.Value, operation.Proof, operation.Descriptions)

	if err != nil {
		return entities.Operation{}, err
	}

	id, _ := res.LastInsertId()
	operation.ID = int(id)

	return operation, nil
}

func (a OperationStorer) Delete(operation entities.Operation) (bool, error) {
	_, err := a.db.Exec("DELETE FROM operation WHERE ROWID = ?", operation.ID)
	if err != nil {
		return false, err
	}
	return true, nil
}
