package main

const (
	SAMPLE_DISK_STATS = `Disk Space on my.server.com

/dev/sda4  
  5GB/50GB (90%%)
/dev/sda4
  25GB/50GB (25%%)`

	SAMPLE_MEM_STATS = `Memory
mem
  Total: 6GB
  Used: 4.3GB
  Rsrv: 1.7GB
swap
  Total: 4GB
  Used: 0GB
  Rsrv: 0GB`

	SAMPLE_NET_STATS = `Network

epn03
  up:  yes
  in:  30GB
  out: 60GB
eth0
  up:  yes
  in:  30GB
  out: 60GB`

	SAMPLE_PROC_STATS = `Processes

top:
1. java (0.2 cpu, 2GB mem)
2. systemd (0.1 cpu, 500MB mem)
3. http (0.1 cpu, 50MB mem)`
)