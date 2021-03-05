/*
Copyright 2021 Genta Kamitani.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package serverupdater contains functions to update internal state of
// SVN servers.
package serverupdater

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"sigs.k8s.io/yaml"

	"github.com/genkami/svn-operator/pkg/svnconfig"
)

// Updater updates SVN repositories and Apache Servers.
type Updater struct {
	// InitdScript is a path to apache init script (e.g. /etc/init.d/httpd)
	InitdScript string

	// SvnAdmin is a path to the `svnadmin` command.
	SvnAdmin string

	// ReposConfig is a path to a set of definitions of repositories that the server has.
	ReposConfig string

	// ReposDir is a path to a directory that SVN repositories resides in.
	ReposDir string

	// Log is a logger.
	Log logr.Logger

	// TimeoutMs is a timeout in milliseconds to run command.
	TimeoutMs int
}

func (u *Updater) OnAuthUserFileChanged() error {
	return u.reloadApache()
}

func (u *Updater) OnAuthzSVNAccessFileChanged() error {
	return u.reloadApache()
}

func (u *Updater) OnReposConfigChanged() error {
	return u.createRepositories()
}

func (u *Updater) reloadApache() error {
	return u.runCommand(u.InitdScript, "reload")
}

func (u *Updater) createRepositories() error {
	reposConfigFile, err := os.Open(u.ReposConfig)
	if err != nil {
		return err
	}
	defer reposConfigFile.Close()
	rawReposConfig, err := ioutil.ReadAll(reposConfigFile)
	if err != nil {
		return err
	}
	var reposConfig svnconfig.ReposConfig
	err = yaml.Unmarshal(rawReposConfig, &reposConfig)
	if err != nil {
		return err
	}
	for i := range reposConfig.Repositories {
		err = u.createRepository(reposConfig.Repositories[i].Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Updater) createRepository(name string) error {
	dest := filepath.Join(u.ReposDir, name)
	if fileExists(dest) {
		return nil
	}
	return u.runCommand(u.SvnAdmin, "create", dest)
}

func (u *Updater) runCommand(cmd ...string) error {
	log := u.Log.WithValues("command", strings.Join(cmd, " "))
	ctx := context.Background()
	if u.TimeoutMs > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, time.Duration(u.TimeoutMs)*time.Millisecond)
		defer cancel()
	}
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	command := exec.CommandContext(ctx, cmd[0], cmd[1:]...)
	command.Stdout = stdout
	command.Stderr = stderr
	err := command.Run()
	log.Info("command output", "stdout", stdout.String(), "stderr", stderr.String())
	if err != nil {
		log.Error(err, "command error")
		return err
	}
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	} else if errors.Is(err, os.ErrExist) {
		return true
	}
	return false
}
