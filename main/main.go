package main

import (
	"log"
	"html/template"
	"net/http"

	"github.com/mubitosh/dmonitor"
)

var templates = template.Must(template.ParseFiles("../html/login.html", "../html/monitor.html"))

//var client *ssh.Client
//var session *ssh.Session

var cp dmonitor.ControlPage

// renderTemplate renders a template with the data in ControlPage struct p.
func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", cp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login")
}

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user := r.Form["username"][0]
	pass := r.Form["password"][0]

	log.Println("Logging in with user:", user)
	err := dmonitor.LoginUser(&cp, user, pass)
	if err != nil {
		log.Println("Login failed:", err)
		log.Println("Redirecting to login page")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	log.Println("Logged in as user:", user)
	log.Println("Redirecting to monitor page")
	http.Redirect(w, r, "/monitor", http.StatusFound)
}

func logout(w http.ResponseWriter, r *http.Request) {
	user := cp.Username
	log.Println("Logging out user:", user)
	dmonitor.LogoutUser(&cp)
	log.Println("Logged out user:", user)
	log.Println("Redirecting to login page")
	http.Redirect(w, r, "/", http.StatusFound)
}

func monitor(w http.ResponseWriter, r *http.Request) {
	if cp.Username == "" {
		log.Println("No user found, redirecting to login page")
		http.Redirect(w, r, "/", http.StatusFound)
        return
	}

	log.Println("Updating status of all daemons")
	log.Printf("Host: %+v\n", cp.CurrentHost)
	log.Printf("Environment: %+v\n", cp.CurrentEnv)
	log.Printf("\n\nControlPage:\n%+v\n\n", cp)
	dmonitor.UpdateDaemonsStatus(&cp)
	log.Printf("\n\nControlPage:\n%+v\n\n", cp)
	log.Println("Update done")
	renderTemplate(w, "monitor")
}

func reloadlist(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	hostValue := r.Form["hostname"][0]
	envValue := r.Form["environment"][0]
	dmonitor.UpdateCurrentHostEnv(&cp, hostValue, envValue)
	log.Println("Reloading daemons status for host:", hostValue, "and env:", envValue)
	http.Redirect(w, r, "/monitor", http.StatusFound)
}

func startOrStop(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Form == nil {
		log.Println("\nEmpty control received\n")
	}
	for daemonName, control := range r.Form {
		log.Println("Request to", control, daemonName)
		if err := dmonitor.StartOrStopDaemon(&cp, daemonName, control[0]); err != nil {
			log.Println("Failed to execute command:", err)
		}
		break;
	}
	http.Redirect(w, r, "/monitor", http.StatusFound)
}

func main() {
	var err error
	cp, err = dmonitor.LoadConfig()

	if err != nil {
		log.Println("Cannot load config file. Exiting application.\n", err)
		return
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/monitor", monitor)
	http.HandleFunc("/reloadlist", reloadlist)
	http.HandleFunc("/startOrStop", startOrStop)
	log.Println("Starting dmonitor at port 8008")
	http.ListenAndServe(":8008", nil)
}