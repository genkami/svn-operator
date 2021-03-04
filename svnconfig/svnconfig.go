// Package svnconfig generates config files for SVN server.
package svnconfig

import (
	"bytes"
	"text/template"
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
	Name         string
	PasswordSHA1 string
}

func (c *Config) AuthzSVNAccessFile() (string, error) {
	buf := bytes.NewBuffer(nil)
	if err := tmplAuthzSVNAccessFile.Execute(buf, c); err != nil {
		return "", err
	}
	return buf.String(), nil
}
