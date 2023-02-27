# extip-svr - return requester's external IP via DNS

Companion server to <https://github.com/hrbrmstr/extip>.

Barring something bad happening to the remote host, you can test it out by doing:

```
dig myip.is @ip.rudis.net
```

It works for `A`, `AAAA`, and `TXT`. Safest bet is to use the `TXT` since it'll handle either IPv4 or IPv6


```
dig myip.is TXT @ip.rudis.net
```

## Installation

```
go install -ldflags "-s -w" github.com/hrbrmstr/extip-svr@latest
```

## Disabling The systemd Stub Resolver

<https://web.archive.org/web/20201203073829/https://www.turek.dev/posts/disable-systemd-resolved-cleanly/>