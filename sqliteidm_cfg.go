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

package sqliteidm

import (
	"database/sql"

	"github.com/avfs/avfs"

	// Sqlite database driver
	_ "github.com/mattn/go-sqlite3"
)

// New create a new identity manager.
func New(db *sql.DB) (*SQLiteIdm, error) {
	ut := avfs.Cfg.Utils()
	adminGroupName := ut.AdminGroupName()
	adminUserName := ut.AdminUserName()

	idm := &SQLiteIdm{
		adminGroup: &Group{
			name: adminGroupName,
			gid:  0,
		},
		adminUser: &User{
			name: adminUserName,
			uid:  0,
			gid:  0,
		},
		db:    db,
		utils: ut,
	}

	err := db.Ping()
	if err != nil {
		return nil, err
	}

	sqlDBCreate := `
	create table if not exists groups
	(
		gid integer primary key autoincrement,
		name text not null unique
	);

	insert into groups(gid, name)
		values
		       (-1, 'invalid group'),
		       (0, ?)
		on conflict do nothing;

	create table if not exists users
	(
		uid integer primary key autoincrement,
		name text not null unique,
		gid integer default -1 not null
			references groups
				on update set default on delete set default
	);

	insert into users(uid, name, gid)
		values (0, ?, 0)
		on conflict do nothing;
	`

	_, err = db.Exec(sqlDBCreate, adminGroupName, adminUserName)
	if err != nil {
		return nil, err
	}

	return idm, nil
}

func (idm *SQLiteIdm) Close() error {
	err := idm.db.Close()

	return err
}

// Type returns the type of the fileSystem or Identity manager.
func (idm *SQLiteIdm) Type() string {
	return "SQLiteIdm"
}

// Features returns the set of features provided by the file system or identity manager.
func (idm *SQLiteIdm) Features() avfs.Features {
	return avfs.FeatIdentityMgr
}

// HasFeature returns true if the file system or identity manager provides a given feature.
func (idm *SQLiteIdm) HasFeature(feature avfs.Features) bool {
	return avfs.FeatIdentityMgr&feature == feature
}
