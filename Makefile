docker : hellointernet.rss
	docker build -t ghcr.io/erwijet/hi-freebooted:latest .

hellointernet.rss : freebooter
	./freebooter \
		-o hellointernet.rss \
		-t "HelloInternet (freebooted)" \
		-l https://hellointernet.holewinski.dev/hellointernet.rss

freebooter :
	go build

clean :
	rm freebooter *.rss

.PHONY : docker clean
