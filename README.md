# tmdbc - The Movie Database Client

Command line utility for downloading metadata and posters from The Movie Database and inserting into movie files.

Supported video containsers: MKV

The tool searches in the database automaticaly by filename. The filename can contain the year in brackets.

Supported movie name example: Alien (1979).mkv

## Compilation Prerequsites

*  Install mkvtoolnix
*  Install go

## Compilation and Installation

```bash
go get -u github.com/disintegration/imaging
go get -u github.com/ryanbradynd05/go-tmdb
go get -u github.com/vysokyj/tmdbc
go install
```

## Example Usage

```bash
tmdbc "Alien (1979).mkv" "Alice in Wonderland.mkv"
tmdbc *.mkv
```
