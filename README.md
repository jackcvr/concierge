# Concierge

Concierge is an unconventional TCP Reverse Proxy designed to hide any TCP services "behind" HTTP server.

It works by starting an HTTP server that dynamically creates TCP listeners 
upon requests to a predefined URLs. It responds with the new port number, which the client must connect to.
Traffic from the first successful connection is then redirected to another(usually internal) 
service on a designated port. The opened port exists only for that one connection.

Other features:
- Only the original requester’s IP is allowed to connect to the provided port.
- Requests to undefined URLs are tarpitted unless `-ntp` argument is provided.
- Can be bound any number of URLs to endpoints.

## Installation

See [Releases](https://github.com/jackcvr/concierge/releases)

## Usage

```shell
Usage of concierge:
  -c string
    	Path to TOML config file (default "/etc/concierge/concierge.toml")
```

### concierge.toml sample
```toml
tz = "Europe/Vilnius"
bind = "0.0.0.0:80"
crtfile = ""  # set crt and key files to enable TLS
keyfile = ""
quiet = false
verbose = false
timeout = 2  # seconds
notarpit = false

[endpoints]
"/ssh" = "0.0.0.0:2222"
```


## Example

On remote machine:
- configure SSH server to bind to localhost.
- start the `concierge` on HTTPS, which responds to requests made to the `/ssh` path:

```shell
$ sudo concierge -c concierge.toml
{"time":"2024-09-21T12:27:36.180365398+03:00","level":"INFO","msg":"http/listening","addr":{"IP":"0.0.0.0","Port":80,"Zone":""}}
{"time":"2024-09-21T12:27:42.710054064+03:00","level":"INFO","msg":"http/connected","remoteAddr":{"IP":"127.0.0.1","Port":58664,"Zone":""},"agent":"curl/7.68.0","method":"GET","url":"/ssh"}
{"time":"2024-09-21T12:27:42.71009367+03:00","level":"INFO","msg":"tcp/listening","addr":{"IP":"::","Port":46381,"Zone":""}}
{"time":"2024-09-21T12:27:42.710104628+03:00","level":"INFO","msg":"http/closed","remoteAddr":{"IP":"127.0.0.1","Port":58664,"Zone":""},"url":"/ssh","lifetime":0}
{"time":"2024-09-21T12:27:42.714576373+03:00","level":"INFO","msg":"tcp/connected","laddr":{"IP":"127.0.0.1","Port":46381,"Zone":""},"raddr":{"IP":"127.0.0.1","Port":33054,"Zone":""}}
{"time":"2024-09-21T12:27:42.714741091+03:00","level":"INFO","msg":"tcp/closed","addr":{"IP":"::","Port":46381,"Zone":""}}
{"time":"2024-09-21T12:27:42.715085155+03:00","level":"INFO","msg":"tcp/connected","laddr":{"IP":"127.0.0.1","Port":43768,"Zone":""},"raddr":{"IP":"127.0.0.1","Port":22,"Zone":""}}
```

On local machine:
- connect to your SSH server (assuming the IP of your remote machine is 8.8.8.8):

```shell
$ ssh root@8.8.8.8 -p $(curl -sk https://8.8.8.8/ssh)
```

## License

[MIT](https://spdx.org/licenses/MIT.html) 