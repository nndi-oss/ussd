package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"bitbucket.org/nndi/phada"
	log "github.com/inconshreveable/log15"
	"github.com/pelletier/go-toml"
)

const (
	STATE_MENU          = 0
	STATE_PROCESS_MENU  = 1010
	STATE_DISK_SPACE    = 1
	STATE_MEMORY        = 2
	STATE_NETWORK       = 3
	STATE_TOP_PROCESSES = 4
	STATE_SERVICES_MENU = 5

	STATE_PROMPT_CHECK_SERVICE   = 51
	STATE_PROMPT_START_SERVICE   = 52
	STATE_PROMPT_STOP_SERVICE    = 53
	STATE_PROMPT_RESTART_SERVICE = 54
	STATE_PROMPT_ENABLE_SERVICE  = 55
	STATE_PROMPT_DISABLE_SERVICE = 56

	STATE_CHECK_SERVICE   = 511
	STATE_START_SERVICE   = 521
	STATE_STOP_SERVICE    = 531
	STATE_RESTART_SERVICE = 541
	STATE_ENABLE_SERVICE  = 551
	STATE_DISABLE_SERVICE = 561

	STATE_EXIT = -1

	USSD_MENU = `USSD Sysadmin
Host: %s

1 Disk Space
2 Memory
3 Network
4 Top Processes
5 Services`

	SERVICES_MENU = `Services
1 Check Status
2 Start Service
3 Stop Service
4 Restart Service
5 Enable Service
6 Disable Service
# Main Menu`
)

var (
	configFile    string
	isDummyServer bool
	hostName      string
)

type UssdApp struct {
	host            string
	sessionStore    phada.SessionStore
	servicectl      ServiceManager
	allowedServices []string
}

// newUssdApp creates a new instance of the USSD handler
func newUssdApp(hostName string, services ServiceManager, sessionStore phada.SessionStore) *UssdApp {
	return &UssdApp{
		host:         hostName,
		sessionStore: sessionStore,
		servicectl:   services,
	}
}

func ussdContinue(text string) string {
	return fmt.Sprintf("CON %s", text)
}

func ussdEnd(text string) string {
	return fmt.Sprintf("END %s", text)
}

func (u *UssdApp) SaveSession(session *phada.UssdRequestSession, nextState int) {
	session.SetState(nextState)
	u.sessionStore.PutHop(session)
}

