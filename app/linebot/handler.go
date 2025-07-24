package linebot

import (
	"fmt"
	"strings"

	"github.com/inkxk/bad-bot/app"
	"github.com/inkxk/bad-bot/constant"
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

				// check prefix
				switch {
				case strings.HasPrefix(text, "!ตีแบดกัน"):
					h.handleBadmintonCommand(event, text)
				default:
					// do nothing
				}
			}

		}
	}

	ctx.OK(nil)
}

func (h *Handler) handleBadmintonCommand(event *linebot.Event, text string) {
	args := strings.TrimSpace(strings.TrimPrefix(text, "!ตีแบดกัน"))

	if args == "" {
		_, err := h.Bot.ReplyMessage(
			event.ReplyToken,
			linebot.NewTextMessage(fmt.Sprintf(constant.DEFAULT_MESSAGE, "วันเสาร์", "N")),
		).Do()

		if err != nil {
			h.Logger.Sugar().Errorf("Reply default message error: %v", err)
		}
	} else {
		// get parts
		partsRaw := strings.Split(args, ",")
		var parts []string
		for _, p := range partsRaw {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				parts = append(parts, trimmed)
			}
		}

		if len(parts) < 2 {
			_, err := h.Bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage(`มึงบอกวันที่ กับ จำนวนคอร์ดด้วยสิ เช่น '!ตีแบดกัน, วันเสาร์ ที่ 26 ก.ค., 2'`),
			).Do()
			if err != nil {
				h.Logger.Sugar().Errorf("Invalid input error: %v", err)
			}
			return
		} else if len(parts) > 2 {
			_, err := h.Bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage(`พิมเหี้ยไรเยอะแยะ`),
			).Do()
			if err != nil {
				h.Logger.Sugar().Errorf("Invalid input error: %v", err)
			}
			return
		}

		date := strings.TrimSpace(parts[0])
		countNumber := strings.TrimSpace(parts[1])

		_, err := h.Bot.ReplyMessage(
			event.ReplyToken,
			linebot.NewTextMessage(fmt.Sprintf(constant.DEFAULT_MESSAGE, date, countNumber)),
		).Do()

		if err != nil {
			h.Logger.Sugar().Errorf("Reply fallback error: %v", err)
		}
	}
}
