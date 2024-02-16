package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
)

func main() {
	host := os.Getenv("HOST")

	pref := tele.Settings{
		Token: os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{
			Timeout: 10 * time.Second,
		},
	}

	b, err := tele.NewBot(pref)

	if err != nil {
		log.Fatal(err)
	}

	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Hello")
	})

	b.Handle("/help", func(c tele.Context) error {
		return c.Send("Help")
	})

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	mux.HandleFunc(fmt.Sprintf("POST /%s", pref.Token), func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			log.Println("Failed to parse update:", err)
			return
		}

		var update tele.Update
		if err := json.Unmarshal(body, &update); err != nil {
			log.Println("Failed to parse update:", err)
			return
		}

		b.ProcessUpdate(update)
	})

	mux.HandleFunc("GET /webhook", func(w http.ResponseWriter, r *http.Request) {
		webhook := tele.Webhook{
			Listen: host,
			Endpoint: &tele.WebhookEndpoint{
				PublicURL: fmt.Sprintf("%s/%s", host, pref.Token),
			},
		}
		b.SetWebhook(&webhook)

		w.Write([]byte("OK, Set"))
	})

	http.ListenAndServe("localhost:8000", mux)
}
