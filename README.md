# www

www is a tool for piping stdout across the web.

## Installation

```
$ go get github.com/Bowery/www
```

## Usage

```
$ echo "Hello World" | www slack --channel=#general
$ cat sample_code.js | www gist
```

## Providers

- Slack
- Gist
- GMail

More coming!

## Contributing

All contributions are welcomed and certainly encouraged! Each `provider` abides by the `Provider` interface defined in `providers/provider.go`.
