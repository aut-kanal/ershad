package telegram

import (
	"strconv"

	"github.com/aryahadii/miyanbor"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/kanalbot/ershad/db"
	"gitlab.com/kanalbot/ershad/models"
	"gitlab.com/kanalbot/ershad/mq"
	telegramAPI "gopkg.in/telegram-bot-api.v4"
)

func sessionStartHandler(userSession *miyanbor.UserSession, matches []string, update interface{}) {
	logrus.WithField("user", *userSession).Debugf("new session started")
}

func unknownMessageHandler(userSession *miyanbor.UserSession, matches []string, update interface{}) {
	logrus.WithField("user", *userSession).Debugf("unknown message received")
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

func messageAcceptHandler(userSession *miyanbor.UserSession, matches []string, update interface{}) {
	// Extract message ID from matches
	id, err := strconv.Atoi(matches[1])
	if err != nil {
		logrus.WithError(err).Error("can't extract ID from matches")
		return
	}

	// Retreive message from DB
	ershadMsg := &models.ErshadMessage{}
	db.GetInstance().First(ershadMsg, matches, id)

	msg := &amqp.Publishing{
		ContentType: "application/x-binary",
		Body:        []byte(ershadMsg.EncMessage),
	}
	mq.PublishAcceptedMsg(msg)

	deleteConfig := telegramAPI.DeleteMessageConfig{
		ChatID:    userSession.ChatID,
		MessageID: (update.(*telegramAPI.Update)).CallbackQuery.Message.MessageID,
	}
	bot.Send(deleteConfig)
}

func messageRejectHandler(userSession *miyanbor.UserSession, matches []string, update interface{}) {
	// Extract message ID from matches
	id, err := strconv.Atoi(matches[1])
	if err != nil {
		logrus.WithError(err).Error("can't extract ID from matches")
		return
	}

	// Retreive message from DB
	ershadMsg := &models.ErshadMessage{}
	db.GetInstance().First(ershadMsg, matches, id)

	msg := &amqp.Publishing{
		ContentType: "application/x-binary",
		Body:        []byte(ershadMsg.EncMessage),
	}
	mq.PublishRejectedMsg(msg)

	deleteConfig := telegramAPI.DeleteMessageConfig{
		ChatID:    userSession.ChatID,
		MessageID: (update.(*telegramAPI.Update)).CallbackQuery.Message.MessageID,
	}
	bot.Send(deleteConfig)
}
