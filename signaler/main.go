package main

import (
	"fmt"
	"math/rand"
	"net/url"

	"github.com/pions/signaler"
)

type MySignalerServer struct {
}

func randSeq(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (m *MySignalerServer) AuthenticateRequest(params url.Values) (apiKey, room, sessionKey string, ok bool) {
	return "ABC", "ABC", randSeq(16), true
}

func (m *MySignalerServer) OnClientMessage(ApiKey, Room, SessionKey string, raw []byte) {
}

func main() {
	fmt.Println(signaler.Start(&MySignalerServer{}))
}