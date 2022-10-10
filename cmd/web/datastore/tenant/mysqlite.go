package tenant

import (
	"apartments/cmd/web/datastore/contract_rent"
	"apartments/cmd/web/datastore/person"
	"apartments/cmd/web/driver"
	"apartments/cmd/web/entities"
	"database/sql"
	"log"
)

type TenantStorer struct {
	db *sql.DB
}

func New(db *sql.DB) TenantStorer {
	return TenantStorer{db: db}
}

func (a TenantStorer) GetByID(id int) (tenant entities.Tenant, err error){
	var row *sql.Row

	row = a.db.QueryRow("SELECT ROWID, * FROM tenant WHERE ROWID = ?", id)

	db, err := driver.ConnectToMySQL()
	if err != nil {
		log.Println("could not connect to sql, err:", err)
		return entities.Tenant{}, err
	}

	contractDB := contract_rent.New(db)
	personDB := person.New(db)

	switch err = row.Scan(&tenant.ID, &tenant.ContractRent.ID, &tenant.Person.ID); err {
	case sql.ErrNoRows:
		return entities.Tenant{}, err
	case nil:
		tenant.ContractRent, _ = contractDB.GetByID(tenant.ContractRent.ID)
		tenant.Person, _ = personDB.GetByID(tenant.Person.ID)
		return tenant, nil
	default:
		return entities.Tenant{}, err
	}
}

func (a TenantStorer) Get(id int) (tenant []entities.Tenant, err error) {
	var rows *sql.Rows

	if id != 0 {
		rows, err = a.db.Query("SELECT * FROM tenant WHERE ROWID = ?", id)
	} else {
		rows, err = a.db.Query("SELECT * FROM tenant")
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

	contractDB := contract_rent.New(db)
	personDB := person.New(db)

	for rows.Next() {
		var a entities.Tenant
		_ = rows.Scan(&a.ID, &a.ContractRent.ID, &a.Person.ID)
		a.ContractRent, _ = contractDB.GetByID(a.ContractRent.ID)
		a.Person, _ = personDB.GetByID(a.Person.ID)
		tenant = append(tenant, a)
	}

	return tenant, nil
}

func (a TenantStorer) Create(tenant entities.Tenant) (entities.Tenant, error) {
	res, err := a.db.Exec("INSERT INTO tenant(contract_rent_id, person_id) VALUES (?, ?)", tenant.ContractRent.ID, tenant.Person.ID)

	if err != nil {
		return entities.Tenant{}, err
	}

	id, _ := res.LastInsertId()
	tenant.ID = int(id)

	return tenant, nil
}
