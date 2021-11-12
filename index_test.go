package main

import "testing"

func TestCommands(t *testing.T) {
	checkInvalidCommand := func(command string) {
		invalidCommand := getMessageForCommand(command)
		if invalidCommand != invalidCommandResponse {
			t.Fatal()
		}
	}

	checkHelpCommand := func(command string) {
		commandResponse := getMessageForCommand(command)
		if commandResponse != helpCommandResponse {
			t.Fatal()
		}
	}

	checkHelpCommand("/start")
	checkHelpCommand("/help")
	checkInvalidCommand("/test")
	checkInvalidCommand("/abcba")
	checkInvalidCommand("/")
}

func TestURLValidator(t *testing.T) {
	isValidUrl := func(url string) {
		_, err := extractURL(url)
		if err != nil {
			t.Fatal()
		}
	}

	isInvalidUrl := func(url string) {
		_, err := extractURL(url)
		if err == nil {
			t.Fatal()
		}
	}

	isValidUrl("https://google.com")
	isValidUrl("https://ya.ru")
	isValidUrl("http://127.0.0.1")
	isInvalidUrl("foo://bar")
	isInvalidUrl("http://")
	isInvalidUrl("vk.com")
	isValidUrl("http://www.россия.рф")
}
