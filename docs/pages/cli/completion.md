# ⌨️ Completion

```bash
crestic completion <shell>
```

Generate shell completion scripts for bash, zsh, fish, and PowerShell.

## Supported Shells

- bash
- zsh
- fish
- powershell

## Bash

### Load Completions for Current Session

```bash
source <(crestic completion bash)
```

### Load Completions Permanently

#### Linux

```bash
crestic completion bash > /etc/bash_completion.d/crestic
```

#### macOS

```bash
crestic completion bash > $(brew --prefix)/etc/bash_completion.d/crestic
```

**Note**: If using Homebrew, the completions directory might be:
- `/opt/homebrew/etc/bash_completion.d/` (Apple Silicon)
- `/usr/local/etc/bash_completion.d/` (Intel)

After installation, restart your terminal or run:
```bash
source $(brew --prefix)/etc/bash_completion.d/crestic
```

## Zsh

### Enable Shell Completion (if not already enabled)

If shell completion is not already enabled in your environment, you will need to enable it. Execute the following once:

```bash
echo "autoload -U compinit; compinit" >> ~/.zshrc
```

Then restart your terminal or run:
```bash
source ~/.zshrc
```

### Load Completions for Current Session

```bash
source <(crestic completion zsh)
```

### Load Completions Permanently

```bash
crestic completion zsh > "${fpath[1]}/_crestic"
```

**Note**: `${fpath[1]}` is typically `~/.zsh/completion/` or `/usr/local/share/zsh/site-functions/`

If the directory doesn't exist, create it first:
```bash
mkdir -p ~/.zsh/completion
crestic completion zsh > ~/.zsh/completion/_crestic
```

Then add to your `~/.zshrc`:
```bash
fpath=(~/.zsh/completion $fpath)
```

**You will need to start a new shell for this setup to take effect.**

## Fish

### Load Completions for Current Session

```bash
crestic completion fish | source
```

### Load Completions Permanently

```bash
crestic completion fish > ~/.config/fish/completions/crestic.fish
```

**Note**: If the directory doesn't exist, create it first:
```bash
mkdir -p ~/.config/fish/completions
crestic completion fish > ~/.config/fish/completions/crestic.fish
```

Fish will automatically load completions from this directory.

## PowerShell

### Load Completions for Current Session

```powershell
crestic completion powershell | Out-String | Invoke-Expression
```

### Load Completions Permanently

1. Generate the completion script:
   ```powershell
   crestic completion powershell > crestic.ps1
   ```

2. Source this file from your PowerShell profile:
   ```powershell
   # Find your profile location
   $PROFILE
   
   # Add to profile (if it exists)
   . $HOME\crestic.ps1
   
   # Or create profile and add
   if (!(Test-Path -Path $PROFILE)) {
     New-Item -ItemType File -Path $PROFILE -Force
   }
   Add-Content -Path $PROFILE -Value ". $HOME\crestic.ps1"
   ```

3. Reload your profile:
   ```powershell
   . $PROFILE
   ```

## See Also

- [CLI Commands](/cli) - All available commands
- [Configuration Guide](/config) - Configuration reference
