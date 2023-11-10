package main

import (
	"context"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/log"
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
			gateway.WithPresenceOpts(gateway.WithWatchingActivity("super reactions"))),
		bot.WithCacheConfigOpts(cache.WithCaches(cache.FlagsNone)),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnGuildMessageReactionAdd: func(event *events.GuildMessageReactionAdd) {
				if !event.Burst {
					return
				}
				rest := event.Client().Rest()
				if err := rest.RemoveAllReactionsForEmoji(event.ChannelID, event.MessageID, event.Emoji.Reaction()); err != nil {
					log.Error("there was an error while removing a burst reaction: ", err)
				}
			},
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
