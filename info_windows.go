// +build windows

package main

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry/gosigar/sys/windows"
	humanize "github.com/dustin/go-humanize"
)

func ReadDiskInfo(diskLabels []string) string {
	var text strings.Builder
	text.WriteString("Disk\n\n")
	templ := "%s\n  %s/%s\n"
	for _, mnt := range diskLabels {
		freeBytesAvailable, totalNumberOfBytes, _, err := windows.GetDiskFreeSpaceEx(mnt)

		if err != nil {
			return "Disk\nFailed to fetch data"
		}

		text.WriteString(fmt.Sprintf(
			templ,
			mnt,
			humanize.Bytes(freeBytesAvailable),
			humanize.Bytes(totalNumberOfBytes),
		))
	}

	return text.String()
}

func ReadMemoryInfo() string {
	mem, err := windows.GlobalMemoryStatusEx()
	if err != nil {
		return "Memory\nFailed to fetch data"
	}
	text := `Memory

mem
  Total: %s
  Avail: %s
page
  Total: %s
  Free: %s
`
	return fmt.Sprintf(text,
		humanize.Bytes(mem.TotalPhys*humanize.KByte),
		humanize.Bytes(mem.AvailPhys*humanize.KByte),
		humanize.Bytes(mem.TotalPageFile*humanize.KByte),
		humanize.Bytes(mem.AvailPageFile*humanize.KByte),
	)
}
