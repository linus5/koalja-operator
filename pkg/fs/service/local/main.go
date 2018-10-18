/*
Copyright 2018 Aljabr Inc.

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
	"context"
	"log"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	fssvc "github.com/AljabrIO/koalja-operator/pkg/fs/service"
	"github.com/AljabrIO/koalja-operator/pkg/util"
)

var (
	cliLog = util.MustCreateLogger()
)

// TODO: Add cleanup of files

func main() {
	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to get kubernetes API server config")
	}

	// Create a new Cmd to provide shared dependencies and start components
	localPathPrefix := "/var/lib/koalja/local-fs"
	storageClassName := "koalja-local-storage"
	scheme := "koalja-file"
	builder := newLocalFileSystemBuilder(localPathPrefix, storageClassName, scheme)
	svc, err := fssvc.NewService(cfg, builder)
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to create FileSystem service")
	}

	log.Printf("Starting the Cmd.")

	// Start the Cmd
	ctx, done := context.WithCancel(context.Background())
	go func() {
		<-signals.SetupSignalHandler()
		done()
	}()
	if err := svc.Run(ctx); err != nil {
		cliLog.Fatal().Err(err).Msg("Service failed")
	}
}
