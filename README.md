# extip-svr - return requester's external IP via DNS

Companion server to <https://github.com/hrbrmstr/extip>.

Barring something bad happening to the remote host, you can test it out by doing:

```
dig myip.is @ip.rudis.net
```

It works for `A`, `AAAA`, and `TXT`. Safest bet is to use the `TXT` since it'll handle either IPv4 or IPv6.


```
dig myip.is TXT @ip.rudis.net
```

## Installation

```
go install -ldflags "-s -w" github.com/hrbrmstr/extip-svr@latest
```

## Usage

You can configure the server to bind to a particular port and use a TLD different than the default `is.` via flags or environment variables.

```
$ ./extip-svr --help
Usage: extip-svr [--port PORT] [--tld TLD] [--quiet]

Options:
  --port PORT, -p PORT   bind port [default: 53, env: EXTIP_PORT]
  --tld TLD, -t TLD      TLD to handle [default: is., env: EXTIP_TLD]
  --quiet, -q            Disable log messages
  --help, -h             display this help and exit
```

## Disabling The systemd Stub Resolver

<https://web.archive.org/web/20201203073829/https://www.turek.dev/posts/disable-systemd-resolved-cleanly/>