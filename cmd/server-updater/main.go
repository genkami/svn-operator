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
	"flag"
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
	var initdScript, svnAdmin string
	var timeoutMs int
	flag.StringVar(&initdScript, "initd-script", "/etc/init.d/apache2", "Path to /etc/init.d/apache2 (or its variant)")
	flag.StringVar(&svnAdmin, "svnadmin", "/usr/bin/svnadmin", "Path to `svnadmin` command")
	flag.IntVar(&timeoutMs, "exec-timeout", 10000, "Timeout to run commands")
	flag.Parse()

	zapLog, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to initialize logger", err)
		os.Exit(1)
	}
	log := zapr.NewLogger(zapLog)

	u := &serverupdater.Updater{
		InitdScript: initdScript,
		SvnAdmin:    svnAdmin,
		ReposConfig: filepath.Join(controllers.VolumePathConfig, controllers.ConfigMapKeyRepos),
		ReposDir:    filepath.Join(controllers.VolumePathRepos, "repos"),
		TimeoutMs:   timeoutMs,
		Log:         log,
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
	// Volumes made from ConfigMaps stores all values inside ConfigMaps in `..data` directory.
	// See https://github.com/kubernetes/kubernetes/blob/master/pkg/volume/util/atomic_writer.go
	dataDir := filepath.Join(controllers.VolumePathConfig, "..data")

	signals := make(chan os.Signal, 2)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	log.Info("initializing")
	err = u.OnConfigChanged()
	if err != nil {
		log.Error(err, "failed to initialize settings")
	}

	for {
		select {
		case ev := <-watcher.Events:
			if ev.Op&(fsnotify.Create|fsnotify.Write) == 0 {
				continue
			}
			if ev.Name != dataDir {
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
