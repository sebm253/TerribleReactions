package main

import (
	"context"
	"encoding/json"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetLevel(log.LevelInfo)
	log.Info("starting the bot...")
	log.Info("disgo version: ", disgo.Version)

	client, err := disgo.New(os.Getenv("TERRIBLE_REACTIONS_TOKEN"),
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuildMessageReactions),
			gateway.WithPresenceOpts(gateway.WithWatchingActivity("super reactions")),
			gateway.WithEnableRawEvents(true)),
		bot.WithCacheConfigOpts(cache.WithCaches(cache.FlagsNone)),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnRaw: onRaw,
		}))
	if err != nil {
		log.Fatal("error while building disgo instance: ", err)
	}

	defer client.Close(context.TODO())

	if err := client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to the gateway: ", err)
	}

	log.Info("terrible reactions bot is now running.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func onRaw(event *events.Raw) {
	if event.EventType != gateway.EventTypeMessageReactionAdd {
		return
	}
	var payload ReactionPayload
	if err := json.NewDecoder(event.Payload).Decode(&payload); err != nil {
		log.Error("there was an error while decoding the payload: ", err)
		return
	}
	if payload.Burst {
		rest := event.Client().Rest()
		if err := rest.RemoveUserReaction(payload.ChannelID, payload.MessageID, payload.Emoji.Reaction(), payload.UserID); err != nil {
			log.Error("there was an error while removing a burst reaction: ", err)
		}
	}
}

type ReactionPayload struct {
	UserID    snowflake.ID         `json:"user_id"`
	MessageID snowflake.ID         `json:"message_id"`
	ChannelID snowflake.ID         `json:"channel_id"`
	Emoji     discord.PartialEmoji `json:"emoji"`
	Burst     bool                 `json:"burst"`
}
