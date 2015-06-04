# www

www is a tool for piping stdout across the web.

## Installation

```
$ go get github.com/Bowery/www
```

## Usage

First you'll need to setup your tokens, authentication, etc. if the provider requires it.

```
$ www setup <provider>
```

Once that's done you can easily pipe any command to that provider:

```
$ echo "Hello World" | www slack --channel=#general
$ cat sample_code.js | www gist
```

You can even chain them together. Here we create a gist and then email it:

```
$ cat www.go | www gist | www gmail -to=hello@bowery.io
```

## Providers

- Slack
- Gist
- GMail

More coming!

## Contributing

All contributions are welcomed and certainly encouraged! Each `provider` abides by the `Provider` interface defined in `providers/provider.go`.
