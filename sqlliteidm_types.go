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

import "database/sql"

type SQLiteIdm struct {
	db          *sql.DB
	userAdd     *sql.Stmt
	userDel     *sql.Stmt
	userLook    *sql.Stmt
	userLookId  *sql.Stmt
	groupAdd    *sql.Stmt
	groupDel    *sql.Stmt
	groupLook   *sql.Stmt
	groupLookId *sql.Stmt
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
