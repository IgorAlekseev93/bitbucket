package vk

type LongPollServer struct {
	Response Response `json:"response"`
}

type Response struct {
	Key    string `json:"key"`
	Server string `json:"server"`
	Ts     string `json:"ts"`
}

type GetResponse struct {
	Ts      string   `json:"ts"`
	Updates []Update `json:"updates"`
}

type Update struct {
	EventType string `json:"type"`
	Object    Object `json:"object"`
	GroupID   int    `json:"group_id"`
}

type Object struct {
	Message Message `json:"message"`
}

type Message struct {
	Id     int    `json:"id"`
	UserID int    `json:"from_id"`
	Text   string `json:"text"`
}
