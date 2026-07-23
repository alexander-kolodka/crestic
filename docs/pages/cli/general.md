# ⚙️ General
## `-c, --config`

Specify the config file to be used.
If omitted `crestic` will search for a `crestic.yaml` in the current directory, your home folder and your config folder:
- ./crestic.yaml
- ~/crestic.yaml
- ~/.crestic/crestic.yaml
- ~/.config/crestic/crestic.yaml


```bash
crestic -c /path/to/my/config.yaml backup --all
```

## `--log-level`

Set the log level. Available levels: `debug`, `info`, `warn`, `error` (default: `info`).

```bash
crestic --log-level debug backup --all
```

## `--ci`

Output logs as plain text without colors. Useful for CI pipelines and log collectors.

Log level is unchanged — use `--log-level` if you need more verbose output.

```bash
crestic --ci backup --all
```

## `--print-commands`

Print executed shell commands. Useful for debugging.

```bash
crestic --print-commands backup --all
```
