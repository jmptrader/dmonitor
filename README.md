## **dmonitor** : A simple daemon monitor web application written in Go

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