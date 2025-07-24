package linebot

import (
	"strings"

	"github.com/inkxk/bad-bot/app"
	"github.com/line/line-bot-sdk-go/v8/linebot"
	"go.uber.org/zap"
)

type Handler struct {
	Bot    *linebot.Client
	Logger *zap.Logger
}

func NewHandler(bot *linebot.Client, logger *zap.Logger) *Handler {
	return &Handler{
		Bot:    bot,
		Logger: logger,
	}
}

func (h *Handler) Callback(ctx app.Context) {
	events, err := h.Bot.ParseRequest(ctx.Request())
	if err != nil {
		h.Logger.Sugar().Errorf("Error parsing LINE webhook request: %v", err)
		ctx.ErrorResponse(err)
	}

	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeJoin:
			_, err := h.Bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage("สวัสดีครับ กลุ่มนี้ใครจะตีแบดบ้าง 🏸"),
			).Do()
			if err != nil {
				h.Logger.Sugar().Errorf("Reply to join event error: %v", err)
			}
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				text := strings.TrimSpace(message.Text)

				// check prefix !
				if strings.HasPrefix(text, "!") {
					switch text {
					case "!ตีแบดกัน":
						_, err = h.Bot.ReplyMessage(
							event.ReplyToken,
							linebot.NewTextMessage("ยินดีรับใช้ 🏸"),
						).Do()
						if err != nil {
							h.Logger.Sugar().Errorf("Reply message error: %v", err)
						}
					default:
						// case unknow command !
						_, err = h.Bot.ReplyMessage(
							event.ReplyToken,
							linebot.NewTextMessage("ขออภัย ไม่เข้าใจคำสั่ง: "+text),
						).Do()
						if err != nil {
							h.Logger.Sugar().Errorf("Unknown command error: %v", err)
						}
					}
				}
			}

		}
	}

	ctx.OK(nil)
}
