package scraper

import (
  "strings"
  "github.com/PuerkitoBio/goquery"
  "net/http"
  "log"
)

var (
  ContestSite = []string{
    "atcoder",
    "yukicoder",
    "codeforces",
    "codechef",
    "dmoj",
    "codejam",
    "usaco",
    "icpc",
    "tlx.toki",
    "leetcode"}
)

type Contest struct {
  Title string
  Link string
  Status string
  Duration string
  Timeleft string
}

func Scrape() []Contest {
  res, err := http.Get("https://clist.by/")
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()

  if res.StatusCode != 200 {
    log.Fatalf("Status code error: %d %s", res.StatusCode, res.Status)
  }
  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    log.Fatal(err)
  }
  contests := make([]Contest, 0)
  doc.Find("#contests").Children().Each(func(i int, sel  *goquery.Selection) {
    row := new(Contest)
    row.Title = sel.Find(".contest_title").Children().Text()
    row.Link, _= sel.Find(".contest_title").Children().Attr("href")
    status, _ := sel.Attr("class")
    if len(strings.Split(status, " ")) > 2 {
      row.Status = strings.ToUpper(strings.Split(status, " ")[2])
    } else {
      row.Status = ""
    }
    row.Duration = sel.Find(".duration").Text()
    row.Timeleft = sel.Find(".timeleft").Text()
    for _, val := range ContestSite {
      if row.Link != "#" && row.Link != "" && strings.Contains(row.Link, val) {
        contests = append(contests, *row)
      }
    }
  })
  return contests
}

