package main

import (
	VKClient "VKBot/client/vk"
	event_consumer "VKBot/consumer/event-consumer"
	"VKBot/events/vk"
	"flag"
	"log"
)

const (
	version = "5.131"
	host    = "api.vk.com"
	//groupID   = "105505449"
)

func main() {
	groupID, token := mustGroupIdToken()

	eventsFetcher := vk.New(VKClient.New(host, token, version, "method"))
	eventsProcessor := vk.New(VKClient.New(host, token, version, "method"))

	consumer := event_consumer.New(eventsFetcher, eventsProcessor)
	consumer.FetcherLongPoll(groupID)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustGroupIdToken() (string, string) {
	groupID := flag.String(
		"bot-groupID",
		"",
		"groupID for access")

	token := flag.String(
		"bot-token",
		"",
		"token for access to bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	if *groupID == "" {
		log.Fatal("host is not")
	}

	return *groupID, *token
}
