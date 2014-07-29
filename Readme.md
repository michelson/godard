
# Godard

Godard is a simple process monitoring tool written in Go , and is a port from Bluepill library written in Ruby.

![Alt text](./assets/gopher-godard.gif)

[![Build Status](https://travis-ci.org/michelson/godard.png)](https://travis-ci.org/michelson/godard)

[![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/michelson/godard).



## Installation
Clone this repo and build dependences and run.

    make deps

    make run


In order to take advantage of logging with syslog, you also need to setup your syslog to log the local6 facility. Edit the appropriate config file for your syslogger (/etc/syslog.conf for syslog) and add a line for local6:

    local6.*          /var/log/godard.log

You&apos;ll also want to add _/var/log/godard.log_ to _/etc/logrotate.d/syslog_ so that it gets rotated.

Lastly, create the _/var/run/godard directory for godard to store its pid and sock files.

## Usage
### Config
Godard organizes processes into 3 levels: application -> group -> process. Each process has a few attributes that tell godard how to start, stop, and restart it, where to look or put the pid file, what process conditions to monitor and the options for each of those.

The minimum config file looks something like this:

```json
{
  "processes":
  [ {
      "name": "process_name",
      "start_command": "/usr/bin/some_start_command",
      "pid_file": "/tmp/some_pid_file.pid"
    }
  ]
}
```

Note that since we specified a PID file and start command, godard assumes the process will daemonize itself.

Unlike Bluepill, Godard will not daemonize processes. so the option ```daemonize: true``` is not, yet, available.


If you don&apos;t specify a stop command, a TERM signal will be sent by default. Similarly, the default restart action is to issue stop and then start.

Now if we want to do something more meaningful, like actually monitor the process, we do:

```json
{
  "processes":
  [ {
      "name": "process_name",
      "start_command": "/usr/bin/some_start_command",
      "pid_file": "/tmp/some_pid_file.pid",
      "checks": {
        "cpu_usage":{ "every": "10.seconds", "below": 5, "times": 3},
      }
    }
  ]
}
```

We added a line that checks every 10 seconds to make sure the cpu usage of this process is below 5 percent; 3 failed checks results in a restart. We can specify a two-element array for the _times_ option to say that it 3 out of 5 failed attempts results in a restart.

To watch memory usage, we just add one more line:

```json

{
  "processes":
  [ {
      "name": "process_name",
      "start_command": "/usr/bin/some_start_command",
      "pid_file": "/tmp/some_pid_file.pid",
      "checks": {
        "cpu_usage":{ "every": "10.seconds", "below": 5, "times": 3},
        "mem_usage":{ "every": "10.secs", "below": "100.megabytes", "times": 3}
      }
    }
  ]
}
```

To watch the modification time of a file, e.g. a log file to ensure the process is actually working add one more line:

```json
{
  "processes":
  [ {
      "name": "process_name",
      "start_command": "/usr/bin/some_start_command",
      "pid_file": "/tmp/some_pid_file.pid",
      "checks": {
        "cpu_usage":{ "every": "10.seconds", "below": 5, "times": 3},
        "mem_usage":{ "every": "10.secs", "below": "100.megabytes", "times": 3},
        "file_time":{ "every": "60.secs", "below": "3.minutes", "times": 3, "filename": "/tmp/some_file.log" }
      }
    }
  ]
}
```

To restart process if it's running too long:

```json
{
  "processes":
  [ {
      "name": "process_name",
      "start_command": "/usr/bin/some_start_command",
      "pid_file": "/tmp/some_pid_file.pid",
      "checks": {
        "running_time":{ "every": "60.secs", "below": "24.hours"}
      }
    }
  ]
}
```



We can tell godard to give a process some grace time to start/stop/restart before resuming monitoring:

```json
{
  "processes":
  [ {
      "name": "process_name",
      "start_command": "/usr/bin/some_start_command",
      "pid_file": "/tmp/some_pid_file.pid",
      "start_grace_time": "3.seconds",
      "stop_grace_time": "5.seconds",
      "restart_grace_time": "8.seconds",
      "checks": {
        "cpu_usage":{ "every": "10.seconds", "below": 5, "times": 3},
        "mem_usage":{ "every": "10.secs", "below": "100.megabytes", "times": 3}
      }
    }
  ]
}
```

We can group processes by name:

```json
{
  "processes":
  [ {
      "name": "process_name_1",
      "start_command": "/usr/bin/some_start_command",
      "pid_file": "/tmp/some_pid_file.pid",
      "group": "mongrels"
    },
    {
      "name": "process_name_2",
      "start_command": "/usr/bin/some_start_command",
      "pid_file": "/tmp/some_pid_file.pid",
      "group": "mongrels"
    }
  ]
}
```

If you want to run the process as someone other than root:

```json
{
  "processes":
  [ {
      "name": "process_name",
      "start_command": "/usr/bin/some_start_command",
      "pid_file": "/tmp/some_pid_file.pid",
      "uid": "deploy",
      "gid": "deploy",
      "checks": {
        "cpu_usage":{ "every": "10.seconds", "below": 5, "times": 3},
        "mem_usage":{ "every": "10.secs", "below": "100.megabytes", "times": 3}
      }
    }
  ]
}
```

If you want to include one or more supplementary groups:

```json
{
  "processes":
  [ {
      "name": "process_name",
      "start_command": "/usr/bin/some_start_command",
      "pid_file": "/tmp/some_pid_file.pid",
      "uid": "deploy",
      "gid": "deploy",
      "supplementary_groups": ["rvm"]
    }
  ]
}
```

You can also set an app-wide uid/gid:

```json
{
  "uid": "deploy",
  "gid": "deploy",
  "processes":
  [ {
      "name": "process_name",
      "start_command": "/usr/bin/some_start_command",
      "pid_file": "/tmp/some_pid_file.pid",
      "supplementary_groups": ["rvm"]
    }
  ]
}
```

To track resources of child processes, use :include_children:
```json
{
  "processes":
  [ {
      "name": "process_name",
      "start_command": "/usr/bin/some_start_command",
      "pid_file": "/tmp/some_pid_file.pid",
      "checks": {
        "mem_usage":{ "every": "10.secs", "below": "100.megabytes", "times": 3, "include_children": true }
      }
    }
  ]
}
```

To check for flapping:

```json
  "flapping":{ "times": 2, "within": "30.seconds", "retry_in": "7.seconds"}
```

To set the working directory to _cd_ into when starting the command:

```json
{
  "processes":
  [ {
      "name": "process_name",
      "start_command": "/usr/bin/some_start_command",
      "pid_file": "/tmp/some_pid_file.pid",
      "working_dir": "/path/to/some_directory"
    }
  ]
}
```

You can also have an app-wide working directory:

```json
{
  "working_dir": "/path/to/some_directory",
  "processes":
  [ {
      "name": "process_name",
      "start_command": "/usr/bin/some_start_command",
      "pid_file": "/tmp/some_pid_file.pid",
    }
  ]
}
```

Note: We also set the PWD in the environment to the working dir you specify. This is useful for when the working dir is a symlink. Unicorn in particular will cd into the environment variable in PWD when it re-execs to deal with a change in the symlink.

By default, godard will send a SIGTERM to your process when stopping.
To change the stop command:

```json
{
  "processes":
  [ {
      "name": "process_name",
      "start_command": "/usr/bin/some_start_command",
      "pid_file": "/tmp/some_pid_file.pid",
      "stop_command": "/user/bin/some_stop_command"
    }
  ]
}
```

If you'd like to send a signal or signals to your process to stop it:

```json
{
  "processes":
  [ {
      "name": "process_name",
      "start_command": "/usr/bin/some_start_command",
      "pid_file": "/tmp/some_pid_file.pid",
      "stop_signals": ["quit", "30.seconds", "term", "5.seconds", "kill"]
    }
  ]
}
```

We added a line that will send a SIGQUIT, wait 30 seconds and check to
see if the process is still up, send a SIGTERM, wait 5 seconds and check
to see if the process is still up, and finally send a SIGKILL.

And lastly, to monitor child processes:

```json
"monitor_children": {
  "checks": {
    "mem_usage":{ "every": "10.secs", "below": "100.megabytes", "times": 3 }
  },
  "stop_command": "kill -QUIT {{PID}}"
  }
}
```

Note {{PID}} will be substituted for the pid of process in both the stop and restart commands.

### CLI

#### Usage

    godard COMMAND [flags]

For the "load" command, the _app_name_ is specified in the config file, and
must not be provided on the command line.

For all other commands, the _app_name_ is optional if there is only
one godard daemon running. Otherwise, the _app_name_ must be
provided, because the command will fail when there are multiple
godard daemons running. The example commands below leaves out the
_app_name_.

#### Commands

    load CONFIG_FILE    Loads new instance of godard using the specified config file
    status              Lists the status of the proceses for the specified app
    start [TARGET]      Issues the start command for the target process or group, defaults to all processes
    stop [TARGET]       Issues the stop command for the target process or group, defaults to all processes
    restart [TARGET]    Issues the restart command for the target process or group, defaults to all processes
    unmonitor [TARGET]  Stop monitoring target process or group, defaults to all processes
    log [TARGET]        Show the log for the specified process or group, defaults to all for app
    quit                Stop godard

### Logging
By default, godard uses syslog local6 facility as described in the installation section. But if for any reason you don&apos;t want to use syslog, you can use a log file. You can do this by setting the :log\_file option in the config:

```json
{
  "log_file": "/path/to/godard.log",
  ...
}
```

Keep in mind that you still need to set up log rotation (described in the installation section) to keep the log file from growing huge.
