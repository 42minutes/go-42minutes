# 42minutes

[![Build Status](https://travis-ci.org/42minutes/go-42minutes.svg?branch=chore%2Fadd-travis-ci)](https://travis-ci.org/42minutes/go-42minutes)
[![Go Report Card](https://goreportcard.com/badge/github.com/42minutes/go-42minutes)](https://goreportcard.com/report/github.com/42minutes/go-42minutes)

42minutes is a collection of tools to manage your tv-series collection.  
It consists of a client, server, and user interface.

In addition to these there is a single user daemon that wraps server and
client into a single binary, requires no authentication, and meant to be
used by a single user.

__The standalone daemon is the first thing 42minutes will focus on.__

## standalone

The standalone client should be able to

* Watch the user's local library and identify existing series and episodes.
* Find missing and new episode for the user's tv series.
* Find magnet links for the missing and new episodes, and add them as `.magnet`
  files in a watch folder that torrent clients can use.
* Create a listing of episodes to download for Torrent RSS downloaders.
* Exposes a simple interactive command line interface with the following
  commands:
  * `list` - lists all shows in user's library
  * `add show-name` - adds show based on name in user's library
  * `watch dir-path` - starts watching a directory recursively (if dir-path
    default to config value)
  * `diff` - runs diff to find missing episodes

### HTTP API

The standalone client comes with an HTTP API for managing the user library.  
More info on the [42minutes HTTP API docs](http://docs.42minutes.apiary.io).

### Getting started

Clone the repo, make sure you have [glide](https://github.com/Masterminds/glide),
do a `glide install` and you are set. The standalone client uses sqlite as its 
primary storage.

Copy the `cmd/standalone/config-sample.json` as `cmd/standalone/config.json`
and modify it to match your settings. Trakt.tv client id can be left to the default.

You can now run the standalone client by `cd cmd/stadalone && go run *.go`.

The client will start an HTTP API on `http://localhost:8081` and you will be
presented with a promt `>>>` where you can now try any of the available commands.

At this point you can `watch` to the client can go through your episodes.
And once this is done, `diff` to find your missing episodes, and their infohashes.