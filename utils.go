package dmonitor

import (
	"log"
	"errors"
	"io/ioutil"
	"encoding/json"

	"golang.org/x/crypto/ssh"
)

// Cache of SSH clients to various hosts for the same user.
var clientCache map[string]*ssh.Client

// LoadConfig loads the configuration data from a JSON file.
// It returns a ControlPage struct with the data from the JSON.
func LoadConfig() (ControlPage, error) {
	cp := ControlPage{}

	data, err := ioutil.ReadFile("config/config.json")
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
// Returns a client object.
func connectToHost(host, user, pass string) (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},
	}

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}

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

func ReloadDaemonsStatus(cp *ControlPage) {
	log.Println("Checking session for user:", cp.Username)
	if client, ok := clientCache[cp.CurrentHost.Value]; ok {
		
		// Check status of all daemons and update the control page
		for i := range cp.Daemons {
			daemon := &cp.Daemons[i]
			session, _ := client.NewSession()
			session.Stdout = nil
			session.Stderr = nil
			log.Println("Running command: ", daemon.StatusCmd)
			out, err := session.CombinedOutput(daemon.StatusCmd)
			if err != nil && string(out) != "" {
				log.Printf("Failed to update daemon status for daemon: %+v", daemon)
				log.Println("Error:", err)
			}
			if string(out) == "" {
				daemon.Status, daemon.Control = "Stopped", "Start"
			} else {
				daemon.Status, daemon.Control = "Running", "Stop"
			}
			log.Printf("Daemon: %s, Status: %s\n", daemon.Name, daemon.Status)
			session.Close()
		}
	}
}

func initSessionCache(hosts []Host, user, pass string) error {
	clientCache = map[string]*ssh.Client{}
	for i := range hosts {
		log.Printf("Connecting to host: %+v\n", hosts[i])
		client, err := connectToHost(hosts[i].Value, user, pass)
		if err != nil {
			return errors.New(string("failed to connect to host " + hosts[i].Name))
		}
		
		clientCache[hosts[i].Value] = client
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
	for _, v := range clientCache {
		v.Close()
	}
	// Empty SSH connection cache
	clientCache = map[string]*ssh.Client{}
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

	client := clientCache[cp.CurrentHost.Value]
	session, _ := client.NewSession()
	session.Stdout = nil
	session.Stderr = nil
	defer session.Close()
	log.Printf("Running command: %q", cmd)
	out, err := session.CombinedOutput(cmd)
	if err != nil {
		log.Println("Command output:", string(out))
		s, _ := client.NewSession()
		s.Stdout = nil
		s.Stderr = nil
		out1, err1 := s.Output("echo $PATH")
		if err1 != nil {
			log.Println("echo $PATH error:", err1)
		}
		log.Println("PATH:", string(out1))
		s.Close()
		return err
	}
	if control == "Start" {
		cp.Daemons[di].Status = "Running"
		cp.Daemons[di].Control = "Stop"
	} else if control == "Stop" {
		cp.Daemons[di].Status = "Stopped"
		cp.Daemons[di].Control = "Start"
	}

	log.Println(string("Executed command: " + cmd))
	log.Println(string("Change control on daemon: " + daemonName))
	return nil
}