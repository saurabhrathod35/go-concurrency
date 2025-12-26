package main

import (
	"database/sql"
	"encoding/gob"
	"final-project/data"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/alexedwards/scs/redisstore"

	scs "github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var webPort = "9091"

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
		Modelas:  data.New(db),
	}
	// setup mail
	app.Mailler = app.CreateMail()
	go app.listenForMail()
	// listen for signal

	go app.listenForShutdown()

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
	gob.Register(data.User{})

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

func (app *Config) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.ShoutDown()
	os.Exit(0)

}

func (app *Config) ShoutDown() {
	app.InfoLog.Println("world run clean up task...")
	// block until

	app.Wait.Wait()
	app.Mailler.DoneChan <- true

	app.InfoLog.Println("Closing Chan and application ....")
	close(app.Mailler.MailerChan)
	close(app.Mailler.ErrorChan)
	close(app.Mailler.DoneChan)

}

func (app *Config) CreateMail() Mail {
	// create channels
	errorChan := make(chan error)
	maillerChan := make(chan Message, 100)
	maillerDone := make(chan bool)

	m := Mail{
		Domain:      "localhost",
		Host:        "localhost",
		Port:        1025,
		Encryption:  "none",
		FromName:    "Info",
		FromAddress: "info@mycompany.com",
		Wait:        app.Wait,
		ErrorChan:   errorChan,
		MailerChan:  maillerChan,
		DoneChan:    maillerDone,
	}
	return m

}
