package main

import (
	"flag"
	"./logsend"
	"os"
	logpkg "log"
)

var (
	logDir  = flag.String("log-dir", "./tmp", "log directories")
	config = flag.String("config", "config.json", "path to config.json file")
	dbhost = flag.String("db-host", "localhost:8086", "db host")
	dbuser = flag.String("db-user", "root", "db user")
	dbpassword = flag.String("db-password", "root", "db-password")
	database = flag.String("database", "test1", "database")
)

var (
	log = logpkg.New(os.Stderr, "", logpkg.Lmicroseconds)
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	config, err := logsend.LoadConfig(*config)
	if err != nil {
		log.Fatalf("can't load config %+v", err)
	}

	logsScopes := make([]*logsend.LogScope, 0)

	for _, group := range config.Groups {
		lsc := logsend.NewLogScope(&group)
		logsScopes = append(logsScopes, lsc)
	}

	logsend.AssociatedLogPerFile(*logDir, logsScopes)
	dbClient, err := logsend.NewDBClient(*dbhost, *dbuser, *dbpassword, *database)
	if err != nil {
		log.Fatalf("NewDBClient %+v", err)
	}

	for _,lsc := range logsScopes {
		go lsc.Tailing(dbClient)
	}

	select{}
}