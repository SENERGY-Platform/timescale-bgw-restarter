/*
 *    Copyright 2024 InfAI (CC SES)
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/SENERGY-Platform/go-service-base/struct-logger/attributes"
	"github.com/SENERGY-Platform/timescale-bgw-restarter/pkg"
	"github.com/SENERGY-Platform/timescale-bgw-restarter/pkg/configuration"
	"github.com/SENERGY-Platform/timescale-bgw-restarter/pkg/log"
)

func main() {
	log.Init()
	configLocation := flag.String("config", "config.json", "configuration file")
	flag.Parse()

	config, err := configuration.Load(*configLocation)
	if err != nil {
		log.Logger.Error("failed to load configuration", attributes.ErrorKey, err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	wg, err := pkg.Start(ctx, config)
	if err != nil {
		log.Logger.Error("failed to start application", attributes.ErrorKey, err)
		os.Exit(1)
	}

	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		sig := <-shutdown
		log.Logger.Info("received shutdown signal", "signal", sig)
		cancel()
	}()

	wg.Wait()
}
