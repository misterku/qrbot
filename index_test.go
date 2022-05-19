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

func UrlShortenerTests(t *testing.T) {
	url, _ := clckApiCheck("http://github.com")
	if url != "https://clck.ru/AJUUf" {
		t.Fatal()
	}
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
	isValidUrl("http://acm.sgu.ru")
	isValidUrl("https://ya.ru")
	isValidUrl("http://127.0.0.1")
	isInvalidUrl("foo://bar")
	isInvalidUrl("http://")
	isValidUrl("vk.com")
	vkUrl, err := extractURL("vk.com")
	if err != nil || vkUrl != "https://vk.com" {
		t.Fatal()
	}
	rusUrl, err := extractURL("www.россия.рф")
	if err != nil || rusUrl != "http://www.россия.рф" {
		t.Fatal()
	}
	isValidUrl("http://www.россия.рф")
}
