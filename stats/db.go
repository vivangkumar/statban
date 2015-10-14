package stats

import (
	r "github.com/dancannon/gorethink"
	"log"
	"time"
)

type Db struct {
	Name    string
	Address string
	Session *r.Session
}

type SummarizedGroup struct {
	Group     string
	Reduction int
}

type SummarizedBatch struct {
	BatchId   string             `gorethink:"batch_id,omitempty" json:"batch_id"`
	States    []*SummarizedState `gorethink:"states",omitempty json:"states"`
	CreatedAt time.Time          `gorethink:"created_at" json:"created_at"`
}

type SummarizedState struct {
	Label string `gorethink:"label,omitempty" json:"label"`
	Count int    `gorethink:"count",omitempty" json:"count" `
}

func (d *Db) Setup() (*Db, error) {
	db := d.Name
	tables := []string{"hourly_state", "daily_state", "hourly_summary"}

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
	for _, tbl := range tables {
		_, err = r.DB(db).TableCreate(tbl).Run(rSession)
		if err != nil {
			log.Printf("Table %v already exists. Skipping..", tbl)
		}
	}

	return d, nil
}

func (d *Db) StoreHourlyState(issues []*StatbanIssue) {
	for _, issue := range issues {
		_, err := r.DB(d.Name).Table("hourly_state").Insert(issue).RunWrite(d.Session)
		if err != nil {
			log.Printf("Error inserting issue %v into table", issue)
		}
	}
}

// TODO link to issues and account for missing states
func (d *Db) SummarizeByBatch(batchId string) {
	cur, err := r.DB(d.Name).Table("hourly_state").Group("label").Count().Run(d.Session)
	if err != nil {
		log.Printf("Error when grouping data: %v", err.Error())
		return
	}

	var res []SummarizedGroup
	err = cur.All(&res)
	if err != nil {
		log.Printf("Error when summarizing: %v", err.Error())
		return
	}

	sumBatch := &SummarizedBatch{
		BatchId:   batchId,
		CreatedAt: time.Now(),
	}

	st := make([]*SummarizedState, len(res))
	for i, sg := range res {
		st[i] = &SummarizedState{
			Label: sg.Group,
			Count: sg.Reduction,
		}
	}

	sumBatch.States = st
	d.writeSummaries(sumBatch)
}

func (d *Db) writeSummaries(summary *SummarizedBatch) {
	_, err := r.DB(d.Name).Table("hourly_summary").Insert(*summary).RunWrite(d.Session)
	if err != nil {
		log.Printf("Error inserting summary %v into table", summary)
	}
}
