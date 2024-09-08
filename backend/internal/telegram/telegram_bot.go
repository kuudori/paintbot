package telegram

import (
	"PaintBackend/internal/config"
	database "PaintBackend/internal/storage/models"
	"context"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"log/slog"
	"strconv"
)

type Handlers struct {
	FileStorage database.FileRepository
}

func NewHandler(fileStorage database.FileRepository) *Handlers {
	return &Handlers{FileStorage: fileStorage}
}

func (h *Handlers) Start(b *gotgbot.Bot, ctx *ext.Context) error {
	cfg := config.GetConfig()
	b.SetChatMenuButton(

		&gotgbot.SetChatMenuButtonOpts{
			ChatId: &ctx.Message.Chat.Id,
			MenuButton: gotgbot.MenuButtonWebApp{Text: "New Paint",
				WebApp: gotgbot.WebAppInfo{Url: cfg.BackendDomain}},
		},
	)
	_, err := ctx.EffectiveMessage.Reply(
		b,
		"Hello! This paint bot lets you draw pics and send them to any chat by mentioning the bot.",
		&gotgbot.SendMessageOpts{})
	if err != nil {
		slog.Error("Send /start message error")
		return err
	}
	return nil
}

func (h *Handlers) Source(b *gotgbot.Bot, ctx *ext.Context) error {
	bgCtx := context.Background()
	images, err := h.FileStorage.GetFilesURLByChatID(bgCtx, ctx.InlineQuery.From.Id, 1)
	if err != nil {
		slog.Error("failed to get file names")
		return err
	}
	offset, _ := strconv.Atoi(ctx.InlineQuery.Offset)
	if offset < 0 {
		offset = 0
	}

	var results []gotgbot.InlineQueryResult
	end := offset + 10
	if end > len(images) {
		end = len(images)
	}

	for i, imageURL := range images[offset:end] {
		results = append(results, gotgbot.InlineQueryResultPhoto{
			Id:           strconv.Itoa(offset + i + 1),
			ThumbnailUrl: imageURL,
			PhotoUrl:     imageURL,
			PhotoWidth:   810,
			PhotoHeight:  1080,
		})
	}

	nextOffset := ""
	if end < len(images) {
		nextOffset = strconv.Itoa(end)
	}

	_, err = ctx.InlineQuery.Answer(b, results, &gotgbot.AnswerInlineQueryOpts{
		IsPersonal: true,
		CacheTime:  5,
		NextOffset: nextOffset,
	})
	if err != nil {
		slog.Error("failed to send inline query answer: ", err)
		return err
	}

	return nil
}
