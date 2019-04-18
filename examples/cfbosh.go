package main

import (
	"log"
	"os"

	"github.com/dpb587/go-slack-topic-bot/message"
	"github.com/dpb587/go-slack-topic-bot/message/boshio"
	"github.com/dpb587/go-slack-topic-bot/message/github"
	"github.com/dpb587/go-slack-topic-bot/message/pairist"
	"github.com/dpb587/go-slack-topic-bot/message/trello"
	"github.com/dpb587/go-slack-topic-bot/slack"
)

func main() {
	msg, err := message.Join(
		" || ",
		message.Prefix(
			"_interrupt_ ",
			message.Join(
				" ",
				message.Conditional(
					pairist.WorkingHours("12:00", "18:00", "America/Los_Angeles"),
					pairist.PeopleInRole{
						Team: "sf-bosh",
						Role: "interrupt",
						People: map[string]string{
							"Luan":    "U02R9SUJX",
							"Josh R":  "U8DLATS12",
							"Josh":    "U0YGVGTM1",
							"Danny":   "U0FUK0EBH",
							"Mike":    "U21JVA9F0",
							"Jim":     "U02QZ1E3G",
							"Morgan":  "U04V9L81Y",
							"Belinda": "U5EJ8MQUW",
							"Max":     "U4FFS1UAK",
						},
					},
				),
				message.Conditional(
					pairist.WorkingHours("06:30", "12:00", "America/Los_Angeles"),
					pairist.PeopleInRole{
						Team: "boshto",
						Role: "Interrupt",
						People: map[string]string{
							"Gaurab":  "U0A0ZUT43",
							"Dale":    "U32RHRLE9",
							"Rebecca": "U8YCN97Q9",
							"Andrew":  "U17K4GAKW",
							"Fred":    "UA3MK3AE7",
							"Jamil":   "U0717EQ04",
						},
					},
				),
				message.Conditional(
					pairist.WorkingHours("09:00", "18:00", "Europe/Berlin"),
					trello.PeopleInRole{
						Team: "bosh-europe",
						Role: "Interrupt",
						People: map[string]string{
							"s4heid":  "U8MRYTRHU",
							"beyhan6": "U0D8E67LZ",
						},
					},
				),
			),
		),
		message.Literal("_docs_ <https://bosh.io|bosh.io>"),
		message.Prefix(
			"_latest_ ",
			message.Join(
				" ",
				boshio.Release{Alias: "bosh", Repository: "github.com/cloudfoundry/bosh"},
				github.Release{Token: os.Getenv("GITHUB_TOKEN"), Alias: "bosh-cli", Owner: "cloudfoundry", Name: "bosh-cli"},
				boshio.Stemcell{Alias: "ubuntu-xenial", Name: "bosh-aws-xen-hvm-ubuntu-xenial-go_agent"},
			),
		),
	).Message()
	if err != nil {
		log.Panicf("ERROR: %v", err)
	}

	log.Printf("DEBUG: expected message: %s", msg)

	err = slack.UpdateChannelTopic(os.Getenv("SLACK_CHANNEL"), msg)
	if err != nil {
		log.Panicf("ERROR: %v", err)
	}
}
