package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"tgent/ent"
	"tgent/internal/config"
	"tgent/internal/repository"
	"time"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	_ "github.com/lib/pq"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	env_data := config.LoadConfig()
	counter := make(chan struct{}, 100)
	defer cancel()
	client, err := ent.Open("postgres", env_data.DB_URL)
	if err != nil {
		log.Fatalf("Произошла ошибка в основной функции при открытии клиента. Ошибка: %s", err)
	}
	defer client.Close()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	handleStart := func (ctx context.Context, b *bot.Bot, update *models.Update) {
		log.Printf("Сработала функция начала бота у пользователя: %v", update.Message.Chat.ID)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: "Спасибо за доверие сервису",
		})
		usr := config.User{TgID: fmt.Sprintf("%d", update.Message.Chat.ID), Name: update.Message.Chat.Username}
		wg.Add(1)
		counter <- struct{}{}
		go repository.MakeUser(ctx, client, counter, usr, &wg)
	}
	handleGet := func (ctx context.Context, b *bot.Bot, update *models.Update) {
		log.Printf("Сработала функция получения бота у пользователя: %v", update.Message.Chat.ID)
		usr := config.User{TgID: fmt.Sprintf("%d", update.Message.Chat.ID), Name: update.Message.Chat.Username}
		wg.Add(1)
		counter <- struct{}{}
		user := repository.SelectUser(ctx, client, usr.TgID, counter, b, &wg)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: fmt.Sprintf("Пользователь: %v", user),
		})
	}
	opts := []bot.Option{
		bot.WithCheckInitTimeout(1 * time.Minute),
		bot.WithDefaultHandler(handleStart),
	}
	b, err := bot.New(env_data.BOT_KEY, opts...)
	if err != nil {
		log.Fatalf("Ошибка при запуске бота: %s\n", err)
	}
	b.RegisterHandler(bot.HandlerTypeMessageText, "/get", bot.MatchTypeExact, handleGet)
	b.Start(ctx)
	wg.Wait()
	close(counter)
}
