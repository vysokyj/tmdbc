.PHONY: all

all:
	go get -u github.com/disintegration/imaging
	go get -u github.com/ryanbradynd05/go-tmdb
	mkdir -p build/linux_386
	GOOS=linux GOARCH=386 go build -o build/linux_386/tmdbc ./
	mkdir -p build/linux_amd64
	GOOS=linux GOARCH=amd64 go build -o build/linux_amd64/tmdbc ./
	mkdir -p build/windows_386
	GOOS=windows GOARCH=386 go build -o build/windows_386/tmdbc.exe ./
	mkdir -p build/windows_amd64
	GOOS=windows GOARCH=amd64 go build -o build/windows_amd64/tmdbc.exe ./
