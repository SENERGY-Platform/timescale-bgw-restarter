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

package pkg

import (
	"context"
	"log"
	"sync"

	"github.com/SENERGY-Platform/timescale-bgw-restarter/pkg/configuration"
	"github.com/SENERGY-Platform/timescale-bgw-restarter/pkg/worker"
)

func Start(ctx context.Context, config configuration.Config) (wg *sync.WaitGroup, err error) {
	wg = &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = worker.Run(ctx, config)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return
}
