package errorController

import (
	"fmt"
	"os"

	"github.com/DisgoOrg/disgohook"
	"github.com/DisgoOrg/disgohook/api"
)

const (
	webhookURLWeb = "997758066130952212/v9V9OgtxkLM-Pw6WDfZ_JfIaL_pI7sOSjJECccBa7CuqoiIBwfdR9b2HAhIxdy6npcWe"
)

var (
	WD    Message
	IsPro bool
)

type Message struct {
	Content string
	Site    string
}

func (m *Message) SendErrorWHWeb() {
	webhook, err := disgohook.NewWebhookClientByToken(nil, nil, webhookURLWeb)
	if err != nil {
		panic(err)
	}
	if IsPro {
		_, err = webhook.SendEmbeds(api.NewEmbedBuilder().
			SetTitle("New error on pro!").
			SetDescription(m.Site + ": " + m.Content + " @here").
			SetColor(14177041).
			Build(),
		)
	} else {
		_, err = webhook.SendEmbeds(api.NewEmbedBuilder().
			SetTitle("New error on the site!").
			SetDescription(m.Site + ": " + m.Content + " @here").
			SetColor(14177041).
			Build(),
		)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
func HandleError(message, site string) {
	t := Message{
		Content: message,
		Site:    site,
	}
	t.SendErrorWHWeb()
	ErrorLogger.Println(message)
}
