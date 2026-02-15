#!/usr/bin/env zsh
go run main.go download --link /data/post.link --proxy http://192.168.5.2:8889 --cookie /data/pornhub.cookie
go run main.go whisper --level medium.en --location /data/models --language English --root /data --format srt
go run main.go trans --root /data/videos --proxy http://192.168.5.2:8889
go run main.go merge --root /data/videos