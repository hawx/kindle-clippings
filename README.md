# kindle-clippings

Contains two tools:

1. kindle-clippings -- reads the "My Clippings.txt" into a JSON format.
2. kindle-clippings-cookie -- reads the JSON format produced by 1) and writes
   it into a format usable by [fortune][].

```sh
$ go get hawx.me/code/kindle-clippings/cmd/...
$ kindle-clippings /media/my/Kindle > clippings.json
$ [sudo] kindle-clippings-cookie /usr/share/games/fortunes/clippings < clippings.json
```

[fortune]: https://en.wikipedia.org/wiki/Fortune_%28Unix%29
