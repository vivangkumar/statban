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

func (d *Db) SummarizeByHour(ghConfig *GithubConfig) {
	now := time.Now().UTC()
	year, month, day := now.Year(), now.Month(), now.Day()
	thisHour := time.Date(year, month, day, now.Hour(), 0, 0, 0, time.UTC)
	nextHour := thisHour.Add(time.Duration(1) * time.Hour)

	log.Printf("Summarizing by hour between %v and %v", thisHour, nextHour)

	cur, err := r.DB(d.Name).Table("hourly_state").Filter(r.Row.Field("created_at").
		During(thisHour, nextHour)).Group("label").Count().Run(d.Session)
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

	sumBatch := &SummarizedHour{
		HourStart: thisHour,
		HourEnd:   nextHour,
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

	d.writeHourlySummary(sumBatch)
	return
}

func (d *Db) SummarizeByDay() {
	now := time.Now().UTC()
	year, month, day := now.Year(), now.Month(), now.Day()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	tomorrow := today.Add(time.Duration(24) * time.Hour)

	log.Printf("Summarizing by day between %v and %v", today, tomorrow)

	cur, err := r.DB(d.Name).Table("hourly_summary").Filter(r.Row.Field("created_at").
		During(today, tomorrow)).Run(d.Session)
	if err != nil {
		log.Printf("Error getting day summary %v", err.Error())
		return
	}
	defer cur.Close()

	var res []SummarizedHour
	err = cur.All(&res)
	if err != nil {
		log.Printf("Error when summarizing by day: %v", err.Error())
		return
	}

	sb := &res[len(res)-1]
	d.writeDailySummary(NewSummarizedDay(sb, today, tomorrow))
	return
}

func (d *Db) GetDailyStats() (res []SummarizedDay, err error) {
	cur, err := r.DB(d.Name).Table("daily_summary").Limit(30).Run(d.Session)
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

func (d *Db) GetHourlyStats() (res []SummarizedHour, err error) {
	now := time.Now().UTC()
	year, month, day := now.Year(), now.Month(), now.Day()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	tomorrow := today.Add(time.Duration(24) * time.Hour)

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

func (d *Db) writeHourlySummary(summary *SummarizedHour) {
	_, err := r.DB(d.Name).Table("hourly_summary").Insert(*summary).RunWrite(d.Session)
	if err != nil {
		log.Printf("Error inserting summary %v into table: %v", summary, err.Error())
		return
	}
	return
}

func (d *Db) writeDailySummary(ds *SummarizedDay) {
	_, err := r.DB(d.Name).Table("daily_summary").Insert(ds).RunWrite(d.Session)
	if err != nil {
		log.Printf("Error inserting day summary %v into table: %v", ds, err.Error())
		return
	}
	return
}
