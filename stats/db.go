package stats

import (
	r "github.com/dancannon/gorethink"
	"log"
)

type DbConfig struct {
	DbName  string
	Address string
	Tables  []string
	Session *r.Session
}

func (d *DbConfig) Setup() (*DbConfig, error) {
	db := d.DbName

	rSession, err := r.Connect(r.ConnectOpts{
		Address:  d.Address,
		Database: db,
	})
	if err != nil {
		return nil, err
	}

	d.Session = rSession
	log.Printf("Setting up database..")

	_, err = r.DBCreate(db).Run(rSession)
	if err != nil {
		log.Printf("Database already exists. Skipping..")
	}

	for _, tbl := range d.Tables {
		_, err = r.DB(db).TableCreate(tbl).Run(rSession)
		if err != nil {
			log.Printf("Table %v already exists. Skipping..", tbl)
		}
	}

	return d, nil
}