func (u *UssdApp) handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	newSession, err := phada.ParseUssdRequest(req)
	if err != nil {
		log.Error("Failed to parse UssdRequest from http.Request.", "err", err)
		fmt.Fprintf(w, ussdEnd("Failed to process request"))
		return
	}
	newSession.SetState(STATE_MENU)
	u.sessionStore.PutHop(newSession)

	session, err := u.sessionStore.Get(newSession.SessionID)
	if err != nil {
		log.Info("Did not find existing session, creating new session", "error", err)
		session = newSession
	}

	sessionLogger := log.New("sessionID", session.SessionID, "phoneNumber", session.PhoneNumber)

	if session.ReadIn() == "" {
		session.SetState(STATE_MENU)
	}

	if session.State == STATE_PROCESS_MENU {
		if session.ReadIn() == "1" {
			u.SaveSession(session, STATE_DISK_SPACE)
		}

		if session.ReadIn() == "2" {
			u.SaveSession(session, STATE_MEMORY)
		}

		if session.ReadIn() == "3" {
			u.SaveSession(session, STATE_NETWORK)
		}

		if session.ReadIn() == "4" {
			u.SaveSession(session, STATE_TOP_PROCESSES)
		}

		if session.ReadIn() == "5" {
			u.SaveSession(session, STATE_SERVICES_MENU)
			fmt.Fprintf(w, ussdContinue(SERVICES_MENU))
			return
		}
	}

	switch session.State {
	case STATE_MENU:
		text := fmt.Sprintf(USSD_MENU, u.host)
		u.SaveSession(session, STATE_PROCESS_MENU)
		fmt.Fprintf(w, ussdContinue(text))
		break

	case STATE_DISK_SPACE:
		sessionLogger.Info("Checking disk usage")
		if isDummyServer {
			fmt.Fprintf(w, ussdEnd(SAMPLE_DISK_STATS))
		} else {
			fmt.Fprintf(w, ussdEnd(ReadDiskInfo(
				[]string{"/", "/home", "/var"},
			)))
		}
		break
	case STATE_MEMORY:
		sessionLogger.Info("Checking memory usage")
		if isDummyServer {
			fmt.Fprintf(w, ussdEnd(SAMPLE_MEM_STATS))
		} else {
			fmt.Fprintf(w, ussdEnd(ReadMemoryInfo()))
		}

		break
	case STATE_NETWORK:
		sessionLogger.Info("Checking network stats")
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
	case STATE_SERVICES_MENU:
		serviceMenuItem := session.ReadIn()

		if serviceMenuItem == "1" {
			u.SaveSession(session, STATE_CHECK_SERVICE)
			fmt.Fprintf(w, ussdContinue("Check Service\nEnter service name: "))
			break
		}

		if serviceMenuItem == "2" {
			u.SaveSession(session, STATE_START_SERVICE)
			fmt.Fprintf(w, ussdContinue("Start Service\nEnter service name: "))
			break
		}

		if serviceMenuItem == "3" {
			u.SaveSession(session, STATE_STOP_SERVICE)
			fmt.Fprintf(w, ussdContinue("Stop Service\nEnter service name: "))
			break
		}

		if serviceMenuItem == "4" {
			u.SaveSession(session, STATE_RESTART_SERVICE)
			fmt.Fprintf(w, ussdContinue("Restart Service\nEnter service name: "))
			break
		}

		if serviceMenuItem == "5" {
			u.SaveSession(session, STATE_ENABLE_SERVICE)
			fmt.Fprintf(w, ussdContinue("Enable Service\nEnter service name: "))
			break
		}

		if serviceMenuItem == "6" {
			u.SaveSession(session, STATE_DISABLE_SERVICE)
			fmt.Fprintf(w, ussdContinue("Disable Service\nEnter service name: "))
			break
		}

		u.SaveSession(session, STATE_SERVICES_MENU)
		break

	case STATE_CHECK_SERVICE:
		serviceName := session.ReadIn()
		sessionLogger.Info("Checking service status", "service", serviceName)
		if output, err := u.servicectl.Status(serviceName); err == nil {
			// TODO: Check for Status active...
			fmt.Fprintf(w, ussdEnd("$ systemctl is-active %s\nResult: Service is %s"), serviceName, output)
		}
		break
	case STATE_START_SERVICE:
		serviceName := session.ReadIn()
		sessionLogger.Info("Attempting to start service", "service", serviceName)
		if _, err := u.servicectl.Start(serviceName); err == nil {
			fmt.Fprintf(w, ussdEnd("$ systemctl start %s\nResult: Service started successfully."), serviceName)
		} else {
			fmt.Fprintf(w, ussdEnd("FAILED to start service %s."), serviceName)
		}
		break
	case STATE_STOP_SERVICE:
		serviceName := session.ReadIn()
		sessionLogger.Info("Attempting to stop service", "service", serviceName)
		if _, err := u.servicectl.Stop(serviceName); err == nil {
			fmt.Fprintf(w, ussdEnd("$ systemctl stop %s\nResult: Service stopped successfully."), serviceName)
		} else {
			fmt.Fprintf(w, ussdEnd("FAILED to stop service %s."), serviceName)
		}
		break
	case STATE_RESTART_SERVICE:
		serviceName := session.ReadIn()
		sessionLogger.Info("Attempting to restart service", "service", serviceName)
		if _, err := u.servicectl.Restart(serviceName); err == nil {
			fmt.Fprintf(w, ussdEnd("$ systemctl restart %s\nResult: Service restarted successfully."), serviceName)
		} else {
			fmt.Fprintf(w, ussdEnd("FAILED to restart service %s."), serviceName)
		}
		break
	case STATE_ENABLE_SERVICE:
		serviceName := session.ReadIn()
		sessionLogger.Info("Attempting to enable service", "service", serviceName)
		if _, err := u.servicectl.Enable(serviceName); err == nil {
			fmt.Fprintf(w, ussdEnd("$ systemctl enable %s\nResult: Service enabled successfully."), serviceName)
		} else {
			fmt.Fprintf(w, ussdEnd("FAILED to enable service %s."), serviceName)
		}
		break
	case STATE_DISABLE_SERVICE:
		serviceName := session.ReadIn()
		sessionLogger.Info("Attempting to disable service", "service", serviceName)
		if _, err := u.servicectl.Disable(serviceName); err == nil {
			fmt.Fprintf(w, ussdEnd("$ systemctl disable %s\nResult: Service disabled successfully."), serviceName)
		} else {
			fmt.Fprintf(w, ussdEnd("FAILED to disable service %s."), serviceName)
		}
		break

	case STATE_EXIT:
	default:
		u.SaveSession(session, STATE_EXIT)
		fmt.Fprintf(w, ussdEnd("Have a nice day :)"))
		break
	}
	u.sessionStore.PutHop(session)
}

func init() {
	flag.BoolVar(&isDummyServer, "dummy", false, "Start the dummy server - uses hardcoded values")
	flag.StringVar(&configFile, "config", "ussd.toml", "Configuration file, defaults to 'ussd.toml'")
}

func main() {
	flag.Parse()

	tree, err := toml.LoadFile(configFile)
	if err != nil {
		fmt.Printf("Failed to read configuration file %v", err)
		os.Exit(-1)
		return
	}

	displayName := tree.Get("server_name").(string)
	host := tree.Get("host").(string)
	port := tree.Get("port").(int64)
	bindAddress := fmt.Sprintf("%s:%d", host, port)

	ussdApp := newUssdApp(displayName, &Systemctl{"/bin/"}, phada.NewRistrettoSessionStore())

	http.HandleFunc("/", ussdApp.handler)
	log.Error("Starting server", "err", http.ListenAndServe(bindAddress, nil))
}
