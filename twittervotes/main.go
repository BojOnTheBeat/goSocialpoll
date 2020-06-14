package main

import (
	"log"

	"github.com/bitly/go-nsq"

	"gopkg.in/mgo.v2"
)

var db *mgo.Session

func dialdb() error {
	var err error
	log.Println("dialing mongodb: localhost")
	db, err = mgo.Dial("localhost") // TODO: make this configurable
	return err
}

func closedb() {
	db.Close()
	log.Println("closed database connection")
}

type poll struct {
	Options []string
}

func loadOptions() ([]string, error) {
	var options []string
	iter := db.DB("ballots").C("polls").Find(nil).Iter()
	var p poll
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}
	iter.Close()
	return options, iter.Err()
}

func publishVotes(votes <-chan string) <-chan struct{} {
	stopchan := make(chan struct{}, 1)

	pub, err := nsq.NewProducer("localhost:4150", nsq.NewConfig())

	if err != nil {
		log.Println("Error creating an nsq producer:", err)
		log.Fatalln(err)

	}
	go func() {
		for vote := range votes {
			pub.Publish("votes", []byte(vote)) // publish vote to nsq "votes" topic
		}

		log.Println("Publisher: stopping")
		pub.Stop()
		log.Println("Publisher: stopped")
		stopchan <- struct{}{}
	}()
	return stopchan

}

func main() {

}
