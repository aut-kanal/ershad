package telegram

import (
	"github.com/aryahadii/miyanbor"
	"github.com/sirupsen/logrus"
	"gitlab.com/kanalbot/ershad/configuration"
	"gitlab.com/kanalbot/ershad/mq"
)

var (
	bot *miyanbor.Bot
)

func StartBot() {
	botDebug := configuration.GetInstance().GetBool("bot.telegram.debug")
	botToken := configuration.GetInstance().GetString("bot.telegram.token")
	botSessionTimeout := configuration.GetInstance().GetInt("bot.telegram.session-timeout")
	botUpdaterTimeout := configuration.GetInstance().GetInt("bot.telegram.updater-timeout")

	var err error
	bot, err = miyanbor.NewBot(botToken, botDebug, botSessionTimeout)
	if err != nil {
		logrus.WithError(err).Fatalf("can't init bot")
	}
	logrus.Infof("telegram bot initialized completely")

	mq.SubscribeMsgs(newMessageHandler)
	logrus.Info("subscribed on msgs queue")
	logrus.Infof("===================================")

	setCallbacks(bot)
	bot.StartUpdater(0, botUpdaterTimeout)
}

func setCallbacks(bot *miyanbor.Bot) {
	bot.SetSessionStartCallbackHandler(sessionStartHandler)
	bot.SetFallbackCallbackHandler(unknownMessageHandler)

	bot.AddCallbackHandler("^a(.+)$", messageAcceptHandler)
	bot.AddCallbackHandler("^r(.+)$", messageRejectHandler)
}
