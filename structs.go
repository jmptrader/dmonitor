package dmonitor

import (
	"golang.org/x/crypto/ssh"
)

type Host struct {
	Name string `json:"name"`
	Value string `json:"value"`
}

type Env struct {
	Name string `json:"name"`
	Value string `json:"value"`
}

type SSHClientSession struct {
	Client *ssh.Client
	Session *ssh.Session
}
	
type Daemon struct {
	Name string `json:"name"`
	Status string `json:"-"`	// dynamic value to be updated at runtime
	Control string `json:"-"`	// dynamic value to be updated at runtime
	StopCmd string `json:"stopcmd"`
	StartCmd string `json:"startcmd"`
}

type ControlPage struct {
	Username string `json:"-"`
	Hosts []Host `json:"hosts"`
	Envs []Env `json:"envs"`
	Daemons []Daemon `json:"daemons"`
	CurrentHost Host
	CurrentEnv Env
}
