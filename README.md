# get-text

![Go](https://github.com/airenas/get-text/workflows/Go/badge.svg) [![Coverage Status](https://coveralls.io/repos/github/airenas/get-text/badge.svg?branch=main)](https://coveralls.io/github/airenas/get-text?branch=main) [![Go Report Card](https://goreportcard.com/badge/github.com/airenas/get-text)](https://goreportcard.com/report/github.com/airenas/get-text) ![CodeQL](https://github.com/airenas/get-text/workflows/CodeQL/badge.svg)

Service to extract text from various formats - wrapper for [ebook-convert](https://manual.calibre-ebook.com/generated/en/ebook-convert.html)

## Building docker image

```bash
make build-docker
```

## Running with compose

Starting service on default port `8003`

```bash
(cd examples/docker-compose && make start)
```

## Testing

Epub file can be converted with:

```bash
curl -X POST http://localhost:8003/text -H 'content-type: multipart/form-data' -F file=@examples/docker-compose/sample.epub
```

Expected result is:

```json
{"text":"Hello book\n\n\n\n\n\n"}
```

## Stopping/removing docker container

```bash
(cd examples/docker-compose && make stop)
```

---

### License

Copyright © 2021, [Airenas Vaičiūnas](https://github.com/airenas).
Released under the [The 3-Clause BSD License](LICENSE).

---
