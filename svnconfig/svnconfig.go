// Package svnconfig generates config files for SVN server.
package svnconfig

import (
	"bytes"
	"text/template"

	"sigs.k8s.io/yaml"
)

var (
	tmplAuthzSVNAccessFile = template.Must(template.New("AuthzSVNAccessFile").Parse(rawTmplAuthzSVNAccessFile))
	tmplAuthUserFile       = template.Must(template.New("AuthUserFile").Parse(rawTmplAuthUserFile))
)

type Config struct {
	Repositories []*Repository
	Groups       []*Group
	Users        []*User
}

type Repository struct {
	Name        string
	Permissions []*Permission
}

type Permission struct {
	Group      string
	Permission string
}

type Group struct {
	Name  string
	Users []string
}

type User struct {
	Name              string
	EncryptedPassword string
}

type ReposConfig struct {
	Repositories []*RepoEntry `json:"repositories"`
}

type RepoEntry struct {
	Name string `json:"name,omitempty"`
}

func (c *Config) AuthzSVNAccessFile() (string, error) {
	buf := bytes.NewBuffer(nil)
	if err := tmplAuthzSVNAccessFile.Execute(buf, c); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *Config) AuthUserFile() (string, error) {
	buf := bytes.NewBuffer(nil)
	if err := tmplAuthUserFile.Execute(buf, c); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *Config) ReposConfig() (string, error) {
	marshaled, err := yaml.Marshal(c.BuildReposConfig())
	if err != nil {
		return "", err
	}
	return string(marshaled), nil
}

func (c *Config) BuildReposConfig() *ReposConfig {
	repos := []*RepoEntry{}
	for _, r := range c.Repositories {
		repos = append(repos, &RepoEntry{Name: r.Name})
	}
	return &ReposConfig{Repositories: repos}
}
