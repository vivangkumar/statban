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

var (
	today    time.Time
	tomorrow time.Time
)

func init() {
	today = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	tomorrow = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 0, 0, 0, 0, time.UTC)
}

func (d *Db) Setup() (*Db, error) {
	db := d.Name
	tables := []string{"hourly_state", "daily_summary", "hourly_summary"}

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
		log.Printf("Creating table %v", tbl)
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

func (d *Db) SummarizeByBatch(batchId string, ghConfig *GithubConfig) {
	log.Printf("Summarizing for batch id: %v", batchId)

	cur, err := r.DB(d.Name).Table("hourly_state").
		Filter(map[string]string{"batch_id": batchId}).
		Group("label").Count().Run(d.Session)
	if err != nil {
		log.Printf("Error when grouping data: %v", err.Error())
		return
	}
	defer cur.Close()

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

	labels := ghConfig.Labels
	var addedLabels []string

	ss := make([]SummarizedState, len(res))
	for i, sg := range res {
		addedLabels = append(addedLabels, sg.Group)
		ss[i] = NewSummarizedState(sg.Group, sg.Reduction)
	}

	missingLabels := arrayDifference(addedLabels, labels)
	addMissingLabels(missingLabels, &ss)
	sumBatch.States = &ss

	d.writeBatchSummary(sumBatch)
	return
}

func (d *Db) SummarizeByDay() {
	log.Printf("Summarizing by day between %v and %v", today, tomorrow)

	cur, err := r.DB(d.Name).Table("hourly_summary").Filter(r.Row.Field("created_at").
		During(today, tomorrow)).Run(d.Session)
	if err != nil {
		log.Printf("Error getting day summary %v", err.Error())
		return
	}
	defer cur.Close()

	var res []SummarizedBatch
	err = cur.All(&res)
	if err != nil {
		log.Printf("Error when summarizing by day: %v", err.Error())
		return
	}

	sb := &res[len(res)-1]
	d.writeDaySummary(NewSummarizedDay(sb, today, tomorrow))
	return
}

func (d *Db) GetDailyStats() (res []SummarizedDay, err error) {
	cur, err := r.DB(d.Name).Table("daily_summary").Filter(r.Row.Field("beginning").
		Eq(today)).Run(d.Session)
	if err != nil {
		log.Printf("Error reading day summary: %v", err.Error())
		return nil, err
	}
	defer cur.Close()

	err = cur.All(&res)
	if err != nil {
		log.Printf("Error when decoding into struct: %v", err.Error())
		return nil, err
	}

	return res, nil
}

func (d *Db) GetBatchStats() (res []SummarizedBatch, err error) {
	cur, err := r.DB(d.Name).Table("hourly_summary").Filter(r.Row.Field("created_at").
		During(today, tomorrow)).Run(d.Session)
	if err != nil {
		log.Printf("Error reading batch summary: %v", err.Error())
		return nil, err
	}
	defer cur.Close()

	err = cur.All(&res)
	if err != nil {
		log.Printf("Error when decoding into struct: %v", err.Error())
		return nil, err
	}

	return res, nil
}

func (d *Db) writeBatchSummary(summary *SummarizedBatch) {
	_, err := r.DB(d.Name).Table("hourly_summary").Insert(*summary).RunWrite(d.Session)
	if err != nil {
		log.Printf("Error inserting summary %v into table: %v", summary, err.Error())
		return
	}
	return
}

func (d *Db) writeDaySummary(ds *SummarizedDay) {
	_, err := r.DB(d.Name).Table("daily_summary").Insert(ds).RunWrite(d.Session)
	if err != nil {
		log.Printf("Error inserting day summary %v into table: %v", ds, err.Error())
		return
	}
	return
}
