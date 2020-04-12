package main

import (
  "github.com/robfig/cron"
  "fmt"
  "github.com/bwmarrin/discordgo"
  "errors"
  "log"
  "net/http"
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

func getPort() string {
  p := os.Getenv("PORT")
  if p != "" {
    return ":" + p
  }
  return ":3000"
}

func main() {
  c := cron.New()
  contests = scraper.Scrape()
  c.AddFunc("@hourly", func() {
    contests = scraper.Scrape()
  })
  c.Start()
  token, err := getToken()
  port := getPort()
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
  log.Fatal(http.ListenAndServe(port, nil))
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
  if m.Author.ID == s.State.User.ID {
    return
  }
  if m.Content == "!schedule" {
    s.ChannelMessageSend(m.ChannelID, "**Here's Your Contests List:**")
    for _, body := range contests {
      var baseColor int
      switch body.Status {
        case "PAST":
          baseColor = 16711680
        case "RUNNING":
          baseColor = 16776960
        case "COMING":
          baseColor = 65280
      }
      embed := new(discordgo.MessageEmbed)
      embed.URL = body.Link
      embed.Color = baseColor
      embed.Title = body.Title + "\n"
      embed.Description = "**(" + body.Status + ")**"
      inlineFields := []*discordgo.MessageEmbedField{
        {Name: "Duration", Value: body.Duration, Inline: true},
        {Name: "Time Left", Value: body.Timeleft, Inline: true},
      }
      embed.Fields = inlineFields
      s.ChannelMessageSendEmbed(m.ChannelID, embed)
    }
  } else if m.Content == "!help" {
    commandList := "**Here's List of Commands:**\n\n"
    commandList += "**1. !help (to see list of commands)**\n"
    commandList += "**2. !schedule (to see contests list)**\n"
    s.ChannelMessageSend(m.ChannelID, commandList)
  }

}
