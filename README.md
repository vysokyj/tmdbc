# tmdbc - The Movie Database Client

[![Build Status](https://secure.travis-ci.org/vysokyj/tmdbc.svg?branch=master)](http://travis-ci.org/vysokyj/tmdbc)

Command line utility for downloading metadata and posters from [The Movie Database](https://www.themoviedb.org) and inserting into movie files.

Supported video containsers: MKV

The tool searches in the database automaticaly by filename. The filename can contain the year in brackets.

Supported movie name example: Alien (1979).mkv

## Compilation and Installation

### Prerequsites

*  Install [mkvtoolnix](https://mkvtoolnix.download/)
*  Install [GO](https://golang.org/)

### Commands

```bash
go get -u github.com/disintegration/imaging
go get -u github.com/ryanbradynd05/go-tmdb
go get -u github.com/vysokyj/tmdbc
go install
```

## Usage

Call utilty with one or more movie file path arguments.

### Examples

#### Example 1:

```bash
tmdbc "Alien (1979).mkv" "Alice in Wonderland.mkv"
```

#### Example 2:

```bash
tmdbc *.mkv
```
