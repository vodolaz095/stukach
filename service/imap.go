package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/vodolaz095/stukach/config"
)

type MailboxService struct {
	Logger *log.Logger
	WG     *sync.WaitGroup

	messagesFound uint32
	client        *imapclient.Client
	cfg           config.ImapConfig
}

func (m *MailboxService) Dial(ctx context.Context, cfg config.ImapConfig) (err error) {
	client, err := imapclient.DialTLS(fmt.Sprintf("%s:%v", cfg.Server, cfg.Port), nil)
	if err != nil {
		return
	}
	m.Logger.Printf("Соединение установлено с %s:%v", cfg.Server, cfg.Port)
	m.client = client
	err = m.client.Login(cfg.Username, cfg.Password).Wait()
	if err != nil {
		return
	}
	m.Logger.Printf("Авторизация %s прошла успешно ", cfg.Username)

	mailboxes, err := m.client.List("", "%", nil).Collect()
	if err != nil {
		return
	}
	m.Logger.Printf("Найдено %v директорий...", len(mailboxes))
	for i, mbox := range mailboxes {
		m.Logger.Printf("%v) %s", i, mbox.Mailbox)
	}
	mailbox, err := m.client.Select(cfg.Directory, nil).Wait()
	if err != nil {
		return
	}
	m.Logger.Printf("Директория %s найдена c %v сообщениями!", cfg.Directory, mailbox.NumMessages)
	m.messagesFound = mailbox.NumMessages
	m.cfg = cfg
	return nil
}

func (m *MailboxService) Fetch(ctx context.Context, feed chan []byte) (err error) {
	if m.messagesFound == 0 {
		return
	}
	seqSet := imap.SeqSetRange(1, m.messagesFound)
	fetchItems := []imap.FetchItem{
		imap.FetchItemUID,
		&imap.FetchItemBodySection{},
	}
	fetchCmd := m.client.Fetch(seqSet, fetchItems, nil)
	for {
		msg := fetchCmd.Next()
		if msg == nil {
			break
		}

		for {
			received := msg.Next()
			if received == nil {
				break
			}

			switch item := received.(type) {
			case imapclient.FetchItemDataUID:
				m.Logger.Printf("Читаем UID %v из %s сервера %s...",
					item.UID, m.cfg.Directory, m.cfg.Server,
				)
			case imapclient.FetchItemDataBodySection:
				b, readErr := io.ReadAll(item.Literal)
				if readErr != nil {
					m.Logger.Printf("failed to read body section: %v", err)
					continue
				}
				m.Logger.Printf("Прочитано %v байт из %s сервера %s...",
					len(b), m.cfg.Directory, m.cfg.Server,
				)
				m.WG.Add(1)
				feed <- b
			}
		}
	}
	err = fetchCmd.Close()
	return
}

func (m *MailboxService) Close() (err error) {
	return m.client.Close()
}
