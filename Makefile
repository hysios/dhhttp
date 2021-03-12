REMOTE_SERVER=x6

build: build-linux

build-linux:
	@GOOS=linux GOARCH=arm GOARM=7 go build -o bin/dhhttp-linux ./example

deploy:
	@ssh $(REMOTE_SERVER) pkill dhttp 
	@scp bin/dhhttp-linux $(REMOTE_SERVER):/home