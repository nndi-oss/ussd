USSD Server Stats Daemon (USSD)
===

This project implements a **proof of concept** USSD application that allows 
System Administrators to check/monitor basic stats of their servers via their mobile phone.

## Background

At [NNDI](https://nndi-tech.com) we like to play around with different technologies to 
see how far we can stretch them for innovative applications. USSD is an old-but-great protocol 
that's popular and serves a lot of applications in Africa including Mobile Money and Banking.

We decided to try using USSD for something atypical - a kind of IOT application;
perhaps you can call this a `USSD of Things` project. ;)

## Usage

This project is intended to be used with USSD APIs provided by [AfricasTalking](https://africastalking.com) - so you will need to get access to their services to run it 
behind an actual USSD Shortcode. However, you can run it locally and test it with
[dialoguss](https://github.com/nndi-oss/dialoguss)

Run it with the following command

```sh
$ go get -u "github.com/nndi-oss/ussd"

$ cd $GOPATH/bin

$ ./ussd -h "my.server.com" -bind "localhost:8000"
```

### Dummy Server

The dummy server provides hard-coded results for all calls and doesn't actually
read stats from the host machine. You can use this to test the server

```sh
$ go run main.go -h "my.server.com" -bind "localhost:8000" -dummy
```

## Basic USSD interaction Concept

1. Dial Short code (e.g. `*384*8327#`)
```
USSD Server Stats Daemon

Host: my.server.com
IP: 192.168.1.1

1. Disk space
2. Memory
3. Network 
4. Processes
#  Quit
```

Input: 1

```
Disk Space on my.server.com

/dev/sda4  
  5GB/50GB (90%)
/dev/sda4
  25GB/50GB (25%)
```

Input: 2

```
Memory
mem
  Total: 6GB
  Used: 4.3GB
  Rsrv: 1.7GB
swap
  Total: 4GB
  Used: 0GB
  Rsrv: 0GB
```
  
Input: 3

```
Network

epn03
  up:  yes
  in:  30GB
  out: 60GB
eth0
  up:  yes
  in:  30GB
  out: 60GB
```

Input: 4

```
Processes

top:
1. java (0.2 cpu, 2GB mem)
2. systemd (0.1 cpu, 500MB mem)
3. http (0.1 cpu, 50MB mem)
```

> NOTE: It's still a work in progress and not (yet) an offical NNDI product

---

Copyright (c) 2018 - 2019, NNDI