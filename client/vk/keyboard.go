package vk

import "encoding/json"

var (
	СolorButton = map[string]string{"зеленый": "positive", "красный": "negative", "белый": "primary", "черный": "secondary"}
)

type Keyboard struct {
	Inline  bool       `json:"inline"`
	Buttons [][]Button `json:"buttons"`
}

type Button struct {
	Action Action `json:"action"`
	Color  string `json:"color"`
}

type Action struct {
	Type  string `json:"type"`
	Label string `json:"label"`
}

func NewKeyboard(Inline bool) *Keyboard {
	return &Keyboard{
		Inline:  Inline,
		Buttons: make([][]Button, 1),
	}
}

func (k *Keyboard) AddLine() {
	k.Buttons = append(k.Buttons, []Button{})
}

func (k *Keyboard) AddButtonInLine(indx int, typeButton, label, color string) {
	k.Buttons[indx] = append(k.Buttons[indx], Button{Action{typeButton, label}, color})
}

func (k *Keyboard) KeyboardString() string {
	keyboardJson, _ := json.Marshal(k)
	return string(keyboardJson)
}

func (k *Keyboard) FirstLayer() string {
	k.AddButtonInLine(0, "text", "Зеленый", "positive")
	k.AddButtonInLine(0, "text", "Красный", "negative")
	k.AddLine()
	k.AddButtonInLine(1, "text", "Белый", "primary")
	k.AddButtonInLine(1, "text", "Черный", "secondary")
	return k.KeyboardString()
}

func (k *Keyboard) SecondLayer(color string) string {
	k.AddButtonInLine(0, "text", "+ "+color, СolorButton[color])
	k.AddButtonInLine(0, "text", "- "+color, СolorButton[color])
	return k.KeyboardString()
}
