// +build linux darwin netbsd freebsd openbsd dragonfly

package main

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry/gosigar"
	humanize "github.com/dustin/go-humanize"
)

func ReadDiskInfo(mountPoints []string) string {
	fslist := sigar.FileSystemList{}
	fslist.Get()

	var text strings.Builder
	text.WriteString("Disk\n\n")
	templ := "%s\n  %s/%s\n"
	for _, mnt := range mountPoints {
		for _, d := range fslist.List {
			if d.DirName == mnt {

				usage := sigar.FileSystemUsage{}
				usage.Get(d.DirName)

				text.WriteString(fmt.Sprintf(
					templ,
					d.DevName,
					humanize.Bytes(usage.Avail),
					humanize.Bytes(usage.Total), // sigar.FormatPercent(usage.UsePercent())
				))
				break
			}
		}
	}

	return text.String()
}

func ReadMemoryInfo() string {
	text := `Memory

mem
  Total: %s
  Free: %s
  Avail: %s
swap
  Total: %s
  Free: %s
`

	mem := sigar.Mem{}
	swap := sigar.Swap{}
	/*
	if err != nil {
		log.Printf("Failed to fetch memory stats. %s", err)
		return "Memory\nFailed to fetch data"
	}
	*/
	mem.Get()
	swap.Get()
	return fmt.Sprintf(text,
		humanize.Bytes(mem.Total*humanize.KByte),
		humanize.Bytes(mem.ActualUsed*humanize.KByte),
		humanize.Bytes(mem.ActualFree*humanize.KByte),
		humanize.Bytes(swap.Total*humanize.KByte),
		humanize.Bytes(swap.Free*humanize.KByte),
	)
}
