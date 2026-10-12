# url-code-checker
A simple app to check HTTP status codes of a group of URLs.  
Built with Go for maximum speed, can be distributed as a Docker image.

[![Go Version](https://img.shields.io/badge/Go-1.27-00ADD8?logo=go)](https://go.dev)
[![Docker](https://img.shields.io/badge/Docker-ready-2496ED?logo=docker)](https://www.docker.com)
[![License](https://img.shields.io/badge/license-MIT-green)](./LICENSE)

## Features
- Check multiple URLs concurrently
- Know when a site is down or up
- Streamed JSON output
- Small Docker image, based on Alpine
- No external dependencies

## Quick Start

### With Docker
```bash
docker build -t url-checker:latest .
docker run -p 8080:8080 url-checker
```
Curl or Postman can be used to access the API `http://localhost:8080`.

### From source
```bash
git clone https://github.com/ak-softwaredev/url-code-checker.git
cd url-code-checker
make run
```

## Usage
### Array of URLs
```bash
curl -X POST -N http://178.253.45.137:8080/api/v1/check \
  -H 'Content-Type: application/json'
  -d '{"links":["https://example.com","https://example.org","https://example.edu]}'
```
Response:
```json
{"url":"https://example.edu","code":200,"ms":50}
{"url":"https://example.com","code":200,"ms":65}
{"url":"https://example.org","code":200,"ms":100}
```