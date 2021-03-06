# dirhttps

![build](https://github.com/maetthu/dirhttps/workflows/build/badge.svg)

Spinning up an HTTPS server from current directory - the TLS equivalent of `python -m http.server`.

## Install

Use [prebuilt binaries](https://github.com/maetthu/dirhttps/releases)

or build from source

```
go get -u github.com/maetthu/dirhttps
$(go env GOPATH)/bin/dirhttps
```

## Setup

_dirhttps_ needs a certificate and corresponding key to operate. Easiest option is to use the excellent [mkcert](https://github.com/FiloSottile/mkcert) tool which creates locally trusted, self-signed development certificates. 

* Install and setup mkcert CA
* Create dirhttps certificate

``` 
$ mkdir ~/.config/dirhttps
$ cd ~/.config/dirhttps
$ mkcert -cert-file cert.pem -key-file key.pem localhost 127.0.0.1 more.hostnames.or.ips.if.needed.example.org
```

* Done

## Usage

``` 
Serving contents of current directory by HTTPS.

Usage:
  dirhttps [flags]

Flags:
      --cache           Enable client side caching
  -c, --cert string     Certificate file (default "/home/maetthu/.config/dirhttps/cert.pem")
  -d, --dump            Dump client request headers to STDOUT
  -h, --help            help for dirhttps
  -k, --key string      Key file (default "/home/maetthu/.config/dirhttps/key.pem")
  -l, --listen string   Listen address (default ":8443")
      --no-cors         Disable CORS handling
      --no-favicon      Disable default favicon delivered when no favicon.ico is present in current directory
  -q, --quiet           Don't log requests to STDOUT
      --version         version for dirhttps
```

### Examples

* Basic usage

```
$ cd /tmp
$ dirhttps  
2019/03/31 13:40:26 Listening for HTTPS connections on :8443
2019/03/31 13:40:26 Serving from directory /tmp
```

* Different listening address/port

```
$ dirhttps -l :1234
$ dirhttps -l 127.0.0.2:8443
```


## License

The default favicon served by dirhttps is generated from the [Font Awesome "code" icon](https://fontawesome.com/icons/code?style=solid) by [FontIcon](https://github.com/devgg/FontIcon) and is [licensed under the CC BY 4.0 License](https://fontawesome.com/license/free).

dirhttps is licensed under the MIT License.
