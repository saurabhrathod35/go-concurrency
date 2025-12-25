package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/redisstore"

	scs "github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var webPort = "9090"

func main() {

	// connect to the database postgres
	db := initDB()
	db.Ping()

	// create session
	session := initSession()

	// create logger
	infoLogs := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLogs := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	// create chennels

	// create wait group
	wg := &sync.WaitGroup{}

	// setup application config
	app := Config{
		Session:  session,
		DB:       db,
		InfoLog:  infoLogs,
		ErrorLog: errorLogs,
		Wait:     wg,
	}
	// setup mail
	app.searve()
	// listen for web requests

}
func (app *Config) searve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	app.InfoLog.Println("Starting web server...")
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
func initDB() *sql.DB {
	// connect to the database
	conn := connectToDB()
	if conn == nil {
		log.Panic("Cannot connect to database")
		return nil
	}

	return conn
}

func connectToDB() *sql.DB {
	// Placeholder for actual database connection logic
	count := 0
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			count++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}
		if count > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Backing off for 2 seconds...")
		// wait for 2 seconds before trying again
		time.Sleep(2 * time.Second)
	}

}

func openDB(dsn string) (*sql.DB, error) {
	// Placeholder for actual database opening logic
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initSession() *scs.SessionManager {

	// setup session
	session := scs.New()
	session.Store = redisstore.New(initRedis())

	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true
	return session
}

func initRedis() *redis.Pool {
	// Placeholder for actual Redis initialization logic
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS"))
		},
	}
	return redisPool
}
