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
)

type SQLiteIdm struct {
	adminGroup  *Group     // Administrator Group.
	adminUser   *User      // Administrator User.
	db          *sql.DB    // db is the SQLite database handle/
	userAdd     *sql.Stmt  // userAdd is the prepared statement to add a user.
	userDel     *sql.Stmt  // userDel is the prepared statement to delete a user.
	userLook    *sql.Stmt  // userLook is the prepared statement to look up a user by name.
	userLookId  *sql.Stmt  // userLookId is the prepared statement to look up a user by id.
	groupAdd    *sql.Stmt  // groupAdd is the prepared statement to add a group.
	groupDel    *sql.Stmt  // groupDel is the prepared statement to delete a group.
	groupLook   *sql.Stmt  // groupLook is the prepared statement to look up a group by name.
	groupLookId *sql.Stmt  // groupLookId is the prepared statement to look up a group by id.
	utils       avfs.Utils // Utils regroups common functions used by emulated file systems.
}

// User is the implementation of avfs.UserReader.
type User struct {
	name string
	uid  int
	gid  int
}

// Group is the implementation of avfs.GroupReader.
type Group struct {
	name string
	gid  int
}
