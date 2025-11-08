package dto

// Config is the YAML configuration structure.
type Config struct {
	Repositories   map[string]Repository `yaml:"repositories"`
	Jobs           Jobs                  `yaml:"jobs"`
	HealthcheckURL string                `yaml:"healthcheck_url"`
}

type Options map[string]any

type BackupJob struct {
	Name                     string   `yaml:"name"`
	Cron                     string   `yaml:"cron"`
	IgnoreMissingXAttrsError bool     `yaml:"ignore_x_attrs_error"`
	From                     []string `yaml:"from"`
	To                       string   `yaml:"to"`
	Options                  Options  `yaml:"options"`
	Hooks                    Hooks    `yaml:"hooks"`
	HealthcheckURL           string   `yaml:"healthcheck_url"`
}

type CopyJob struct {
	Name           string  `yaml:"name"`
	Cron           string  `yaml:"cron"`
	From           string  `yaml:"from"`
	To             string  `yaml:"to"`
	Options        Options `yaml:"options"`
	Hooks          Hooks   `yaml:"hooks"`
	HealthcheckURL string  `yaml:"healthcheck_url"`
}

type Repository struct {
	Path          string  `yaml:"path"`
	PasswordCMD   string  `yaml:"password_command"`
	ForgetOptions Options `yaml:"forget_options"`
}

type Hooks struct {
	Before  []string `yaml:"before"`
	Failure []string `yaml:"failure"`
	Success []string `yaml:"success"`
}
