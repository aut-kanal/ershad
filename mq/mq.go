package mq

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/kanalbot/ershad/configuration"
)

var (
	conn    *amqp.Connection
	channel *amqp.Channel

	qMsgs    amqp.Queue
	qAccepts amqp.Queue
	qRejects amqp.Queue

	msgs <-chan amqp.Delivery
)

func SubscribeMsgs(callback func(amqp.Delivery)) {
	go func() {
		for msg := range msgs {
			go callback(msg)
		}
	}()
}

func PublishAcceptedMsg(data *amqp.Publishing) error {
	return channel.Publish("", qAccepts.Name, false, false, *data)
}

func PublishRejectedMsg(data *amqp.Publishing) error {
	return channel.Publish("", qRejects.Name, false, false, *data)
}

func InitMessageQueue() {
	// Connection
	var err error
	conn, err = amqp.Dial(configuration.GetInstance().GetString("rabbit-mq.url"))
	if err != nil {
		logrus.WithError(err).Fatalln("can't connect to message queue")
	}

	// Channel
	channel, err = conn.Channel()
	if err != nil {
		logrus.WithError(err).Fatalln("can't create mq channel")
	}

	// Queue
	qMsgs, err = channel.QueueDeclare(
		configuration.GetInstance().GetString("rabbit-mq.msg-queue-name"), // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		logrus.WithError(err).Fatalln("can't create messages queue")
	}

	qAccepts, err = channel.QueueDeclare(
		configuration.GetInstance().GetString("rabbit-mq.accept-queue-name"), // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		logrus.WithError(err).Fatalln("can't create accepts queue")
	}

	qRejects, err = channel.QueueDeclare(
		configuration.GetInstance().GetString("rabbit-mq.reject-queue-name"), // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		logrus.WithError(err).Fatalln("can't create rejects queue")
	}

	// Consumer
	msgs, err = channel.Consume(
		qMsgs.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		logrus.WithError(err).Fatal("can't init msg consumer")
	}

	logrus.Info("message queue initialized")
}

func Close() {
	conn.Close()
}
