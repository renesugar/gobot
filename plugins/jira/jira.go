package jira

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/gabeguz/gobot"
	gb "github.com/gabeguz/gobot/bots/gobot"
	sb "github.com/gabeguz/gobot/bots/slack"
	"github.com/nlopes/slack"
	"otremblay.com/jkl"
)

type Jira struct{}

func init() {
	jkl.FindRCFile()
}

func (p Jira) Name() string {
	return "Jira v1.0"
}

var ticketRE = regexp.MustCompile("[A-Z]+-[0-9]+")

func (p Jira) Execute(msg gobot.Message, bot gobot.Bot) error {
	b2 := bot.(gb.Gobot)
	if msg.From() != bot.FullName() {
		matches := ticketRE.FindAllString(msg.Body(), -1)
		if len(matches) > 0 {
			issues, err := jkl.List(fmt.Sprintf("key in (%s)", strings.Join(matches, ",")))
			if err != nil {
				bot.Send(fmt.Sprintf("I AM ERROR: %s", err.Error()))
			}

			switch b3 := b2.InternalBot().(type) {
			case *sb.Bot:
				c := b3.Client()
				p := slack.NewPostMessageParameters()
				p.EscapeText = false
				p.Username = b3.Opt.Name
				p.Attachments = make([]slack.Attachment, 0, len(issues))
				for _, issue := range issues {
					a := slack.Attachment{
						Title:     fmt.Sprintf("%s : %s", issue.Key, issue.Fields.Summary),
						TitleLink: issue.URL(),
						Text:      issue.Fields.Description,
					}
					p.Attachments = append(p.Attachments, a)
				}

				c.PostMessage(b3.Opt.Room, "", p)
			default:
				b := bytes.NewBuffer(nil)
				for _, issue := range issues {
					fmt.Fprintln(b, fmt.Sprintf("%s|%s : %s", issue.URL(), issue.Key, issue.Fields.Summary))
				}
				bot.Send(b.String())
			}
		}
	}

	return nil
}