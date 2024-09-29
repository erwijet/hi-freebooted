# HelloInternet FreeBooted

There is no RSS feed that includes all episodes of the [Hello Internet](https://hellointernet.fm) podcast, as the first 1-50 were removed due to squarespace limitations. Although many podcasting apps have hard-coded the episodes back in, there is no RSS feed that contains all of them.

Until now! This project scrapes the hellointernet website to collect download links and the publishes those links alongside their episode title, publish dates, and show notes as an RSS feed.

### Building

Since the scraper uses a webview to aggregate the RSS feed, the actual scraping is done as part of the build step and the produced artifact is just a static docker image.

```
$ make # builds ghcr.io/erwijet/hi-freebooted
$ docker run --rm -p 8080:8080 ghcr.io/erwijet/hi-freebooted:latest
```

---

### Deployment

The RSS feed is published here:

https://hellointernet.holewinski.dev/hellointernet.rss
