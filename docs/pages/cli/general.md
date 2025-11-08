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

Run the CLI in CI Mode, which means there will be no interactivity, no colors and automatically sets verbose logging.

This can be useful when you want to run cron e.g. as all the output will be saved.

```bash
crestic --ci backup --all
```

## `--json`

Output logs in JSON format. Useful for log aggregation systems.

```bash
crestic --json backup --all
```

## `--print-commands`

Print executed shell commands. Useful for debugging.

```bash
crestic --print-commands backup --all
```
