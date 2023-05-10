package vk

import (
	"VKBot/client/vk"
	"strings"
)

const (
	HelpCmd  = "help"
	StartCmd = "start"
)

var (
	SelectedСolor = ""
)

func (p *Processor) doCmd(text string, messageID, userID int) error {
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	keyboard := vk.NewKeyboard(true)

	switch text {
	case HelpCmd:
		SelectedСolor = ""
		keyboardString := keyboard.FirstLayer()
		return p.Vk.SendMessage(userID, msgHelp, keyboardString)
	case StartCmd:
		SelectedСolor = ""
		keyboardString := keyboard.FirstLayer()
		return p.Vk.SendMessage(userID, msgHello, keyboardString)
	case "зеленый", "красный", "белый", "черный":
		SelectedСolor = text
		keyboardString := keyboard.SecondLayer(text)
		return p.Vk.SendMessage(userID, "Выберите темный(+) или светлый(-) оттенок цвета.", keyboardString)
	case "+ " + SelectedСolor:
		SelectedСolor = ""
		return p.Vk.SendMessage(userID, "Вы выбрали темный оттенок, цвет "+text, "")
	case "- " + SelectedСolor:
		SelectedСolor = ""
		return p.Vk.SendMessage(userID, "Вы выбрали светлый оттенок, цвет "+text, "")
	default:
		SelectedСolor = ""
		keyboardString := keyboard.FirstLayer()
		return p.Vk.SendMessage(userID, msgUnknown, keyboardString)
	}
}
