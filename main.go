package main

import (
	"encoding/json"
	"go-elastic-reindex/internal/person"
	"go-elastic-reindex/internal/pool"
	"go-elastic-reindex/internal/session"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	//logger
	l := log.New(os.Stdout, "elastic-reindex ", log.LstdFlags)

	// Create an AWS session
	sess, err := session.NewSession()
	endpoint := os.Getenv("AWS_ENDPOINT")
	sess.AwsSession.Config.Endpoint = &endpoint

	if err != nil {
		l.Fatal(err)
	}

	collector := pool.StartDispatcher(4, l) // start up worker pool
	l.Println("Logging as I work....")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		messageWorker(*sess, l, collector.Work)
	}()

	wg.Wait()
}

func messageWorker(sess session.Session, l *log.Logger, c chan *person.Person) {
	// Create the AWS Service session
	svc := sqs.New(sess.AwsSession)

	for {
		l.Println("Getting a message...")
		res, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl: aws.String(os.Getenv("QUEUE_URL")),
		})

		if err != nil {
			l.Println("failed to fetch sqs message ", err)
		}

		for _, msg := range res.Messages {
			l.Println("Sent to channel, Message ID: ", *msg.MessageId)
			p := person.Person{}
			json.Unmarshal([]byte(*msg.Body), &p)

			c <- &p
			deleteMessage(sess, *msg.ReceiptHandle, l)

		}
	}
}

func deleteMessage(sess session.Session, msgHandle string, l *log.Logger) {
	svc := sqs.New(sess.AwsSession)

	l.Println("Deleteing the message from the queue ", msgHandle)
	_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(os.Getenv("QUEUE_URL")),
		ReceiptHandle: aws.String(msgHandle),
	})

	if err != nil {
		l.Println("Error deleting message from the queue ", err)
	}
}
