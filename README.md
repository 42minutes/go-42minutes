# 42minutes

[![Build Status](https://travis-ci.org/42minutes/go-42minutes.svg?branch=chore%2Fadd-travis-ci)](https://travis-ci.org/42minutes/go-42minutes)
[![Go Report Card](https://goreportcard.com/badge/github.com/42minutes/go-42minutes)](https://goreportcard.com/report/github.com/42minutes/go-42minutes)

42minutes is a collection of tools to manage your tv-series collection.  
It consists of a client, server, and user interface.

In addition to these there is a single user daemon that wraps server and
client into a single binary, requires no authentication, and meant to be
used by a single user.

__The standalone daemon is the first thing 42minutes will focus on.__

## client

The client should be able to

* Watch the user's local library and identify existing series and episodes.
* Upload the user's library metadata to the 42minutes server and keep it up
  to date.
* Download episodes it is told to by the 42minutes server, either using an
  internal downloader or an external application.

The client only sends file metadata to the server, and nothing else. It will
only send the file's relative path to the watched directory, file size, and
file checksum.

The server is the one responsible for understanding which series and episode
this file belongs to.

## Server

The server's core responsibilities are to

* Keep a listing of each user's tv-series, seasons, and episodes.
* Allow the users to specify which series they want to check for completenes
  and new episodes.
* Create a listing of episodes to download for the 42 minutes client, or
  other RSS downloader.

To accomplish these tasks the server will be able understand which episodes
the user is missing as well as to find sources where they can be downloaded
from.

* Identify series name, season, and episode from file name or path.
* Match the files against a tv-series provider such as trakt, tvdb, or other.
* Identify missing episodes.
* Watch for newly released episodes for tv series users are watching. 
* Find missing or new episode download sources from torrents, nbz, or other.

## User interface

We need a "simple" user interface to allow the users to

* View their libraries and mark which series they want to download missing or
  new episodes.