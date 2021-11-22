---
id: user
title: Query an User Information
---

To query credential associated with an user, for now, we need to use a user administration transaction context.
The `GetUser()` method in the `UsersTxContext` retrieve a given key from a given database.

```go
type UsersTxContext interface {
    .
    ...
	// GetUser obtain user's record from database
	GetUser(userID string) (*types.User, error)
    ...
    .
}
```

:::info pre-requisite
For the example shown here to work, we need to have

 1. Two databases named `db1` and `db2` in the orion-server. If you have not created these two databases,
refer to [creates databases using SDK](./dbtx#creation-of-databases) to create `db1` and `db2`.
 2. Two users named `alice` and `bob`. If you have not created these users already, refer to [create users using SDK](./usertx#2-addition-of-users).

Finally, [Create a connection](../../pre-requisite/gosdk#creating-a-connection-to-the-orion-cluster) and
[Open a database session](../../pre-requisite/gosdk#opening-a-database-session).
:::

### Source Code

```go
package main

import (
	"fmt"
	"strconv"
)

func main() {
	bcdb, err := createConnection()
    // if err is not nil, print and return

	session, err := openSession(bcdb)
    // if err is not nil, print and return

	tx, err := session.UsersTx()
    // if err is not nil, print and return

	alice, err := tx.GetUser("alice")
    // if err is not nil, print and return

	err = dbtx.Abort()
    // if err is not nil, print and return
}
```

### (3.2) Source Code Commentry
For simplicity, not all errors are handled in this code. Further, the implementation of `createConnection()` and `openSession()` can be found
[here](../../pre-requisite/gosdk).

The `session.UsersTx()` starts a new user administration transaction and returns the user administration transaction context. We can then query
credentials using this transaction context.

The `GetUser("alice")` fetches the alice credentials from Orion server. The structure of the returned user is shown below:

```go
type User struct {
	Id                   string
	Certificate          []byte
	Privilege            *Privilege
}

type Privilege struct {
	DbPermission map[string]Privilege_Access
    // admin has privileges to submit a user administration transaction,
	// cluster configuration transaction, and database administration
	// transaction. Further, admin has permission to read-write states
	// from any database provided that the state has no ACL defined. If
	// a state has a read and write ACL, the admin can read or write to
	// the state only if the admin is listed in the read or write ACL list.
	Admin                bool
}
```
