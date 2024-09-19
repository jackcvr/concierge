# Concierge

TODO

## Installation

See [Releases](https://github.com/jackcvr/reverssh/releases)

## Usage

```shell
Usage of concierge:
  -a value
    	Endpoint in format 'url:host:port' (e.g. /ssh:localhost:22)
  -b string
    	Local address to listen on (default "0.0.0.0:80")
  -crt string
    	Crt file for TLS
  -f string
    	Log file (default stdout)
  -key string
    	Key file for TLS
  -q	Do not print anything
  -t duration
    	Timeout for new connections (default 2s)
  -v	Verbose mode
```

## Examples

Configure SSH server binding it to localhost.

Start Concierge HTTPS server which replies on `/ssh` URL: 

```shell
$ sudo concierge -a /ssh:localhost:22 -crt server.crt -key server.key
{"time":"2024-09-19T22:38:16.973240816+03:00","level":"INFO","msg":"http/listening","addr":{"IP":"0.0.0.0","Port":443,"Zone":""}}
{"time":"2024-09-19T22:38:18.867798191+03:00","level":"INFO","msg":"http","remoteAddr":{"IP":"127.0.0.1","Port":45304,"Zone":""},"agent":"curl/7.68.0","method":"GET","url":"/ssh"}
{"time":"2024-09-19T22:38:18.867841689+03:00","level":"INFO","msg":"tcp/listening","addr":{"IP":"::","Port":39575,"Zone":""}}
{"time":"2024-09-19T22:38:18.870319064+03:00","level":"INFO","msg":"connected","laddr":{"IP":"127.0.0.1","Port":39575,"Zone":""},"raddr":{"IP":"127.0.0.1","Port":48886,"Zone":""}}
{"time":"2024-09-19T22:38:18.870525647+03:00","level":"INFO","msg":"connected","laddr":{"IP":"127.0.0.1","Port":46658,"Zone":""},"raddr":{"IP":"127.0.0.1","Port":22,"Zone":""}}
```

Connect to your server:

```shell
$ ssh root@localhost -p $(curl -sk https://localhost/ssh)
```

## License

[MIT](https://spdx.org/licenses/MIT.html) 