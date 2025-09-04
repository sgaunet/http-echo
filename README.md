[![GitHub release](https://img.shields.io/github/release/sgaunet/http-echo.svg)](https://github.com/sgaunet/http-echo/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/sgaunet/http-echo)](https://goreportcard.com/report/github.com/sgaunet/http-echo)
![GitHub Downloads](https://img.shields.io/github/downloads/sgaunet/http-echo/total)
![Coverage Badge](https://raw.githubusercontent.com/wiki/sgaunet/http-echo/coverage-badge.svg)
[![linter](https://github.com/sgaunet/http-echo/actions/workflows/coverage.yml/badge.svg)](https://github.com/sgaunet/http-echo/actions/workflows/coverage.yml)
[![coverage](https://github.com/sgaunet/http-echo/actions/workflows/coverage.yml/badge.svg)](https://github.com/sgaunet/http-echo/actions/workflows/coverage.yml)
[![Snapshot Build](https://github.com/sgaunet/http-echo/actions/workflows/snapshot.yml/badge.svg)](https://github.com/sgaunet/http-echo/actions/workflows/snapshot.yml)
[![Release Build](https://github.com/sgaunet/http-echo/actions/workflows/release.yml/badge.svg)](https://github.com/sgaunet/http-echo/actions/workflows/release.yml)
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

## Project Status

🟨 **Maintenance Mode**: This project is in maintenance mode.

While we are committed to keeping the project's dependencies up-to-date and secure, please note the following:

- New features are unlikely to be added
- Bug fixes will be addressed, but not necessarily promptly
- Security updates will be prioritized

## Issues and Bug Reports

We still encourage you to use our issue tracker for:

- 🐛 Reporting critical bugs
- 🔒 Reporting security vulnerabilities
- 🔍 Asking questions about the project

Please check existing issues before creating a new one to avoid duplicates.

## Contributions

🤝 Limited contributions are still welcome.

While we're not actively developing new features, we appreciate contributions that:

- Fix bugs
- Update dependencies
- Improve documentation
- Enhance performance or security

If you're interested in contributing, please read our [CONTRIBUTING.md](link-to-contributing-file) guide for more information on how to get started.

## Support

As this project is in maintenance mode, support may be limited. We appreciate your understanding and patience.

Thank you for your interest in our project!