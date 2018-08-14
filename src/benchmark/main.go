package main

import (
	"flag"
	"time"
	"fmt"
	"gopkg.in/mgo.v2"
	"database/sql"

	"benchmark/mongo"
	"benchmark/postgres"
)

var (
	actions = map[string]func(){}
	actionsV2Mongo = map[string]func(c *mgo.Collection){}
	actionsV2Postgres = map[string]func(db *sql.DB){}
)

func init() {
	actions = make(map[string]func())
	actionsV2Mongo = make(map[string]func(c *mgo.Collection))
	actionsV2Postgres = make(map[string]func(db *sql.DB))

	actions["v1.mongo.query1"] = mongo.Query1
	actions["v1.mongo.query2"] = mongo.Query2
	actions["v1.mongo.query3"] = mongo.Query3
	actions["v1.mongo.query4"] = mongo.Query4

	actions["v1.postgres.query1"] = postgres.Query1
	actions["v1.postgres.query2"] = postgres.Query2
	actions["v1.postgres.query3"] = postgres.Query3
	actions["v1.postgres.query4"] = postgres.Query4

	actionsV2Mongo["query1"] = mongo.ExecQuery1
	actionsV2Mongo["query2"] = mongo.ExecQuery2
	actionsV2Mongo["query3"] = mongo.ExecQuery3
	actionsV2Mongo["query4"] = mongo.ExecQuery4

	actionsV2Postgres["query1"] = postgres.ExecQuery1
	actionsV2Postgres["query2"] = postgres.ExecQuery2
	actionsV2Postgres["query3"] = postgres.ExecQuery3
	actionsV2Postgres["query4"] = postgres.ExecQuery4
}

func main() {
	version := "v1"
	flag.StringVar(&version, "v", "v1", "benchmark version: v1 or v2")

	dbName := "mongo"
	flag.StringVar(&dbName, "db", "mongo", "database name: mongo or postgres")

	action := "query1"
	flag.StringVar(&action, "q", "query1", "query number")

	flag.Parse()

	switch version {
	case "v1":
		runV1(version, dbName, action)
	case "v2":
		runV2(dbName, action)
	}
}

func runWithTime(actionName string, f func()) {
	startedAt := time.Now().UnixNano()
	f()
	finishedAt := time.Now().UnixNano()

	fmt.Printf("Action: %s, took: %d microseconds \n\n", actionName, (finishedAt - startedAt) / 1000)
}

func runV1(version string, dbName string, action string) {
	actionName := version+"."+dbName+"."+action
	runWithTime(actionName, func() {
		actions[actionName]()
	})
}

func runV2(dbName string, action string) {
	actionName := "v2."+dbName+"."+action
	switch dbName {
	case "mongo":
		session, c := mongo.GetDBC()
		defer session.Close()
		runWithTime(actionName, func() {
			actionsV2Mongo[action](c)
		})

	case "postgres":
		db := postgres.GetDBC()
		defer db.Close()
		runWithTime(actionName, func() {
			actionsV2Postgres[action](db)
		})
	}
}
