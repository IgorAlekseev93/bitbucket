package event_consumer

import (
	"VKBot/events"
	"VKBot/events/vk"
	"log"
	"net/url"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
}

func New(fetcher events.Fetcher, processor events.Processor) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch()
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}
		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)
			continue
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)
		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())
			continue
		}
	}

	return nil
}

func (c *Consumer) FetcherLongPoll(groupID string) {
	lps, err := c.fetcher.(*vk.Processor).Vk.GetLongPollServer(groupID)

	if err != nil {
		log.Printf("[ERR] FetcherLongPoll consumer: %s", err.Error())
	}

	u, _ := url.Parse(lps.Server)

	c.fetcher.(*vk.Processor).Vk.Host = u.Host
	c.fetcher.(*vk.Processor).Vk.BasePath = u.Path
	c.fetcher.(*vk.Processor).Ts = lps.Ts
	c.fetcher.(*vk.Processor).Key = lps.Key
}
