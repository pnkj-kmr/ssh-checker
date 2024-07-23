# ssh-checker

ssh-checker helps to run multiple commands on end router (switch) in bulk.

_Basically, if you want to run few commands on end router and you want to do for multiple routers then ssh-checker helps to perform the action with input file and generate a output file as result_

### HOW TO USE

_Download the relevent os package from [here](https://github.com/pnkj-kmr/ssh-checker/releases)_

_create a **input.json** file_

```
[
    {
        "host": "127.0.0.1",
        "Tag": "",                  # omitempty
        "commands": ["cmd 1", "cmd 2"],     
        "port": 22,                 # omitempty default: 22
        "username": "admin",        # omitempty default: admin
        "password": "admin",        # omitempty default: admin
        "timeout": 5                # omitempty default: 30
    },
    ...
]
```

_After creating the file run the executable binary as_

```
./sshchecker
```

### OUTPUT


_As a result **output.json** file will be created after completion_

```
[
  {
    "input": {
      "host": "127.0.0.1",
      "Tag": "",
      "commands": ["cmd 1", "cmd 2"],     
      "port": 22,
      "username": "admin",
      "password": "admin",
      "timeout": 5
    },
    "error": ["ssh: handshake failed: EOF"],
    "output": []        # result success result if any
  }
]

```

### HELP

```
./sshchecker --help

----------------------
Usage of ./sshchecker:
  -f string
        input file name (default "input.json")
  -knownhosts string
        .ssh known hosts file path:($HOME)/.ssh/known_hosts
  -o string
        output file name (default "output.json")
  -p int
        generic default ssh port (default 22)
  -passwd string
        generic password for connection (default "admin")
  -rsa string
        .ssh file path: ($HOME)/.ssh/id_rsa
  -t int
        timeout [secs] (default 30)
  -usr string
        generic username for connection (default "admin")
  -w int
        number of workers (default 4)

-------
Example:

./sshchecker -f x.json -t 30 -w 20

```

:)
