package main

import (
	"bitbucket.org/nndi/phada"
	"flag"
	"fmt"
	"log"
	"net/http"
)

const (
	STATE_MENU          = 0
	STATE_DISK_SPACE    = 1
	STATE_MEMORY        = 2
	STATE_NETWORK       = 3
	STATE_TOP_PROCESSES = 4
	STATE_NOOP          = 5

	USSD_MENU = `USSD Server Stats Daemon
Host: %s
IP: %s

1. Disk Space
2. Memory
3. Network
4. Top Processes
#. Quit`

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

var (
	hostName    string
	bindAddress string
)

type UssdApp struct {
	sessionStore phada.SessionStore
	CurrentState int
}

func newUssdApp(sessionStore phada.SessionStore) *UssdApp {
	return &UssdApp{
		sessionStore: sessionStore,
		CurrentState: STATE_NOOP,
	}
}

func ussdContinue(text string) string {
	return fmt.Sprintf("CON %s", text)
}

func ussdEnd(text string) string {
	return fmt.Sprintf("END %s", text)
}

func (u *UssdApp) handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	session, err := phada.ParseUssdRequest(req)
	if err != nil {
		log.Fatalf("Failed to parse UssdRequest from http.Request. Error {}", err)
	}
	session.SetState(STATE_NOOP)

	u.sessionStore.PutHop(session)

	session, err = u.sessionStore.Get(session.SessionID)
	if err != nil {
		log.Fatalf("Failed to read session %s", err)
	}

	if session.ReadIn() == "" {
		session.SetState(STATE_MENU)
	}
	if session.ReadIn() == "1" {
		session.SetState(STATE_DISK_SPACE)
	}

	if session.ReadIn() == "2" {
		session.SetState(STATE_MEMORY)
	}

	if session.ReadIn() == "3" {
		session.SetState(STATE_NETWORK)
	}

	if session.ReadIn() == "4" {
		session.SetState(STATE_TOP_PROCESSES)
	}

	switch session.State {
	case STATE_MENU:
		text := fmt.Sprintf(USSD_MENU, hostName, bindAddress)
		fmt.Fprintf(w, ussdContinue(text))
		break
	case STATE_DISK_SPACE:
		fmt.Fprintf(w, ussdEnd(SAMPLE_DISK_STATS))
		break
	case STATE_MEMORY:
		fmt.Fprintf(w, ussdEnd(SAMPLE_MEM_STATS))
		break
	case STATE_NETWORK:
		fmt.Fprintf(w, ussdEnd(SAMPLE_NET_STATS))
		break
	case STATE_TOP_PROCESSES:
		fmt.Fprintf(w, ussdEnd(SAMPLE_PROC_STATS))
		break
	case STATE_NOOP:
	default:
		fmt.Fprintf(w, ussdEnd("# exit()"))
		break
	}
}

func init() {
	flag.StringVar(&hostName, "h", "example.com", "Hostname")
	flag.StringVar(&bindAddress, "b", "127.0.0.1:8000", "Bind address")
}

func main() {
	ussdApp := newUssdApp(phada.NewInMemorySessionStore())
	http.HandleFunc("/", ussdApp.handler)
	log.Fatalf("Failed to start server. Error %s", http.ListenAndServe(bindAddress, nil))
}
