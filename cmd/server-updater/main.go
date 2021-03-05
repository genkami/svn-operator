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
	"path/filepath"
	"time"

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

	for {
		// TODO: watch files and reload only if necessary
		err = u.OnReposConfigChanged()
		if err != nil {
			log.Error(err, "failed to update repos")
		}
		err = u.OnAuthUserFileChanged()
		if err != nil {
			log.Error(err, "failed to reload apache")
		}
		err = u.OnAuthzSVNAccessFileChanged()
		if err != nil {
			log.Error(err, "failed to reload apache")
		}
		time.Sleep(10 * time.Second)
	}
}
