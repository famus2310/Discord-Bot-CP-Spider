package main

import (
  "fmt"
  "github.com/bwmarrin/discordgo"
  "errors"
  "os"
)

var (
  token string
)

func getToken() (string, error) {
  t := os.Getenv("BOT_TOKEN")
  if t == "" {
    return "", errors.New("Invalid Token")
  }
  return t, nil
}

func main() {
  token, err := getToken()
  if err != nil {
    fmt.Println(err.Error())
    return
  }

  dg, err := discordgo.New("Bot " + token)
  if err != nil {
    fmt.Println(err.Error())
    return
  }

  dg.AddHandler(messageHandler)
  err = dg.Open()

  if err != nil {
    fmt.Println(err.Error())
    return
  }
  fmt.Println("Bot is Running")
  sc := make(chan os.Signal, 1)
  <-sc

  dg.Close()
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
  if m.Author.ID == s.State.User.ID {
    return
  }

  s.ChannelMessageSend(m.ChannelID, m.Content)
}
