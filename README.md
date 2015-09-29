## **dmonitor** : A simple daemon monitor web application written in Go

dmonitor is a simple web application for monitoring daemons. Provided for a daemon commands to start it, to stop it, to check its status, it can be monitored with a web interface.
The daemons are expected to be running on a remote host which can be connected to with SSH.

### Project layout

Below is a screenshot of the project layout.

![project layout](https://github.com/mubitosh/dmonitor/blob/master/main/images/project-layout.jpeg "project layout")


### Build it
The following commands will build the project.

1. Build the dmonitor package

```bash
$ cd $GOPATH/src
$ go build github.com/mubitosh/dmonitor
```

2. Build the main package

```bash
$ cd $GOPATH/src/github.com/mubitosh/dmonitor/main
$ go build main.go
```

The above command creates an executable ```main``` in the directory ```$GOPATH/src/github.com/mubitosh/dmonitor/main```. The application is expected to be located at ```$GOPATH/src/github.com/mubitosh/dmonitor/main```.

### Configure it
The configuration file ```config.json``` at ```dmonitor/main/config/``` directory can be modified as needed.

A sample content of ```config.json``` is below.

```hosts``` is the list of hosts. For all the hosts listed a common username password pair with SSH access should be there. Not necessary though and the source code can be modified as needed.

```envs``` is the list of envs. This is an option in case multiple copies of the same daemon run in the same host in different environments. The start/stop command for a daemon in different environments maybe different depending on the requirement.

```daemons``` is the list of daemons. The commnands to start/stop a daemon and command to check the status of the daemon are given in this list. The example setup below assumes a very simple daemon. If the commands are very complex, it would be better to replace it with a shell script and execute the shell script instead.

```json
{
    "hosts": [
        {
            "name": "host01",
            "value": "127.0.0.1:22"
        }
    ],

    "envs": [
        {
            "name": "env01",
            "value": "env01"
        },
        {
            "name": "env02",
            "value": "env02"
        },
        {
            "name": "env03",
            "value": "env03"
        }
    ],
    
    "daemons": [
        {
            "name": "daemon01",
            "startcmd": "daemon01",
            "stopcmd": "pkill daemon01",
            "statuscmd": "pidof daemon01"
        },
        {
            "name": "daemon02",
            "startcmd": "daemon02",
            "stopcmd": "pkill daemon02",
            "statuscmd": "pidof daemon02"
        },
        {
            "name": "daemon03",
            "startcmd": "daemon03",
            "stopcmd": "pkill daemon03",
            "statuscmd": "pidof daemon03"
        }
    ]
}
```

### Run it
The following commands will run the application.

```bash
$ cd $GOPATH/src/mubitosh/dmonitor/main
$ ./main
```

###	Use it
Open a web browser and go the url ```localhost:8008/```(it is the default port number). It shows a login page.

![dmonitor login page](https://github.com/mubitosh/dmonitor/blob/master/main/images/dmonitor-login-page-screenshot.jpeg "dmonitor login page")

Provide a username and password for SSH connection to the hosts provided in ```config.json```. After a successful login, the monitor page will show the current status of the daemons listed in config.json. 

![dmonitor monitor page](https://github.com/mubitosh/dmonitor/blob/master/main/images/dmonitor-monitor-page-screenshot.jpeg "dmonitor monitor page")
Start/Stop a daemon, refresh the list as needed.

### How does it work?
A web server running locally maintains a cache of SSH clients. Whenever a command is to be executed on one of the clients, it creates a new session with the client. It runs the command and according to the result the webpage is updated.

![dmonitor architecture](https://github.com/mubitosh/dmonitor/blob/master/main/images/dmonitor-architecture.jpeg "dmonitor architecture")


### Notes

**[About creating a daemon in Linux](http://www.netzmafia.de/skripten/unix/linux-daemon-howto.html)**

May be you need a simple daemon to test out the application. The project includes a simple daemon written in C for Linux. It is aptly named ```noob_daemon.c```(located under ```dmonitor/config/```). Compile it and place in the $PATH which will be available in a SSH session.

**[Swirl pattern background and other backgrounds at http://subtlepatterns.com/](http://subtlepatterns.com/)**


### TODO

**Add support for session handling**

Currently the application does not have support for session handling. So if you are logged in from a web browser. You can open a different web browser, open the monitor page and will be able to access it.