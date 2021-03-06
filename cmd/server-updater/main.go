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

package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"

	"github.com/genkami/svn-operator/controllers"
	"github.com/genkami/svn-operator/pkg/serverupdater"
)

func main() {
	zapLog, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to initialize logger", err)
		os.Exit(1)
	}
	log := zapr.NewLogger(zapLog)
	// TODO: add command line flag
	u := &serverupdater.Updater{
		InitdScript: "/etc/init.d/apache2",
		SvnAdmin:    "/usr/bin/svnadmin",
		ReposConfig: filepath.Join(controllers.VolumePathConfig, controllers.ConfigMapKeyRepos),
		ReposDir:    filepath.Join(controllers.VolumePathRepos, "repos"),
		Log:         log,
		TimeoutMs:   10000,
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error(err, "failed to initialize watcher")
		os.Exit(1)
	}
	if err := watcher.Add(controllers.VolumePathConfig); err != nil {
		log.Error(err, "failed to watch config files")
		os.Exit(1)
	}

	signals := make(chan os.Signal, 2)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case ev := <-watcher.Events:
			if ev.Op&(fsnotify.Create|fsnotify.Write) == 0 {
				continue
			}
			log.Info("detected config change", "filename", ev.Name)
			err = u.OnConfigChanged()
			if err != nil {
				log.Error(err, "failed to update repository settings")
			}
		case sig := <-signals:
			log.Info("caught signal; quitting", "signal", sig.String())
			break
		}
	}
}
