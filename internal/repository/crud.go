package repository

import (
	"context"
	"log"
	"sync"
	"tgent/ent"
	"tgent/ent/user"
	"tgent/internal/config"
	"github.com/go-telegram/bot"
)


func MakeUser(ctx context.Context, client *ent.Client, counter chan struct{}, u config.User, wg *sync.WaitGroup) error {
	defer wg.Done()
	defer func() { <- counter }()
	user, err := client.User.Create().SetName(u.Name).SetTgID(u.TgID).Save(ctx)
	if err != nil {
		log.Printf("Ошибка во время создания пользователя: %s", err)
		return nil
	}
	log.Printf("Пользователь был успешно создан: %v", user)
	return nil
}


func SelectUser(ctx context.Context, client *ent.Client, ID string, counter chan struct{}, selected chan <- *ent.User, b *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() { <- counter }()
	user, err := client.User.Query().Where(user.TgIDEQ(ID)).WithPosts().First(ctx)
	if err != nil {
		log.Fatalf("Ошибка при выборе пользователя: %s", err)
	}
	log.Printf("Пользователь найден: %v \n", user)
	selected <- user
}


