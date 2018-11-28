package main

import (
	"github.com/mattn/go-xmpp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"fmt"
)

func main() {
	var xmppClient *xmpp.Client

	viper.SetConfigName("xmpp-echo-bot")
	//viper.AddConfigPath("/etc/xmpp-echo-bot/")
	//viper.AddConfigPath("$HOME/.xmpp-echo-bot")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("xeb")
	viper.AutomaticEnv()
	
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		log.Error(fmt.Sprintf("Fatal error config file: %s \n", err))
	}

	host := viper.GetString("Host")
	user := viper.GetString("User")
	pass :=  viper.GetString("Password")

	if host == "" || user == "" || pass == "" {
		panic(fmt.Errorf("Missing config variables"))
	}

	options := xmpp.Options{
		Host: host,
		User: user,
		Password: pass,
		Resource: "xmpp-echo-bot",
		NoTLS: true,
		StartTLS: true,
	}

	xmppClient, _ = options.NewClient()

	messages := make(chan xmpp.Chat)

	go func() {
		log.Info("Starting listening for messages ...")

		for {
			chat, err := xmppClient.Recv()
			if err != nil {
				log.Fatal(err)
			}
			switch v := chat.(type) {
				case xmpp.Chat:
					messages <- v
			}
		}
	}()

	for message := range messages {
		if message.Text != "" {
			log.WithFields(log.Fields{
				"jid": message.Remote,
			}).Info("Got message, replying")
	
			xmppClient.Send(xmpp.Chat{
				Remote: message.Remote,
				Type: "chat",
				Text: message.Text,
			})
		}
	}
}