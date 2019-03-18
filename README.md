# tmdbc - The Movie Database Client

[![Build Status](https://secure.travis-ci.org/vysokyj/tmdbc.svg?branch=master)](http://travis-ci.org/vysokyj/tmdbc)

Command line utility downloads metadata and posters from [The Movie Database](https://www.themoviedb.org) and inserts them into movie files.

## Main Features

* Downloads metadata in predefined language (en, de, fr....)
* Uses your own API key from [The Movie Database](https://www.themoviedb.org).
* Supported video containers: MKV (for now)

## Compilation and Installation

### Prerequisites

* Install [mkvtoolnix](https://mkvtoolnix.download/)
* Install [GO](https://golang.org/)
* Obtain your API key from [The Movie Database](https://www.themoviedb.org)
* Support year in filaname - supported format example: `Alien (1979).mkv`

Application will ask for the API key and preferred language when you start it for the first time. Configuration is kept in a `.tmdbc` file stored in your home directory.

### Installation


```bash
go get -u github.com/disintegration/imaging
go get -u github.com/ryanbradynd05/go-tmdb
go get -u github.com/vysokyj/tmdbc
go install
```

or use make for cross compilation

```bash
make
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
