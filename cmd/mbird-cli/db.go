package main

import (
	"encoding/json"
	"github.com/sdomino/scribble"
	"math/rand"
	"time"
)

type Data struct {
	Id      string              `json:"id"`
	Headers map[string][]string `json:"headers"`
	Status  int                 `json:"status_code"`
	Body    string              `json:"body"`
}

type DB struct {
	db *scribble.Driver
}

func newDB(dir string) (*DB, error) {
	db, err := scribble.New(dir, nil)
	d := &DB{db: db}
	return d, err
}

func (self *DB) record(collection string, data Data) error {
	if err := self.db.Write(collection, data.Id, data); err != nil {
		return err
	}
	return nil
}

func (self *DB) replayRandom(collection string) (error, *Data) {
	records, err := self.db.ReadAll(collection)
	if err != nil {
		return err, nil
	}
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	record := records[rand.Intn(len(records))]
	var randData Data
	if err := json.Unmarshal([]byte(record), &randData); err != nil {
		return err, nil
	}
	return nil, &randData
}
