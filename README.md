# goConnectPool  
[![Build Status](http://img.shields.io/travis/fatih/pool.svg?style=flat-square)](https://travis-ci.org/dubin555/goConnectPool)
[![Go Report Card](https://goreportcard.com/badge/github.com/dubin555/goConnectPool)](https://goreportcard.com/report/github.com/dubin555/goConnectPool)
![Project Status](https://img.shields.io/badge/version-1.0-green.svg)

## What is goConnectPool
A go net.conn pool for Golang. Inspired by the Pool(fatih/pool), add some code and function for limit the active connection numbers.

## Install and Usage

```bash
go get github.com/dubin555/goConnectPool
```

Please make sure you use the master branch code.

## Example

```go
// a factory func to generate a net connection.
factory := func() (net.Conn, error) { return net.Dial("tcp", "127.0.0.1:9999") }

// create a new channel based pool with an initial capacity of 5, maximum capacity of 30,
// and maximum actives of 15.
pool, err := pool.NewChannelPool(5, 30, 15, factory)

// get a connection from the pool, non blocking mode, if reach the limit, return nil instead. 

// Block Mode
conn, err := pool.Get()

// balabalabala
// close the conn, will release the conn to the pool by calling Close()
conn.Close()

// Non Blocking Mode, will return nil when no permission
conn, err := pool.TryGet()

// close the pool
pool.Close()

// return the connections length
p.Len()

// return the current active permissions
p.LenActives()
```
