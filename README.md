# portscan
A simple parallel port scanner.
Connect to port and close connection immidiately.
You can define number of parallel connections.
Support Json output
------

## Usage:
```
  portscan [flags]

Flags:
  -c, --connections int   Connection pool (default 50)
  -d, --dst string        Destination address (default "localhost")
  -h, --help              help for portscan
  -j, --json              Json output (ignore -v and -s)
  -p, --port string       Port or range. Can be several ranges/ports. Example: '2,80-100,8080'
  -r, --retries int       How many times check unavailable port (default 2)
  -s, --sort              Sort output ports (print when checks all)
  -t, --timeout int       Timeout in milliseconds for each connection (default 1000)
  -v, --verbose           Print failed ports
```