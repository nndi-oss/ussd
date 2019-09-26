USSD Server Stats Daemon (USSD)
===

At [NNDI](https://nndi-tech.com) we believe that USSD is a great protocol for the
African market and we like to play around with such technologies
to see how far we can stretch them for innovative applications. 

This project implements a **proof of concept** USSD application that allows system administrators
to check some basic stats about their servers via their mobile phone. It's a 
proof of concept of an idea we're calling the `USSD of Things` ;)

> NOTE: It's still a work in progress and not an offical NNDI product

## Usage

Run it with the following command

```sh
$ go get -u "github.com/nndi-oss/ussd"

$ cd $GOPATH/src/github.com/nndi-oss/ussd

$ go run main.go -h "my.server.com" -bind "localhost:8000" -sample
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

---

Copyright (c) 2018, NNDI