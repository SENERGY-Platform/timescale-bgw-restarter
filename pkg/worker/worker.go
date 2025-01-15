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

package worker

import (
	"context"
	"log"
	"time"

	"github.com/SENERGY-Platform/timescale-bgw-restarter/pkg/configuration"
	"github.com/jackc/pgx"
)

func Run(ctx context.Context, config configuration.Config) error {
	conn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     config.PostgresHost,
			Port:     config.PostgresPort,
			Database: config.PostgresDb,
			User:     config.PostgresUser,
			Password: config.PostgresPw,
		},
		MaxConnections: 10,
		AcquireTimeout: 0})
	if err != nil {
		return err
	}
	defer conn.Close()

	row := conn.QueryRow("SELECT next_start from timescaledb_information.jobs WHERE proc_name = 'policy_refresh_continuous_aggregate' AND next_start != '-infinity' ORDER BY next_start asc LIMIT 1;")
	var t time.Time
	err = row.Scan(&t)
	if err != nil {
		return err
	}
	timeBorder := time.Now().Add(-12 * time.Hour)
	if t.After(timeBorder) {
		log.Printf("Next job schedule is set for %v, which is after border time of %v. Not performing any action.", t.Format(time.RFC3339), timeBorder.Format(time.RFC3339))
		return nil
	}
	log.Printf("Next job schedule is set for %v, which is before border time of %v. Restarting background workers!", t.Format(time.RFC3339), timeBorder.Format(time.RFC3339))
	_, err = conn.Exec("SELECT _timescaledb_internal.restart_background_workers();")
	if err != nil {
		return err
	}
	return nil
}
