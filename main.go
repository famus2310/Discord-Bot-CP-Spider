package main

import (
  "fmt"
  "os/signal"
  "github.com/bwmarrin/discordgo"
  "errors"
  "syscall"
  "os"
  scraper "Discord-Bot-CP-Spider/scraper"
)

var (
  token string
  contests []scraper.Contest
)

func getToken() (string, error) {
  t := os.Getenv("BOT_TOKEN")
  if t == "" {
    return "", errors.New("Invalid Token")
  }
  return t, nil
}

func main() {
  contests = scraper.Scrape()
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
  defer dg.Close()
  if err != nil {
    fmt.Println(err.Error())
    return
  }
  fmt.Println("Bot is Running")
  sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
  <-sc

}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
  if m.Author.ID == s.State.User.ID {
    return
  }
  if m.Content == "!schedule" {
    message := ""
    for _, body := range contests {
      message = "(" + body.Status + ") " + body.Title + " " +  body.Link + "\n"
      _, err := s.ChannelMessageSend(m.ChannelID, message)
      if err != nil {
        fmt.Println(err.Error())
      }
    }
  }
}
