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
