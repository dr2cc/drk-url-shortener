package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	// ENV=local используется для разработки на компьютере.
	// На реальном сервере значение "local" перезапишется реальным значением переменной ENV (например, "prod").
	// Точно также его использует логгер- от максимално разверноутого лога при разработке, до строго json на сервере
	Env           string `env:"ENV" envDefault:"local"`
	ServAddres    string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	CacheDumpPath string `env:"FILE_STORAGE_PATH"`
}

func New() (Config, error) {
	cfg := Config{}
	// 1. Присваиваем полям Config флаги
	// Разбираем флаги в конфигурацию
	flag.StringVar(&cfg.ServAddres, "a", ":8080", "HTTP server startup address")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "base URL")
	flag.StringVar(&cfg.CacheDumpPath, "f", "./short-url-db.json", "path to dump with addresses")
	flag.Parse()

	// // 2. Загужаем в окружение переменные из файла .env в корне.
	// // Не соответствует заданию yp, но для других проектов достаточно снять комментарий.
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatalf("Error loading env variables: %s", err.Error())
	// }

	// parse читает переменные окружения и заполняет структуру cfg . Подробнее:
	// - Библиотека берёт структуру cfg и с помощью механизма рефлексии (способность программы изучать саму себя) смотрит на её поля.
	// - Она ищет поля, у которых есть тег env:"ИМЯ_ПЕРЕМЕННОЙ".
	// - Если в самой операционной системе (в окружении) есть переменная с таким же именем,
	// библиотека caarlos0/env берёт её значение и записывает в это конкретное поле.
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Config parsing error: %+v\n", err)
		return Config{}, err
	}

	return cfg, nil
}
