package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"go-jobpass-bot/srcs/entities"
	"go-jobpass-bot/srcs/store"
	"go-jobpass-bot/srcs/usecase"
)

func AutoSave() {
	for {
		time.Sleep(time.Minute * 5)
		log.Infof("Auto-saved data")
		store.SaveStoreInfo()
	}
}

func main() {
	err := store.InitStore()
	if err != nil {
		log.Errorf("couldn't recover store info")
		log.Exit(1)
	}

	defer func() {
		err = store.SaveStoreInfo()
		if err != nil {
			log.Errorf("Error while saving store info")
			return
		}
		log.Infof("JobpassBot turned off")
	}()

	// Create Discord Bot Session
	sess, err := discordgo.New("Bot " + entities.Data.Discord.APIKey)
	if err != nil {
		log.Errorf("couldn't create discord bot session")
		log.Exit(1)
	}

	err = usecase.InitCommands(sess)
	if err != nil {
		return
	}

	sess.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := usecase.CommandsHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	//sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
	//	if m.Author.ID == s.State.User.ID {
	//		return
	//	}
	//
	//	if strings.HasPrefix(m.Content, entities.Data.Discord.Prefix) {
	//		newCommand, err := usecase.NewCommand(strings.TrimPrefix(m.Content, entities.Data.Discord.Prefix), s, m)
	//		if err != nil {
	//			return
	//		}
	//
	//		err = newCommand.ParseCommand()
	//		if err != nil {
	//			return
	//		}
	//
	//	}
	//	if err != nil {
	//		panic(err)
	//	}
	//})
	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	if err != nil {
		log.Errorf("Error while adding CMD | %s", err)
		return
	}

	err = sess.Open()
	if err != nil {
		return
	}
	defer func(sess *discordgo.Session) {
		err := sess.Close()
		if err != nil {
			log.Errorf("Something fucked up while closing the bot | %s", err)
		}
	}(sess)

	log.Infof("JobpassBot is online!")

	go AutoSave()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
