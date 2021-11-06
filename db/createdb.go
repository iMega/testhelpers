// Copyright © 2020 Dmitry Stoletov <info@imega.ru>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/imega/txwrapper"

	// init driver.
	_ "github.com/mattn/go-sqlite3"
)

type CloseFn func() error

// Create creates the sqlite database in temporary file.
func Create(dbName string, txFunc txwrapper.TxFunc) (*sql.DB, CloseFn, error) {
	if dbName == "" {
		dbName = "db"
	}

	file, err := ioutil.TempFile("", dbName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create tmp file, %w", err)
	}

	filename := file.Name()

	if err := file.Close(); err != nil {
		return nil, nil, fmt.Errorf("failed to close file, %w", err)
	}

	dbSL, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open db, %w", err)
	}

	wrapper := txwrapper.New(dbSL)

	err = wrapper.Transaction(context.Background(), nil, txFunc)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute, %w", err)
	}

	closeFn := func() error {
		errDB := dbSL.Close()
		if err := os.Remove(filename); err != nil || errDB != nil {
			return fmt.Errorf(
				"failed to close db or remove temp file, %s, %w", errDB, err,
			)
		}

		return nil
	}

	return dbSL, closeFn, nil
}
