package service

import (
	"bytes"
	"context"
	"log"
	"sync"

	"github.com/Shopify/go-rspamd"
)

type RspamdReporterService struct {
	Logger *log.Logger
	Client rspamd.Client
	WG     *sync.WaitGroup
}

func (r *RspamdReporterService) Ping(ctx context.Context) (err error) {
	_, err = r.Client.Ping(ctx)
	return
}

func (r *RspamdReporterService) Start(ctx context.Context, feed chan []byte) (err error) {
	for {
		select {
		case msg := <-feed:
			eml := rspamd.NewEmailFromReader(bytes.NewReader(msg))
			res, checkErr := r.Client.Check(ctx, eml)
			if checkErr != nil {
				r.Logger.Printf("%s : при проверке сообщения", checkErr)
				r.WG.Done()
				continue
			}
			r.Logger.Println("message_id", res.MessageID)
			r.Logger.Println("score", res.Score)
			//for s := range res.Symbols {
			//	r.Logger.Printf("Name: %s. Score: %v. Metric score: %v. Description: %s",
			//		res.Symbols[s].Name,
			//		res.Symbols[s].Score,
			//		res.Symbols[s].MetricScore,
			//		res.Symbols[s].Description,
			//	)
			//}
			r.WG.Done()
			break
		case <-ctx.Done():
			return nil
		}
	}
}
