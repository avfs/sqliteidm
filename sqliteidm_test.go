//
//  Copyright 2020 The AVFS authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//  	http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//

package sqliteidm_test

import (
	"database/sql"
	"testing"

	"github.com/avfs/avfs"
	"github.com/avfs/avfs/test"
	"github.com/avfs/avfs/vfs/memfs"
	"github.com/avfs/sqliteidm"
)

var (

	// MemIdm implements avfs.IdentityMgr interface.
	_ avfs.IdentityMgr = &sqliteidm.SQLiteIdm{}

	// User implements avfs.UserReader interface.
	_ avfs.UserReader = &sqliteidm.User{}

	// Group implements avfs.GroupReader interface.
	_ avfs.GroupReader = &sqliteidm.Group{}
)

// InitDB initialize the database.
func InitDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("sql.Open : want error to be nil, got %v", err)
	}

	return db
}

// TestSqliteIdmAll run all tests.
func TestSqliteIdmAll(t *testing.T) {
	db := InitDB(t)

	idm, err := sqliteidm.New(db)
	if err != nil {
		t.Fatalf("New : want error to be nil, got %v", err)
	}

	defer idm.Close()

	sidm := test.NewSuiteIdm(t, idm)
	sidm.All(t)
}

func TestMemFsWithSqliteIdm(t *testing.T) {
	db := InitDB(t)

	idm, err := sqliteidm.New(db)
	if err != nil {
		t.Fatalf("sqliteidm.New : want error to be nil, got %v", err)
	}

	defer idm.Close()

	fs, err := memfs.New(memfs.WithMainDirs(), memfs.WithIdm(idm))
	if err != nil {
		t.Fatalf("New : want err to be nil, got %s", err)
	}

	sfs := test.NewSuiteFS(t, fs)
	sfs.All(t)
}
