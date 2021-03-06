# kvbase
A simple abstraction layer for key value stores.

Currently supported stores:
- [BadgerDB](https://github.com/dgraph-io/badger)
- [BoltDB](https://github.com/boltdb/bolt)
- [BboltDB](https://github.com/etcd-io/bbolt)

## Getting Started

### Installing

To start using KVBase, install Go and run `go get`:

```sh
$ go get github.com/Wolveix/kvbase/..
```

This will retrieve the library. You can interact with each supported store in the exact same way, simply swap out the `New` initialisation.

<hr>

### Opening a database

All databases are stored within a `Backend` interface, and have the following functions available:

- `Count(bucket string) (int, error)`
- `Create(bucket string, key string, model interface{}) error`
- `Delete(bucket string, key string) error`
- `Get(bucket string, model interface{}) (*map[string]interface{}, error)`
- `Read(bucket string, key string, model interface{}) error`
- `Update(bucket string, key string, model interface{}) error`

To open a database instance, call the package with `New` followed by the type. E.g: `kvbase.NewBadgerDB("data")`:

```go
package main

import (
	"log"

	"github.com/Wolveix/kvbase"
)

func main() {
    db, err := kvbase.NewBadgerDB("data")
    if err != nil {
        log.Fatal(err)
    }
}
```

<hr>

### Counting entries within a bucket

The `Count()` function expects a bucket (as a `string`),:

```go
counter, err := db.Count("users")
if err != nil {
    log.Fatal(err)
}

fmt.Print(counter) //This will output 1.
```

### Creating an entry

The `Create()` function expects a bucket (as a `string`), a key (as a `string`) and a struct containing your data (as an `interface{}`):

```go
type User struct {
	Password string
	Username string
}

user := User{
    "Password123",
    "JohnSmith123",
}

if err := db.Create("users", user.Username, &user); err != nil {
    log.Fatal(err)
}
```
If the key already exists, this will **fail**.

### Deleting an entry

The `Delete()` function expects a bucket (as a `string`), a key (as a `string`):

```go
if err := db.Delete("users", "JohnSmith01"); err != nil {
    log.Fatal(err)
}
```

### Getting all entries

The `Get()` function expects a bucket (as a `string`), and a struct to unmarshal your data into (as an `interface{}`):

```go
results, err := db.Get("users", User{})
if err != nil {
    log.Fatal(err)
}

s, _ := json.MarshalIndent(results, "", "\t")
fmt.Print(string(s))
```

`results` will now contain a `*map[string]interface{}` object. Note that the object doesn't support indexing, so `results["JohnSmith01"]` won't work; however, you can loop through the mapo to find specific keys.

### Reading an entry

The `Read()` function expects a bucket (as a `string`), a key (as a `string`) and a struct to unmarshal your data into (as an `interface{}`):

```go
user := User{}

if err := db.Read("users", "JohnSmith01", &user); err != nil {
    log.Fatal(err)
}

fmt.Print(user.Password) //This will output Password123
```

### Updating an entry

The `Update()` function expects a bucket (as a `string`), a key (as a `string`) and a struct containing your data (as an `interface{}`):

```go
user := User{
    "Password456",
    "JohnSmith123",
}

if err := db.Update("users", user.Username, &user); err != nil {
    log.Fatal(err)
}
```
If the key doesn't already exist, this will **fail**.

<hr>

## Notes

- Since BadgerDB doesn't actually support buckets, prefixes are used in their place.
- While BadgerDB creates a directory for its database, BoltDB and BboltDB create a single file.

## Credits
- Creator: [Robert Thomas](https://github.com/Wolveix)
- License: [GNU General Public License v3.0](https://github.com/Wolveix/kvbase/blob/master/LICENSE)
