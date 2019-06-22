# Goms

Extensible tool to generate microservices in Go using [go-kit](https://gokit.io/), by specified service interface. It is primarily inspired by [microgen](https://github.com/devimteam/microgen).
The goal is to generate all the code you need to write but has nothing to do with your business logic.

## Install
```
go get -u github.com/wlMalk/goms
```

## Usage
``` sh
goms
```
goms tool will look for any service interface declared in `service.go` file inside `CWD`.
