go-simple-server
================
Simple HTTP server serving domain redirects and simple pages, allowing custom headers.   

Configuration
-------------
### config.yaml
```
listen: ":8888"
loginmatched: true
hosts:
  foobar.example:
    body: "Hello World"
  foo.example:
    redirect: "https://bar.example/foo"
    log: true
  bar.example:
    redirect: "https://foobar.example{{ .RequestURI }}"
  root.example:
    redirect: "https://domain.example{{ .RequestURI }}"
    code: 302
    headers:
      Strict-Transport-Security: "max-age=63072000; includeSubDomains; preload"
```

### Defaults
- Default return code is 200 for an existing host.
- When a redirect rule is present, defaults to HTTP 301.
- Host `status` defaults to HTTP 200 with body response `ok` (for healthchecks).
- Unknown hosts return HTTP 404.
- Listens on `0.0.0.0:80`
