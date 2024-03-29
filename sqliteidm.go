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

// Package sqliteidm is an Identity manager backed by a SQLite database for AVFS.
package sqliteidm

import (
	"database/sql"

	"github.com/avfs/avfs"
	"github.com/mattn/go-sqlite3"
)

// AdminGroup returns the administrator (root) group.
func (idm *SQLiteIdm) AdminGroup() avfs.GroupReader {
	return idm.adminGroup
}

// AdminUser returns the administrator (root) user.
func (idm *SQLiteIdm) AdminUser() avfs.UserReader {
	return idm.adminUser
}

// GroupAdd adds a new group.
func (idm *SQLiteIdm) GroupAdd(name string) (avfs.GroupReader, error) {
	stmt, err := idm.db.Prepare("insert into groups(name) values (?)")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	r, err := stmt.Exec(name)
	if err != nil {
		if e, ok := err.(sqlite3.Error); ok && e.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, avfs.AlreadyExistsGroupError(name)
		}

		return nil, err
	}

	gid, _ := r.LastInsertId()

	g := &Group{
		name: name,
		gid:  int(gid),
	}

	return g, nil
}

// GroupDel deletes an existing group.
func (idm *SQLiteIdm) GroupDel(name string) error {
	stmt, err := idm.db.Prepare("delete from groups where name = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	r, err := stmt.Exec(name)
	if err != nil {
		return err
	}

	n, _ := r.RowsAffected()
	if n == 0 {
		return avfs.UnknownGroupError(name)
	}

	return nil
}

// LookupGroup looks up a group by name.
// If the group cannot be found, the returned error is of type UnknownGroupError.
func (idm *SQLiteIdm) LookupGroup(name string) (avfs.GroupReader, error) {
	stmt, err := idm.db.Prepare("select gid from groups where name = ?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(name)

	var gid int

	err = row.Scan(&gid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, avfs.UnknownGroupError(name)
		}

		return nil, err
	}

	g := &Group{
		name: name,
		gid:  gid,
	}

	return g, nil
}

// LookupGroupId looks up a group by groupid.
// If the group cannot be found, the returned error is of type UnknownGroupIdError.
func (idm *SQLiteIdm) LookupGroupId(gid int) (avfs.GroupReader, error) {
	stmt, err := idm.db.Prepare("select name from groups where gid = ?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(gid)

	var name string

	err = row.Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, avfs.UnknownGroupIdError(gid)
		}

		return nil, err
	}

	g := &Group{
		name: name,
		gid:  gid,
	}

	return g, nil
}

// LookupUser looks up a user by username.
// If the user cannot be found, the returned error is of type UnknownUserError.
func (idm *SQLiteIdm) LookupUser(name string) (avfs.UserReader, error) {
	stmt, err := idm.db.Prepare("select uid, gid from users where name = ?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(name)

	var uid, gid int

	err = row.Scan(&uid, &gid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, avfs.UnknownUserError(name)
		}

		return nil, err
	}

	u := &User{
		name: name,
		uid:  uid,
		gid:  gid,
	}

	return u, nil
}

// LookupUserId looks up a user by userid.
// If the user cannot be found, the returned error is of type UnknownUserIdError.
func (idm *SQLiteIdm) LookupUserId(uid int) (avfs.UserReader, error) {
	stmt, err := idm.db.Prepare("select name, gid from users where uid = ?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(uid)

	var (
		name string
		gid  int
	)

	err = row.Scan(&name, &gid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, avfs.UnknownUserIdError(uid)
		}

		return nil, err
	}

	u := &User{
		name: name,
		uid:  uid,
		gid:  gid,
	}

	return u, nil
}

// UserAdd adds a new user.
func (idm *SQLiteIdm) UserAdd(name, groupName string) (avfs.UserReader, error) {
	g, err := idm.LookupGroup(groupName)
	if err != nil {
		return nil, err
	}

	stmt, err := idm.db.Prepare("insert into users(name, gid) values (?, ?)")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	r, err := stmt.Exec(name, g.Gid())
	if err != nil {
		if e, ok := err.(sqlite3.Error); ok && e.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, avfs.AlreadyExistsUserError(name)
		}

		return nil, err
	}

	uid, _ := r.LastInsertId()

	u := &User{
		name: name,
		uid:  int(uid),
		gid:  g.Gid(),
	}

	return u, nil
}

// UserDel deletes an existing group.
func (idm *SQLiteIdm) UserDel(name string) error {
	stmt, err := idm.db.Prepare("delete from users where name = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	r, err := stmt.Exec(name)
	if err != nil {
		return err
	}

	n, _ := r.RowsAffected()
	if n == 0 {
		return avfs.UnknownUserError(name)
	}

	return nil
}

// User

// Name returns the user name.
func (u *User) Name() string {
	return u.name
}

// Gid returns the primary group ID of the user.
func (u *User) Gid() int {
	return u.gid
}

// IsAdmin returns true if the user has administrator (root) privileges.
func (u *User) IsAdmin() bool {
	return u.uid == 0 || u.gid == 0
}

// Uid returns the user ID.
func (u *User) Uid() int {
	return u.uid
}

// Group

// Gid returns the group ID.
func (g *Group) Gid() int {
	return g.gid
}

// Name returns the group name.
func (g *Group) Name() string {
	return g.name
}
