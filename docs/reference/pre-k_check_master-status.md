## pre-k check master-status

Checks whether master(s) are running and ready

### Synopsis

Checks whether master(s) are running and ready

```
pre-k check master-status [flags]
```

### Options

```
  -h, --help                help for master-status
      --interval duration   Interval between checks (default 2s)
      --kubeconfig string   Path to kubeconfig file with authorization information (the master location is set by the master flag).
      --master string       The address of the Kubernetes API server (overrides any value in kubeconfig)
      --timeout duration    Timeout for check master status
```

### Options inherited from parent commands

```
      --alsologtostderr                  log to standard error as well as files
      --analytics                        Send analytical events to Google Guard (default true)
      --log-flush-frequency duration     Maximum number of seconds between log flushes (default 5s)
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --stderrthreshold severity         logs at or above this threshold go to stderr
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO

* [pre-k check](pre-k_check.md)	 - Check stuff

