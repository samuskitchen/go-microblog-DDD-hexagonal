package main

import (
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
	"microblog/infrastructure"
	data "microblog/infrastructure/database"
	"os"
	"os/signal"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("error in main")

			os.Exit(1)
		}
	}()

	// connection to the database.
	db := data.New()
	if err := db.DB.Ping(); err != nil {
		log.Fatal(err)
	}

	conn := &data.Data{
		DB: db.DB,
	}

	DaemonPort := os.Getenv("DAEMON_PORT")
	serv := infrastructure.NewApplication(DaemonPort, conn)

	// start the server.
	go serv.Start()

	// Wait for an in interrupt.
	// If you ask about <- look here https://tour.golang.org/concurrency/2
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Attempt a graceful shutdown.
	_ = serv.Close()
	_ = data.Close()
}
