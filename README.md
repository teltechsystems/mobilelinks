Mobile Links
=========

A simple service that allows for the appropriate redirection of browser-specific endpoints from a single short-url

Running The Service
-------------------
```bash
go get github.com/teltechsystems/mobilelinks
mobilelinks -public_url=http://shrtnr.io -listen=:8000
```

Using The Service
-----------------
Once the service is running, creating a short url is as simple as a curl!

```bash
curl -X POST http://localhost:8000/create \
    -d"default=http://www.google.com" \
    -d"android=http://www.android.com" \
    -d"ios=http://www.apple.com"
> http://shrtnr.io/a3g

curl --header "User-Agent: Android" http://shrtnr.io/a3g # Redirects To android.com!
curl --header "User-Agent: iPod" http://shrtnr.io/a3g # Redirects To apple.com!
curl --header "User-Agent: Safari" http://shrtnr.io/a3g # Redirects To google.com!
```