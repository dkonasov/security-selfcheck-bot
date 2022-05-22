package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var token = os.Getenv("BOT_TOKEN")

var scenarios Scenarios = Scenarios{
	Scenario{"Надежный пароль", []string{
		"Пароль должен быть длиной хотя бы 6 символов",
		"Пароль не должен быть словарным словом, датой рождения или любым другим осмысленным сочетанием букв и цифр",
		"Пароль должен содержать хотя бы один символ из следущих категорий: строчные и прописные буквы, цифры и специальные символы (знаки препинания, символ подчеркивания и пр.)",
		"Пароль должен быть сохранен в надежном хранилище паролей (например связка ключей Apple или браузерное паролехранилище)",
	}},
	Scenario{"Защита аккаунта на npmjs.com", []string{
		"Задайте надежный пароль (см. соответствующий чеклист)",
		"Установите двухфакторную аутентификацию",
		"Если вы деплоите в npm пакеты через публичные инструменты для CI (например Github action), установите двухфакторную аутентификацию и в этих сервисах",
	}},
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})
	handler := func(update Update, bot *Bot) {
		ctx := context.Background()
		raw_user_state, err := rdb.Get(ctx, strconv.FormatUint(update.Message.From.Id, 10)).Result()
		_ = raw_user_state
		if err == nil || err.Error() == "redis: nil" {
			err = nil
			var user_state UserState
			if raw_user_state != "" {
				err = json.Unmarshal([]byte(raw_user_state), &user_state)
			}
			if err == nil {
				keyboard, _ := json.Marshal(scenarios.MakeKeyboard(3))
				params := make(url.Values)
				message_text := "Выберите сценарий для самопроверки"
				if user_state.Scenario == "" {
					scenario := scenarios.GetScenario(update.Message.Text)
					if scenario.Name != "" {
						user_state.Scenario = update.Message.Text
						message_text = scenario.Steps[0]
						keyboard, _ = json.Marshal(ReplyKeyboardMarkup{[][]string{{"Далее"}}})
						user_state.Step++
					}
				} else {
					scenario := scenarios.GetScenario(user_state.Scenario)
					if user_state.Step < len(scenario.Steps) {
						message_text = scenario.Steps[user_state.Step]
						keyboard, _ = json.Marshal(ReplyKeyboardMarkup{[][]string{{"Далее"}}})
						user_state.Step++
					} else {
						message_text = "Поздравляем! Вы завершили самопроверку"
						user_state.Scenario = ""
						user_state.Step = 0
					}
				}
				params.Set("chat_id", strconv.FormatUint(update.Message.From.Id, 10))
				params.Set("text", message_text)
				params.Set("reply_markup", string(keyboard))
				bot.doMethod("sendMessage", params)
				serialized_state, err := json.Marshal(user_state)
				if err == nil {
					err = rdb.Set(ctx, strconv.FormatUint(update.Message.From.Id, 10), string(serialized_state), 0).Err()
					if err != nil {
						fmt.Println("Write to database error")
						fmt.Println(err.Error())
					}
				} else {
					fmt.Println("User state serializtion error")
					fmt.Println(err.Error())
				}
				// err = rdb.Set(strconv.FormatUint(update.Message.From.Id, 10))
			} else {
				fmt.Println("User state error")
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("Redis read error")
			fmt.Println(err.Error())
		}
	}
	bot := Bot{token, os.Getenv("HOST"), os.Getenv("PORT"), handler}
	bot.Start()
}
