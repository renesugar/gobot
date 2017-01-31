package main

import (
	"flag"
	"log"

	"github.com/gabeguz/gobot"
	"github.com/gabeguz/gobot/bot"
	"github.com/gabeguz/gobot/bots/slack"
	"github.com/gabeguz/gobot/plugins/beer"
	"github.com/gabeguz/gobot/plugins/chatlog"
	"github.com/gabeguz/gobot/plugins/dm"
	"github.com/gabeguz/gobot/plugins/echo"
	"github.com/gabeguz/gobot/plugins/quote"
	"github.com/gabeguz/gobot/plugins/rickroll"
	"github.com/gabeguz/gobot/plugins/stathat"
	"github.com/gabeguz/gobot/plugins/troll"
	"github.com/gabeguz/gobot/plugins/url"
)

func main() {
	var host, user, pass, room, name string
	flag.StringVar(&host, "host", "", "Hostname:port of the XMPP server")
	flag.StringVar(&user, "user", "", "Username of XMPP server (i.e.: foo@hostname.com")
	flag.StringVar(&pass, "pass", "", "Password for XMPP server")
	flag.StringVar(&room, "room", "", "Room to join (i.e.: #myroom@hostname.com")
	flag.StringVar(&name, "name", "gobot", "Name of the bot")
	flag.Parse()

	//TODO:Add some validation...but whatever for now

	bot := Gobot{
		slack.New(pass, room, name),
		[]gobot.Plugin{
			echo.Echo{},
			beer.Beer{},
			quote.Quote{},
			dm.DirectMessage{},
			stathat.StatHat{},
			chatlog.ChatLog{},
			troll.Troll{},
			rickroll.RickRoll{},
			url.Url{},
		},
	}

	/*
		bot := Gobot{
			xmpp.New(host, user, pass, room, name),
			[]gobot.Plugin{
				echo.Echo{},
				beer.Beer{},
				quote.Quote{},
				dm.DirectMessage{},
				stathat.StatHat{},
				chatlog.ChatLog{},
				troll.Troll{},
				rickroll.RickRoll{},
				url.Url{},
			},
		}
	*/

	err := bot.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	//go bot.PingServer(30)
	var msg bot.Message
	var plugin bot.Plugin
	for msg = range gobot.Listen() {
		for _, plugin = range gobot.plugins {
			go executePlugin(plugin, msg, gobot)
		}
	}

}

func executePlugin(p bot.Plugin, m bot.Message, b bot.Bot) {
	err := p.Execute(m, b)
	if err != nil {
		b.Log(p.Name() + " => " + err.Error())
	}
}

type chatbot struct {
	bot.Bot
	plugins []bot.Plugin
}