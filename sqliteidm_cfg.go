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
	_ "github.com/mattn/go-sqlite3" // Sqlite database driver
)

// New create a new identity manager.
func New(db *sql.DB) (*SQLiteIdm, error) {
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
		       (0, 'root')
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
		values (0, 'root', 0)
		on conflict do nothing;
	`

	_, err = db.Exec(sqlDBCreate)
	if err != nil {
		return nil, err
	}

	groupAdd, err := db.Prepare("insert into groups(name) values (?)")
	if err != nil {
		return nil, err
	}

	groupDel, err := db.Prepare("delete from groups where name = ?")
	if err != nil {
		return nil, err
	}

	groupLook, err := db.Prepare("select gid from groups where name = ?")
	if err != nil {
		return nil, err
	}

	groupLookId, err := db.Prepare("select name from groups where gid = ?")
	if err != nil {
		return nil, err
	}

	userAdd, err := db.Prepare("insert into users(name, gid) values (?, ?)")
	if err != nil {
		return nil, err
	}

	userDel, err := db.Prepare("delete from users where name = ?")
	if err != nil {
		return nil, err
	}

	userLook, err := db.Prepare("select uid, gid from users where name = ?")
	if err != nil {
		return nil, err
	}

	userLookId, err := db.Prepare("select name, gid from users where uid = ?")
	if err != nil {
		return nil, err
	}

	idm := &SQLiteIdm{
		db:          db,
		userAdd:     userAdd,
		userDel:     userDel,
		userLook:    userLook,
		userLookId:  userLookId,
		groupAdd:    groupAdd,
		groupDel:    groupDel,
		groupLook:   groupLook,
		groupLookId: groupLookId,
	}

	return idm, nil
}

func (idm *SQLiteIdm) Close() error {
	_ = idm.groupAdd.Close()
	_ = idm.groupDel.Close()
	_ = idm.groupLook.Close()
	_ = idm.groupLookId.Close()
	_ = idm.userAdd.Close()
	_ = idm.userDel.Close()
	_ = idm.userLook.Close()
	_ = idm.userLookId.Close()
	err := idm.db.Close()

	return err
}

// Type returns the type of the fileSystem or Identity manager.
func (idm *SQLiteIdm) Type() string {
	return "SQLiteIdm"
}

// Features returns the set of features provided by the file system or identity manager.
func (idm *SQLiteIdm) Features() avfs.Feature {
	return avfs.FeatIdentityMgr
}

// HasFeature returns true if the file system or identity manager provides a given feature.
func (idm *SQLiteIdm) HasFeature(feature avfs.Feature) bool {
	return avfs.FeatIdentityMgr&feature == feature
}
