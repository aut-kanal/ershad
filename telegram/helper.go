package telegram

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"io/ioutil"
	"net/http"

	"gitlab.com/kanalbot/ershad/models"
	"gitlab.com/kanalbot/ershad/ui/keyboard"

	"gitlab.com/kanalbot/ershad/configuration"
	"gitlab.com/kanalbot/ershad/ui/text"

	"github.com/sirupsen/logrus"

	telegramAPI "gopkg.in/telegram-bot-api.v4"
)

func decodeBinary(enc string, out interface{}) {
	b64, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		logrus.WithError(err).Error("base64 decode failed")
	}
	buf := bytes.Buffer{}
	buf.Write(b64)
	dec := gob.NewDecoder(&buf)
	err = dec.Decode(out)
	if err != nil {
		logrus.WithError(err).Error("can't decode message")
	}
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func generateErshadMessage(msg *models.Message, id string) telegramAPI.Chattable {
	chatID := configuration.GetInstance().GetInt64("ershad-chatid")

	// Send messages with media attached
	if msg.FileURL != "" {
		// Download attached media
		fileBytes, err := downloadFile(msg.FileURL)
		if err != nil {
			return telegramAPI.NewMessage(chatID, text.MsgCannotDownloadFile)
		}
		mediaFile := telegramAPI.FileBytes{
			Bytes: fileBytes,
		}

		// Create message
		if msg.Audio != nil {
			mediaFile.Name = msg.Audio.Title
			audio := telegramAPI.NewAudioUpload(chatID, mediaFile)
			audio.Caption = msg.Caption
			audio.ReplyMarkup = keyboard.NewErshadInlineKeyboard(id)
			return audio
		}
		if msg.Voice != nil {
			mediaFile.Name = "voice"
			voice := telegramAPI.NewVoiceUpload(chatID, mediaFile)
			voice.Caption = msg.Caption
			voice.ReplyMarkup = keyboard.NewErshadInlineKeyboard(id)
			return voice
		}
		if msg.Video != nil {
			mediaFile.Name = "video"
			video := telegramAPI.NewAudioUpload(chatID, mediaFile)
			video.Caption = msg.Caption
			video.ReplyMarkup = keyboard.NewErshadInlineKeyboard(id)
			return video
		}
		if msg.Document != nil {
			mediaFile.Name = msg.Document.FileName
			document := telegramAPI.NewDocumentUpload(chatID, mediaFile)
			document.Caption = msg.Caption
			document.ReplyMarkup = keyboard.NewErshadInlineKeyboard(id)
			return document
		}
		if msg.Photo != nil {
			mediaFile.Name = "photo"
			photo := telegramAPI.NewPhotoUpload(chatID, mediaFile)
			photo.Caption = msg.Caption
			photo.ReplyMarkup = keyboard.NewErshadInlineKeyboard(id)
			return photo
		}
	}

	// Message without media
	replyMsg := telegramAPI.NewMessage(chatID, msg.Text)
	replyMsg.ReplyMarkup = keyboard.NewErshadInlineKeyboard(id)
	return replyMsg
}
