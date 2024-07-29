[![GitHub release](https://img.shields.io/github/release/sgaunet/http-echo.svg)](https://github.com/sgaunet/http-echo/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/sgaunet/http-echo)](https://goreportcard.com/report/github.com/sgaunet/http-echo)
![GitHub Downloads](https://img.shields.io/github/downloads/sgaunet/http-echo/total)
[![Maintainability](https://api.codeclimate.com/v1/badges/82203a820d424b9b12e9/maintainability)](https://codeclimate.com/github/sgaunet/http-echo/maintainability)
[![GoDoc](https://godoc.org/github.com/sgaunet/http-echo?status.svg)](https://godoc.org/github.com/sgaunet/http-echo)
[![License](https://img.shields.io/github/license/sgaunet/http-echo.svg)](LICENSE)

# http-echo

Just a webserver that prints miscellaneous informations about the http request. It can be used to play with docker and kubernetes.

```
$ docker-compose up -d
...
$  curl -i http://localhost:8080/hello?t=toto
HTTP/1.1 200 OK
Date: Thu, 10 Mar 2022 16:21:11 GMT
Content-Length: 374
Content-Type: text/plain; charset=utf-8

r.URL.Query() :
t => [toto]
End r.URL.Query()

Headers :
User-Agent => [curl/7.68.0]
Accept => [*/*]
End Headers

r.PostForm: map[]
r.RequestURI: /hello?t=toto
r.URL.Query(): map[t:[toto]]
r.Form: map[]
body: 
url.ParseQuery(string(body))End url.ParseQuery(string(body))

Method: GET
Host: localhost:8080
Proto: HTTP/1.1
Remote Addr: 172.18.0.1:43454
Hostname: 633b9fd87484
```
