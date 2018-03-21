package telegram

import (
	"github.com/aryahadii/miyanbor"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func sessionStartHandler(userSession *miyanbor.UserSession, input interface{}) {
	logrus.Debugf("new session started")
}

func unknownMessageHandler(userSession *miyanbor.UserSession, input interface{}) {
	logrus.Debugf("unknown message received, %+v", input)
}

func newMessageHandler(msg amqp.Delivery) {
	logrus.Debug("new message arrived")

	// Decode message
	decodedMsg := &Message{}
	decodeBinary(string(msg.Body), decodedMsg)

	// Send to Ershad group
	bot.Send(generateErshadMessage(decodedMsg))
}
