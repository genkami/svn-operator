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

// Generator generates configuration files for SVN server.
//
// NB: Generator assumes that all fields and parameters are VALIDATED ELSEWHERE and does not validate nor escape
// any fields. Be careful if you want to generate configuration files from untrusted source.
type Generator struct {
	Repositories []Repository
	Groups       []Group
	Users        []User
}

// Repository is a definition of a repository.
type Repository struct {
	Name        string
	Permissions []Permission
}

// Permission configurates permission to a specific repository.
type Permission struct {
	Group      string
	Permission string
}

// Group is a definitions of a group.
type Group struct {
	Name  string
	Users []string
}

// User is a definition of a user.
type User struct {
	Name              string
	EncryptedPassword string
}

// ReposConfig is a special configuration structure that is used to create SVN repositories.
type ReposConfig struct {
	Repositories []RepoEntry `json:"repositories"`
}

// RepoEntry is an entry for SVN repository.
type RepoEntry struct {
	Name string `json:"name,omitempty"`
}

// AuthzSVNAccessFile is an authorization configuration file for mod_authz_svn.
//
// See https://svn.apache.org/repos/asf/subversion/trunk/subversion/mod_authz_svn/INSTALL for more details.
func (g *Generator) AuthzSVNAccessFile() (string, error) {
	buf := bytes.NewBuffer(nil)
	if err := tmplAuthzSVNAccessFile.Execute(buf, g); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// AuthUserFile is an authentication configuration file for mod_authn_file.
//
// See https://httpd.apache.org/docs/2.4/en/mod/mod_authn_file.html for more details.
func (g *Generator) AuthUserFile() (string, error) {
	buf := bytes.NewBuffer(nil)
	if err := tmplAuthUserFile.Execute(buf, g); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (g *Generator) ReposConfig() (string, error) {
	marshaled, err := yaml.Marshal(g.BuildReposConfig())
	if err != nil {
		return "", err
	}
	return string(marshaled), nil
}

func (g *Generator) BuildReposConfig() *ReposConfig {
	repos := []RepoEntry{}
	for _, r := range g.Repositories {
		repos = append(repos, RepoEntry{Name: r.Name})
	}
	return &ReposConfig{Repositories: repos}
}
