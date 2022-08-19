.PHONY: build-pi deploy-pi

build-pi:
	GOOS=linux GOARCH=arm GOARM=5 go build .

deploy-pi: build-pi
	scp gasmeter logpi:~
