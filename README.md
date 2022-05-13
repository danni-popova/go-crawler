# Danni's Web Crawler

## How to run

You can run it by building it and then running the binary and providing it with the starting URL to crawl.

Note: Please provide the URL as shown below - starting with `https://` without a trailing slash. 

```shell
go build crawler
./crawler --starting-url "https://improbable.io"
```

## How to test
```shell
go test -v ./...
```

