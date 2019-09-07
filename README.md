# gomediaserver

![image.png](https://github.com/leonardoderoy/phpmediaserver/blob/master/imgs/logo/1.png?raw=true)

## Screenshots

![gomediaserver](https://media.giphy.com/media/8A7qPRuF8jzRcaQiyj/giphy.gif)

[IMG1](http://i67.tinypic.com/33opaiv.png)

## Description

FOR TESTING, for now only few options

A bunch of utilities for:
 - Html5 web player with go + sqlite + jquery + ffmpeg
 - WebPlayer support for audio and subs tracks selector
 - Filter list by search and list by genres
 - Easy configuration with `config.toml`
 - Admins and player users
 - Realtime ffmpeg transcoding of any type of video supported by ffmpeg, not needed to reencode before play or create temp files
 - Identify media files thanks to: [pymediaident](https://github.com/EsTass/pymediaident), www.filebot.net, www.omdbapi.com, www.thetvdb.com (cron, manual and helped).
 - Country IP block thanks to www.geoplugin.net.
 - IP whitelist/blacklist (autoban non included countrys)
 - Media info in configured language (if possible).
 - Logo thanks to [leonardoderoy](https://github.com/leonardoderoy)
 - Groups by premiere, continue, recomended and last added (frontal page)
 - IPTV List and import
 - Poster list with search by genres, actors, years or rating (complex search)
 - IPTV from urls on cron and import
 - Extract files on cron
 - Filtered remove files to recover extra free space (manual, helped and cron)
 - Clean duplicates by quality with safe seeding (min days to seed) and max filesize to maintanin
 
## Working on
 - Search and download new media from web adding scrappers to configuration (youtube, elinks, magnets, torrents and dd supported, cron or manual with any external program like transmission, jdownloader, amule, qbittorent, etc).
 - Stop adding downloads on min space config
 - Mini dlna server
 - [Kodi pluging](https://github.com/EsTass/phpmediaserver-kodi)
 - Multilanguaje WebUI
 
## Default User (Important: change pass on first login)
 - User: admin
 - Pass: admin01020304
 
## Install
 - Download `https://github.com/EsTass/gomediaserver/archive/master.zip` 
 - Extract and edit `config.toml`
 - Build or use `gms` or `gms.exe`
 - run `./gms`(linux) or `gms.exe` (windows, sorry not tested)
 - Database compatible with [phpmediaserver](https://github.com/EsTass/phpmediaserver) can be replaced

## Build

 - Download needed from `main.go` imports header (commented are used by other files but needed)
 - `go build .`

## Needed
 - [ffmpeg and ffprobe](https://ffmpeg.org/)
 
## Recomended
- [pymediaident](https://github.com/EsTass/pymediaident)
- [Filebot](https://www.filebot.net)
- [omdbapi APIKEY](https://www.omdbapi.com)
- [www.thetvdb.com APIKEY](https://www.thetvdb.com)
- [googler](https://pypi.python.org/pypi/googler)
- [ddgr](https://github.com/jarun/ddgr)
