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
	"testing"

	"github.com/avfs/avfs/fs/memfs"

	"github.com/avfs/avfs"
	"github.com/avfs/avfs/test"
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

// TestSqliteIdmAll run all tests.
func TestSqliteIdmAll(t *testing.T) {
	idm, err := sqliteidm.New(":memory:")
	if err != nil {
		t.Fatalf("New : want error to be nil, got %v", err)
	}

	defer idm.Close()

	ci := test.NewConfigIdm(t, idm)
	ci.SuiteAll()
}

func TestMemFsWithSqliteIdm(t *testing.T) {
	idm, err := sqliteidm.New(":memory:")
	if err != nil {
		t.Fatalf("sqliteidm.New : want error to be nil, got %v", err)
	}

	defer idm.Close()

	fs, err := memfs.New(memfs.OptMainDirs(), memfs.OptIdm(idm))
	if err != nil {
		t.Fatalf("New : want err to be nil, got %s", err)
	}

	cf := test.NewConfigFs(t, fs)
	cf.SuiteAll()
}
