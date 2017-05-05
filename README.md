# mackerel-plugin-chrony

chrony custom metrics plugin for mackerel.io agent.

Inspired by [telegraf chrony Input Plugin](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/chrony).

## Requirements

- [chrony](https://chrony.tuxfamily.org/)
  - Executable `chronyc` is required

## Synopsis

```shell
mackerel-plugin-chrony [-c <path to chrony>]
```

```shell
Usage of ./mackerel-plugin-chrony:
  -command string
    	path to chronyc (default "/usr/bin/chronyc")
```


## Example of mackerel-agent.conf
```toml
[plugin.metrics.chrony]
command = "/usr/local/bin/mackerel-plugin-chrony"
```

## Metrics
See [the official docs](https://chrony.tuxfamily.org/doc/3.1/chronyc.html#_system_clock) for the various headers returned by this plugin.
