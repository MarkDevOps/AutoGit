package types

import (
	"time"
)

// Used for the fetch Command
type Config struct {
	Org   string              `yaml:"org"`
	Repos map[string][]string `yaml:"repos"`
}

// Used for outputs
type OutputData struct {
	Organization string                        `yaml:"organization"`
	Repositories map[string]map[string]EnvData `yaml:"repositories"`
}

type EnvData struct {
	DeploymentID  int    `yaml:"deployment_id,omitempty"`
	Ref           string `yaml:"ref,omitempty"`
	Description   string `yaml:"description,omitempty"`
	CreatedAt     string `yaml:"created_at,omitempty"`
	Status        string `yaml:"status,omitempty"`
	StatusTime    string `yaml:"status_time,omitempty"`
	DeploymentURL string `yaml:"deployment_url,omitempty"`
}

type ReleaseData struct {
	TagName     string `yaml:"tag_name"`
	Name        string `yaml:"name"`
	PublishedAt string `yaml:"published_at"`
	HTMLURL     string `yaml:"html_url"`
}

type WorkflowData struct {
	Name      string `yaml:"name"`
	State     string `yaml:"state"`
	HTMLURL   string `yaml:"html_url"`
	UpdatedAt string `yaml:"updated_at"`
}

type Deployment struct {
	ID          int       `json:"id"`
	Ref         string    `json:"ref"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	StatusesURL string    `json:"statuses_url"`
	Timestamp   int64     `json:"timestamp"`
	Environment string    `json:"environment"`
}

// Used for the secrets Command
type ConfigSecrets struct {
	Org   string                        `yaml:"org"`
	Repos map[string]map[string]Secrets `yaml:"repos"`
}

// Used for the secrets Command
type Secrets struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
