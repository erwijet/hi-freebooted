package lib

import (
	"fmt"
	"log"
	"strings"
	"time"

	webview "github.com/webview/webview_go"
)

type HelloInternetArchiveInfo struct {
	Title     string
	Url       string
	PubDate   time.Time
	ShowNotes string
}

type HiPodcast struct {
	Episodes []HelloInternetArchiveInfo
}

func CrawlHelloInternetArchive(num int) (ret *HelloInternetArchiveInfo) {
	ret = new(HelloInternetArchiveInfo)
	w := webview.New(true)
	defer w.Destroy()

	t := NewTaskmaster(&w)

	if num == 116 {
		return // this is the 12 days of christmas nonsense that is ALL encoded as "one" episode :sigh:
	}

	t.Task("crawlMp3", func(url string) {
		ret.Url = strings.Split(url, "?")[0] // remove the trailing ?download=true
	})

	t.Task("crawlTitle", func(title string) {
		ret.Title = title
	})

	t.Task("crawlShownotes", func(shownotes string) {
		ret.ShowNotes = shownotes
	})

	t.Task("crawlPubDateStr", func(pubdatestr string) {
		parsed, err := time.Parse("2006-01-02", pubdatestr)
		if err != nil {
			log.Fatal(err)
		}
		ret.PubDate = parsed
	})

	w.Bind("catch", func(err string) {
		log.Printf("When processing URL: %v", ResolveHelloInternetArchiveUrl(num))
		log.Fatal(err)
	})

	w.Navigate(ResolveHelloInternetArchiveUrl(num))

	w.Init(`
		window.onload = function() {
			try {
				var title = document.querySelector('.entry-title').textContent;
				window.crawlTitle(title);

				var shownotes = document.querySelector('div[data-block-type="44"] .sqs-block-content')?.innerHTML ?? '<div />';
				window.crawlShownotes(shownotes);

				if (title == 'Six Geese A-laying')
					window.crawlMp3('https://traffic.libsyn.com/hellointernet/Six_Geese_A-laying.mp3');
				else if (title == 'Eleven Pipers Piping')
					window.crawlMp3('https://traffic.libsyn.com/hellointernet/Eleven_Pipers_Piping.mp3');
				else {
					var link = document.querySelector('.sqs-audio-embed div.secondary-controls > div.download > a').href;
					window.crawlMp3(link);
				}

				var pubdatestr = document.querySelector('.published').getAttribute('datetime');
				window.crawlPubDateStr(pubdatestr);
			} catch (ex) { window.catch(ex.toString()) }
		};
	`)

	go func() {
		t.WaitForAllTasks()
		w.Terminate()
	}()

	w.Run() // blocks until w.Terminate() is called

	return
}

func ResolveHelloInternetArchiveUrl(num int) string {
	if num == 93 {
		return "https://www.hellointernet.fm/podcast/2017/12/24/hi-star-wars-the-last-jedi-christmas-special"
	}

	if num == 100 {
		return "https://www.hellointernet.fm/podcast/onehundred"
	}

	if num == 115 {
		return "https://www.hellointernet.fm/podcast/2018/12/23/hi-115-pink-flamingo"
	}

	if num == 107 {
		return "https://www.hellointernet.fm/podcast/hi-107-one-year-of-weird"
	}

	if num == 122 {
		return "https://www.hellointernet.fm/podcast/2019/4/24/hi-122-wax-cylinders"
	}

	if num == 125 {
		return "https://www.hellointernet.fm/podcast/2019/6/30/hi-125-the-spice-must-flow"
	}

	if num == 133 {
		return "https://www.hellointernet.fm/podcast/2019/12/25/star-wars-the-rise-of-skywalker-hello-internet-christmas-special"
	}

	if num == -1 {
		return "https://www.hellointernet.fm/podcast/2018/12/25/a-partridge-in-a-pear-tree"
	}

	if num == -2 {
		return "https://www.hellointernet.fm/podcast/2018/12/26/two-turtle-doves"
	}

	if num == -3 {
		return "https://www.hellointernet.fm/podcast/2018/12/27/three-french-hens"
	}

	if num == -4 {
		return "https://www.hellointernet.fm/podcast/2018/12/28/four-calling-birds"
	}

	if num == -5 {
		return "https://www.hellointernet.fm/podcast/2018/12/29/five-gold-rings"
	}

	if num == -6 {
		return "https://www.hellointernet.fm/podcast/2018/12/30/six-geese-a-laying"
	}

	if num == -7 {
		return "https://www.hellointernet.fm/podcast/2018/12/31/seven-swans-a-swimming"
	}

	if num == -8 {
		return "https://www.hellointernet.fm/podcast/2019/1/1/8-maids-a-milking"
	}

	if num == -9 {
		return "https://www.hellointernet.fm/podcast/2019/1/2/nine-ladies-dancing"
	}

	if num == -10 {
		return "https://www.hellointernet.fm/podcast/2019/1/3/ten-lords-a-leaping-really"
	}

	if num == -11 {
		return "https://www.hellointernet.fm/podcast/2019/1/4/eleven-pipers-piping"
	}

	if num == -12 {
		return "https://www.hellointernet.fm/podcast/2019/1/5/twelve-drummers-drumming"
	}

	return fmt.Sprintf("https://www.hellointernet.fm/podcast/%v", num)
}
