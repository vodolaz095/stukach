package main

import (
	"context"
	"flag"
	"log"
	"os"
	"sync"

	"github.com/Shopify/go-rspamd"
	"github.com/vodolaz095/stukach/config"
	"github.com/vodolaz095/stukach/data"
	"github.com/vodolaz095/stukach/service"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	logger := log.New(os.Stdout, "stukach ", log.Lshortfile|log.Ltime)

	var err error
	var pathToConfig string
	ctx, cancel := context.WithCancel(context.Background())

	flag.StringVar(&pathToConfig, "config", "./stukach.yaml", "путь к файлу конфигурации")
	flag.Parse()

	cfg, err := config.LoadFromFile(pathToConfig)
	if err != nil {
		log.Fatalf("%s : при загрузке конфигурации из %s", err, pathToConfig)
	}
	client := rspamd.New(cfg.Rspamd.URL, rspamd.Credentials(cfg.Rspamd.Username, cfg.Rspamd.Password))
	srv := service.RspamdReporterService{
		Logger: logger,
		Client: client,
		WG:     &wg,
	}
	err = srv.Ping(ctx)
	if err != nil {
		logger.Fatalf("%s : при проверке доступа к %s", err, cfg.Rspamd.URL)
	}
	feed := make(chan []byte, 10)
	feed <- data.TestEmail

	for i := range cfg.Inputs {
		wg.Add(1)
		fetcher := service.MailboxService{Logger: logger, WG: &wg}
		err = fetcher.Dial(ctx, cfg.Inputs[i])
		if err != nil {
			log.Printf("%s : при соединении с сервером %s как %s",
				err, cfg.Inputs[i].Server, cfg.Inputs[i].Username)
			continue
		}
		logger.Printf("Получаем письма с %s из директории %s",
			cfg.Inputs[i].Server, cfg.Inputs[i].Directory)

		err = fetcher.Fetch(ctx, feed)
		if err != nil {
			log.Printf("%s : при получении данных с %s из %s",
				err, cfg.Inputs[i].Server, cfg.Inputs[i].Directory)
			continue
		}

		logger.Printf("Письма загружены с %s из директории %s",
			cfg.Inputs[i].Server, cfg.Inputs[i].Directory)
		wg.Done()
	}

	go func() {
		wg.Wait()
		logger.Printf("Все сообщения обработаны!")
		cancel()
	}()

	err = srv.Start(ctx, feed)
	if err != nil {
		logger.Fatalf("%s : запуске сервиса проверки сообщений", err)
	}
}
