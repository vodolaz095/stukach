package main

import (
	"context"
	"flag"
	"log"
	"os"
	"sync"

	"github.com/Shopify/go-rspamd/v3"

	"github.com/vodolaz095/stukach/config"
	"github.com/vodolaz095/stukach/service"
)

const bufferSize = 1000

func main() {
	wg := sync.WaitGroup{}
	wg1 := sync.WaitGroup{}
	wg.Add(1)
	logger := log.New(os.Stdout, "stukach ", log.Lshortfile|log.Ltime)
	var dryRun, learnSpam bool
	var err error
	var pathToConfig string
	ctx, cancel := context.WithCancel(context.Background())

	flag.StringVar(&pathToConfig, "config", "./config.yaml", "путь к файлу конфигурации")
	flag.BoolVar(&learnSpam, "learn", false, "загружать сообщения как спам")
	flag.BoolVar(&dryRun, "dry", false, "имитировать отправку данных в rspamd")
	flag.Parse()

	cfg, err := config.LoadFromFile(pathToConfig)
	if err != nil {
		log.Fatalf("%s : при загрузке конфигурации из %s", err, pathToConfig)
	}
	srv := service.RspamdReporterService{
		Logger: logger,
		Client: rspamd.New(cfg.Rspamd.URL),
		WG:     &wg,
		DryRun: dryRun,
		Learn:  learnSpam,
		Config: cfg.Rspamd,
	}
	err = srv.Ping(ctx)
	if err != nil {
		logger.Fatalf("%s : при проверке доступа к %s", err, cfg.Rspamd.URL)
	}
	feed := make(chan []byte, bufferSize)

	go func() {
		err = srv.Start(ctx, feed)
		if err != nil {
			logger.Fatalf("%s : запуске сервиса проверки сообщений", err)
		}
		logger.Printf("Все сообщения отправлены в rspamd...")
		wg.Done()
	}()

	for i := range cfg.Inputs {
		wg1.Add(1)
		go func(j int) {
			fetcher := service.MailboxService{Logger: logger, WG: &wg}
			err = fetcher.Dial(ctx, cfg.Inputs[j])
			if err != nil {
				log.Printf("%s : при соединении с сервером %s как %s",
					err, cfg.Inputs[j].Server, cfg.Inputs[j].Username)
				return
			}
			logger.Printf("Получаем письма с %s из директории %s...",
				cfg.Inputs[j].Server, cfg.Inputs[j].Directory)

			err = fetcher.Fetch(ctx, feed)
			if err != nil {
				log.Printf("%s : при получении данных с %s из %s",
					err, cfg.Inputs[j].Server, cfg.Inputs[j].Directory)
				return
			}

			logger.Printf("Письма загружены с %s из директории %s!",
				cfg.Inputs[j].Server, cfg.Inputs[j].Directory)

			wg1.Done()
		}(i)
	}
	wg1.Wait()
	close(feed)

	wg.Wait()
	logger.Printf("Работа завершена!")
	cancel()
}
