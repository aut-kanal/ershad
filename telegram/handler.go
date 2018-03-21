package telegram

import (
	"strconv"

	"github.com/aryahadii/miyanbor"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/kanalbot/ershad/db"
	"gitlab.com/kanalbot/ershad/models"
	"gitlab.com/kanalbot/ershad/mq"
)

func sessionStartHandler(userSession *miyanbor.UserSession, input interface{}) {
	logrus.WithField("user", *userSession).Debugf("new session started")
}

func unknownMessageHandler(userSession *miyanbor.UserSession, input interface{}) {
	logrus.WithField("user", *userSession).Debugf("unknown message received, %+v", input)
}

func newMessageHandler(msg amqp.Delivery) {
	logrus.Debug("new message arrived")

	// Decode message
	decodedMsg := &models.Message{}
	decodeBinary(string(msg.Body), decodedMsg)

	// Save encMessage in DB
	encMsg := models.ErshadMessage{
		EncMessage: string(msg.Body),
	}
	db.GetInstance().Create(&encMsg)

	// Send to Ershad group
	bot.Send(generateErshadMessage(decodedMsg, strconv.FormatUint(uint64(encMsg.ID), 10)))
}

func messageAcceptHandler(userSession *miyanbor.UserSession, matches interface{}) {
	// Extract message ID from matches
	input := matches.([]string)
	id, err := strconv.Atoi(input[1])
	if err != nil {
		logrus.WithError(err).Error("can't extract ID from matches")
		return
	}

	// Retreive message from DB
	ershadMsg := &models.ErshadMessage{}
	db.GetInstance().First(ershadMsg, matches.([]string), id)

	msg := &amqp.Publishing{
		ContentType: "application/x-binary",
		Body:        []byte(ershadMsg.EncMessage),
	}
	mq.PublishAcceptedMsg(msg)
}

func messageRejectHandler(userSession *miyanbor.UserSession, matches interface{}) {
	// Extract message ID from matches
	input := matches.([]string)
	id, err := strconv.Atoi(input[1])
	if err != nil {
		logrus.WithError(err).Error("can't extract ID from matches")
		return
	}

	// Retreive message from DB
	ershadMsg := &models.ErshadMessage{}
	db.GetInstance().First(ershadMsg, matches.([]string), id)

	msg := &amqp.Publishing{
		ContentType: "application/x-binary",
		Body:        []byte(ershadMsg.EncMessage),
	}
	mq.PublishRejectedMsg(msg)
}
