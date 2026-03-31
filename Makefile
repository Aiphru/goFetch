install:
	go build 
	sudo cp gofetch /usr/local/bin
uninstall:
	sudo rm /usr/local/bin/gofetch