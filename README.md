# ujsonize

This is a cli tool that converts Go's `url.Values`(i.e. `map[string][]string`)
and `json` to each other.

**That's all.**

```bash
$ echo "foo=bar" | ujsonize encode
{"foo":["bar"]}
$ echo "foo=bar" | ujsonize encode | ujsonize decode
foo=bar
```

# Installation

```
$ go get github.com/imishinist/ujsonize
```

# Author

Taisuke Miyazaki (a.k.a imishinist)
