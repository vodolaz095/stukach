package service

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/Shopify/go-rspamd/v3"
	"github.com/vodolaz095/stukach/config"
)

type RspamdReporterService struct {
	Logger *log.Logger
	Client rspamd.Client
	WG     *sync.WaitGroup
	DryRun bool
	Learn  bool
	Config config.RspamdConnectionConfig
}

func (r *RspamdReporterService) Ping(ctx context.Context) (err error) {
	_, err = r.Client.Ping(ctx)
	return
}

func (r *RspamdReporterService) Start(ctx context.Context, feed chan []byte) (err error) {
B:
	for {
		select {
		case msg, ok := <-feed:
			if !ok {
				r.Logger.Printf("Канал кончился...")
				break B
			}

			if r.DryRun {
				r.Logger.Printf("Имитация отправки %v байт сообщения в rspamd", len(msg))
				r.WG.Done()
			} else {
				if r.Learn {
					// https://rspamd.com/doc/architecture/protocol.html#controller-http-endpoints
					_, learnErr := r.Client.LearnSpam(ctx, &rspamd.LearnRequest{
						Message: bytes.NewReader(msg),
						Header: http.Header{
							"Password": []string{
								r.Config.Password,
							},
						},
					})
					if learnErr != nil {
						if strings.HasPrefix(learnErr.Error(), "Unexpected response code: 208") {
							r.Logger.Printf("Сообщение уже обработано!")
						} else {
							r.Logger.Printf("%s : при обучении на спам сообщении", learnErr)
						}
					}
					r.WG.Done()
				} else {
					res, checkErr := r.Client.Check(ctx, &rspamd.CheckRequest{
						Message: bytes.NewReader(msg),
						Header: http.Header{
							"Password": []string{
								r.Config.Password,
							},
						},
					})
					if checkErr != nil {
						r.Logger.Printf("%s : при проверке сообщения", checkErr)
						r.WG.Done()
					} else {
						r.Logger.Printf("Сообщение %s проверено, его счёт %v!", res.MessageID, res.Score)
						r.WG.Done()
					}
				}
			}
		case <-ctx.Done():
			return nil
		}
	}
	return nil
}
