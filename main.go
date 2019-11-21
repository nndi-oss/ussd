package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"bitbucket.org/nndi/phada"
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
)

var (
	hostName      string
	bindAddress   string
	isDummyServer bool
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
		log.Printf("Failed to parse UssdRequest from http.Request. Error %s", err)
		fmt.Fprintf(w, ussdEnd("Failed to process request"))
		return
	}
	session.SetState(STATE_NOOP)

	u.sessionStore.PutHop(session)

	session, err = u.sessionStore.Get(session.SessionID)
	if err != nil {
		log.Printf("Failed to read session %s", err)
		fmt.Fprintf(w, ussdEnd("Failed to process request"))
		return
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
		if isDummyServer {
			fmt.Fprintf(w, ussdEnd(SAMPLE_DISK_STATS))
		} else {
			fmt.Fprintf(w, ussdEnd(ReadDiskInfo(
				[]string{"/", "/home", "/var"},
			)))
		}
		break
	case STATE_MEMORY:
		if isDummyServer {
			fmt.Fprintf(w, ussdEnd(SAMPLE_MEM_STATS))
		} else {
			fmt.Fprintf(w, ussdEnd(ReadMemoryInfo()))
		}

		break
	case STATE_NETWORK:
		if isDummyServer {
			fmt.Fprintf(w, ussdEnd(SAMPLE_NET_STATS))
		} else {
			fmt.Fprintf(w, ussdEnd(SAMPLE_NET_STATS))
		}
		break
	case STATE_TOP_PROCESSES:
		if isDummyServer {
			fmt.Fprintf(w, ussdEnd(SAMPLE_PROC_STATS))
		} else {
			fmt.Fprintf(w, ussdEnd(SAMPLE_PROC_STATS))
		}
		break
	case STATE_NOOP:
	default:
		fmt.Fprintf(w, ussdEnd("# exit()"))
		break
	}
}

func init() {
	flag.BoolVar(&isDummyServer, "dummy", false, "Start the dummy server - uses hardcoded values")
	flag.StringVar(&hostName, "h", "example.com", "Hostname")
	flag.StringVar(&bindAddress, "b", "127.0.0.1:8000", "Bind address")
}

func main() {
	ussdApp := newUssdApp(phada.NewInMemorySessionStore())
	http.HandleFunc("/", ussdApp.handler)
	log.Fatalf("Failed to start server. Error %s", http.ListenAndServe(bindAddress, nil))
}
