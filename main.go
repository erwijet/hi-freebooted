package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	podcasts "github.com/jbub/podcasts"
	gofeed "github.com/mmcdole/gofeed"
	progressbar "github.com/schollz/progressbar/v3"
	lib "holewinski.dev/freebooter/lib"
)

const TOTAL_EPISODES = 136

func main() {
	var outfile string
	flag.StringVar(&outfile, "o", "", "The file to write the RSS feed to")
	flag.Parse()

	if outfile == "" {
		flag.Usage()
		os.Exit(1)
	}

	o, err := os.Create(outfile)
	defer o.Close()

	if err != nil {
		log.Panic(err)
	}

	var pod = new(lib.HiPodcast)
	bar := progressbar.Default(TOTAL_EPISODES + 12) // for the 12 days of christmas stuff
	bar.Describe("Scraping Episodes")

	for i := 1; i <= TOTAL_EPISODES; i++ {
		pod.Episodes = append(pod.Episodes, *lib.CrawlHelloInternetArchive(i))
		bar.Add(1)
	}

	// and then here, we represent a negative episode number as a "bonus" episode.
	// each negative episode number is handled explicitly as a "bonus episode id"
	// by the crawler
	for i := 1; i <= 12; i++ {
		pod.Episodes = append(pod.Episodes, *lib.CrawlHelloInternetArchive(-i))
		bar.Add(1)
	}

	bar.Finish()

	fmt.Println("Fetching original rss feed metadata...")

	fp := gofeed.NewParser()
	orig, _ := fp.ParseURL("https://www.hellointernet.fm/podcast?format=rss")

	p := &podcasts.Podcast{
		Title:       "Hello Internet (freebooted)",
		Description: orig.Description,
		Language:    "EN",
		Link:        "https://hi.holewinski.dev/.rss",
		Copyright:   orig.Copyright,
	}

	fmt.Println("Constructing feed...")

	for _, item := range pod.Episodes {
		p.AddItem(&podcasts.Item{
			Title: item.Title,
			GUID:  item.Url,
			Image: &podcasts.ItunesImage{
				Href: orig.Image.URL,
			},
			Summary: &podcasts.ItunesSummary{
				Value: item.ShowNotes,
			},
			PubDate: &podcasts.PubDate{Time: item.PubDate},
			Enclosure: &podcasts.Enclosure{
				URL:  item.Url,
				Type: "mp3",
			},
		})
	}

	feed, err := p.Feed(
		podcasts.Author(strings.Join(lib.Map(orig.Authors, func(each *gofeed.Person) string {
			return each.Name
		}), ", ")),
		podcasts.Block,
		podcasts.Complete,
		podcasts.Summary(orig.Description),
		podcasts.Subtitle(orig.ITunesExt.Subtitle),
		podcasts.Image(orig.Image.URL),
	)

	if err != nil {
		log.Fatal(err)
	}

	feed.Write(o)

	// http.HandleFunc("/.rss", func(w http.ResponseWriter, r *http.Request) {
	// 	feed.Write(w)
	// })

	// fmt.Println("Servering RSS feed on 0.0.0.0:8080/.rss")
	// http.ListenAndServe(":8080", nil)
}
