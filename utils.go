package dmonitor

import (
	"log"
	"errors"
	"io/ioutil"
	"encoding/json"

	"golang.org/x/crypto/ssh"
)

// Cache of SSH connections and sessions to various hosts for the same user
var connSessionCache map[string]SSHClientSession

// LoadConfig loads the configuration data from a JSON file.
// It returns a ControlPage struct with the data from the JSON.
func LoadConfig() (ControlPage, error) {
	cp := ControlPage{}

	data, err := ioutil.ReadFile("../config/config.json")
	if err != nil {
		log.Println("Error reading config file:", err)
		return cp, err
	}
	
	err = json.Unmarshal(data, &cp)
	if err != nil {
		return cp, err
	}

	return cp, nil
}

// ConnectToHost sets up a SSH connection to the host with username and password.
// Returns a client object and a session object.
func connectToHost(host, user, pass string) (*ssh.Client, *ssh.Session, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},
	}

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, session, nil
}

func GenControlPage() {}

func UpdateCurrentHostEnv(cp *ControlPage, hostValue, envValue string) {
	for i := range cp.Hosts {
		if cp.Hosts[i].Value == hostValue {
			cp.Hosts[0], cp.Hosts[i] = cp.Hosts[i], cp.Hosts[0]
		}
	}
	for i := range cp.Envs {
		if cp.Envs[i].Value == envValue {
			cp.Envs[0], cp.Envs[i] = cp.Envs[i], cp.Envs[0]
		}
	}
	cp.CurrentHost = cp.Hosts[0]
	cp.CurrentEnv = cp.Envs[0]
}

func UpdateDaemonsStatus(cp *ControlPage) {
	// check if there is a SSH session for the current host
	// if not set up SSH connection
	log.Println("Checking session for user:", cp.Username)
	if _, ok := connSessionCache[cp.CurrentHost.Value]; !ok {
		
	}

	for i := range cp.Daemons {
		cp.Daemons[i].Status = "Running"
		cp.Daemons[i].Control = "Stop"
	}
}

func initSessionCache(hosts []Host, user, pass string) error {
	// Set up a connection session for all the hosts.
	for i := range hosts {
		log.Printf("Connecting to host: %+v\n", hosts[i])
		client, session, err := connectToHost(hosts[i].Value, user, pass)
		if err != nil {
			return errors.New(string("failed to connect to host " + hosts[i].Name))
		}
		sc := SSHClientSession{client, session}
		connSessionCache[hosts[i].Value] = sc
		log.Println("Connected to host:", hosts[i].Value)
	}
	return nil
}

func LoginUser(cp *ControlPage, user, pass string) error {
	cp.Username = user
	cp.CurrentHost = cp.Hosts[0]
	cp.CurrentEnv = cp.Envs[0]

	log.Println("Setting up SSH connection...")
	if err := initSessionCache(cp.Hosts, user, pass); err != nil {
		cp.Username = ""
		return err
	}

	return nil
}

func LogoutUser(cp *ControlPage) {
	// Remove user
	cp.Username = ""
	// Empty daemon details
	for i := range cp.Daemons {
		cp.Daemons[i].Status = ""
		cp.Daemons[i].Control = ""
	}
	// Close all connections and sessions
	for _, v := range connSessionCache {
		v.Session.Close()
		v.Client.Close()
	}
	// Empty SSH connection cache
	connSessionCache = map[string]SSHClientSession{}
}

func runCmd(cp *ControlPage, cmd string) bool {

	return true
}

func StartOrStopDaemon(cp *ControlPage, daemonName, control string) error {
	var cmd string
	var di int
	for i := range cp.Daemons {
		if cp.Daemons[i].Name == daemonName {
			if control == "Stop" {
				cmd = cp.Daemons[i].StopCmd
			} else if control == "Start" {
				cmd = cp.Daemons[i].StartCmd
			}
			di = i
		}
	}

	if cmd == "" {
		return errors.New(string("unable to find command to " + control + " " + daemonName))
	}

	runCmd(cp, cmd)
	log.Println(string("Executed command: " + cp.Daemons[di].StartCmd))
	log.Println(string("Change control on daemon: " + daemonName))
	return nil
}