## **Dmonitor** : A simple daemon monitor web application written in Go.

### Build it
The following commands will build the project-

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

The application is expected to locate at ```$GOPATH/src/github.com/mubitosh/dmonitor/main```.

### Configure it
The configuration file ```config.json``` at ```dmonitor/main/config/``` directory modifieds as needed.

### Run it
The following commands will run the application.

```bash
$ cd $GOPATH/src/mubitosh/dmonitor/main
$ ./main
```

###	Use it
Open a web browser and go the url ```localhost:8008/```(it is the default one). It shows a login page.

![dmonitor login page](https://github.com/mubitosh/dmonitor/blob/master/main/images/dmonitor-login-page-screenshot.png "dmonitor login page")

Provide a username and password for SSH connection to the hosts provided in ```config.json```. After a successful login, the monitor page will show the current status of the daemons listed in config.json. 

![dmonitor monitor page](https://github.com/mubitosh/dmonitor/blob/master/main/images/dmonitor-monitor-page-screenshot.png "dmonitor monitor page")
Start/Stop a daemon, refresh the list as needed.