# Go-Couchbase File Server (GCFS)

GCFS is a metadata management tool for file servers.

## Summary

- [About](#about)
- [GCFS Usage](#gcfs-usage)
- [Configuration](#configuration)
    - [Prerequisite](#prerequisite)
        - [Couchbase instance](#couchbase-instance)
        - [Server instance](#server-instance)
    - [Basic setup](#basic-setup)
    - [Api mode](#api-mode)
    - [Extended setup](#extended-setup)
        - [Database](#database)
        - [Global](#global)
            - [AutoProvide](#autoprovide)
        - [Server](#server)
            - [Port](#port)
            - [Logs](#logs)
        - [Metadata](#metadata-advanced)
            - [Understanding metadata](#understanding-metadata)
            - [Adding your own metadata](#adding-your-own-metadata)
            - [Checking for required metadata](#checking-for-required-metadata)
- [Methods](#methods)
    - [Insert](#insert)
    - [InsertF](#insertf)
    - [Get](#get)
    - [Update](#update)
        - [Update specs](#update-specs)
            - [Remove](#remove)
            - [Upsert](#upsert)
            - [Append](#append)
    - [Delete](#delete)
    - [AutoProvide (method)](#autoprovide-method)
    - [CheckIntegrity](#checkintegrity)
- [Error handling](#error-handling)
- [Upcoming features](#upcoming-features)

## About

GCFS is a tool for building cloud file systems. It provides default Couchbase
utilities to store data about your files, located on a distant server.

GCFS is the link between your UI and the distant file. It will provide client
adequate information for retrieving any files stored on your custom server.

## GCFS Usage

You can use GCFS in two modes : api mode and methods mode. The first gives you
more easy setups, while the second aims for deeper integration in your server
application.

You can use both modes together, depending on your configuration.

## Configuration

### Prerequisite

Before using GCFS, you need a minimal setup.

#### Couchbase instance

You need a running Couchbase instance. If you don't want to provide a
cluster configuration, you need to setup a local Couchbase instance with
the following defaults :

- a server running on 127.0.0.1
- empty username and password to access cluster
- one bucket named "metadata"

To provide custom cluster configuration, please refer to the [database](#database)
section below.

#### Server instance

For Api mode, you need your go app to run in server mode on a port.
By default and for local development, this port is set to 8080.

### Basic setup

```go
package my_package

import "github.com/Alvarios/gcfs"

func main() {
    gcfs.Setup(gcfs.Configuration{})
}
```

Configuration interface can be left empty. However, it is required to fill it
if you want to use the package in api mode.

### Api mode

*coming soon*

### Extended setup

```go
package my_package

import "github.com/Alvarios/gcfs"

func main() {
    gcfs.Setup(gcfs.Configuration{
        Database: gcfs.DbConfig{
            BucketName: "metadata",
            Address: "couchbase://127.0.0.1",
            Username: "",
            Password: "",
            Bucket: nil,
        },
        Global: gcfs.GlobalConfig{
            AutoProvide: false,
        },
        Server: gcfs.ServerConfig{
            Port: "8080",
        },
        Metadata: nil,
    })
}
```

#### Database

Provide information to connect to your Couchbase cluster. See [Couchbase official documentation](https://docs.couchbase.com/go-sdk/current/hello-world/start-using-sdk.html#hello-couchbase)
for more information.

Alternatively, you can do the setup on your own and pass a Bucket pointer
into the Database configuration. You can then ignore other configuration
arguments :

```go
package my_package

import (
    "github.com/Alvarios/gcfs"
    "github.com/couchbase/gocb/v2"
)

func main() {
    var myClust *gocb.Cluster
    var myBucket *gocb.Bucket

    myClust = gocb.Connect(...) // Enter your cluster configuration here.
    myBucket = myClust.Bucket(bucketName)

    gcfs.Setup(gcfs.Configuration{
        Database: gcfs.Database{
            Bucket: myBucket,
        },
    })
}
```

#### Global

Configuration of flags to set the package default behaviors.

##### AutoProvide

Provides a default behavior for the [insert method](#insert-method). Default is
false.

AutoProvide mode allows insert method to automatically fill up some metadata
when not specified in the newly created document.

Only 2 fields currently support the autofill mode : `general.creation_time` and
`general.modification_time`. Both are automatically set to current date.

For more detail about metadata, please refer to [this section](#metadata).

#### Server

Provide server configuration for Api mode.

##### Port

The port to run the server on. Running a server on a port from any domain
makes it accessible via `domain:port` url. Domain is `localhost` for local
server, or `127.x.x.x` (range from `127.0.0.1` to `127.255.255.254`).

#### Metadata (advanced)

##### Understanding metadata

GCFS is a metadata utility, which stores textual information about a distant
file. Metadata is used to represent a file without having to load it.

As a service, GCFS provides some pre-configured metadata. This metadata is
present on every file handled by GCFS, and serves as a standard representation.

Some pre-configured metadata is required, while the rest can be left blank.
They are listed below.

| Metadata | Required | AutoFillable | Type | Description |
| :--- | :--- | :--- | :--- | :--- |
| url | true | - | string | url pointing to the actual file. |
| general.name | true | - | string | name of the file. |
| general.format | true | - | string | format of the file. |
| general.size | - | - | int | size of the file in bytes. |
| general.creation_time | true | true | uint64 | date of the file creation. **(1)** |
| general.modification_time | - | true | uint64 | date of the file last modification. **(1)** |

> **(1)** In milliseconds since January 1st, 1970, 00:00:00 UTC. Be careful as this
format is different from the one used by Go and UNIX systems, calculated from
the same date but using nanoseconds. The reason we use milliseconds is to keep
consistent with javascript and other web standards (since Couchbase is a JSON
document database system).

##### Adding your own metadata

As Couchbase works with the very permissive JSON format, you are totally free
to add any metadata to your file, as long as you leave the default provided
ones. 

You can even add some metadata under the general key. For example, a general.author
to keep track of the file owner. GCFS treats a document like a basic interface{},
so you have total freedom.

Now you may want to work with your own standards. Your system will grow up in
a specific direction, and maybe you'd like to require some more metadata. You
can do a check on your own in methods mode, but that's time and processor
costly. Instead, GCFS provides you a prebuilt and efficient solution, that
also works in Api mode.

##### Checking for required metadata

When adding a new document, GCFS will perform an integrity check to ensure
every required field is present with a correctly typed value.

You can add some required fields to check, with the Metadata field inside the
Configuration interface.

```go
package my_package

import "github.com/Alvarios/gcfs"

func main() {
    gcfs.Setup(gcfs.Configuration{
        Metadata: map[string]interface{}{
            "general": map[string]interface{}{
                "author": "string",
            },
            "version": "string",
        },
    })
}
```

For example, the above configuration will force every document to have a
`general.author` and a `version` field, both of type string. If a document
miss one of the above fields, it will be refused and the insertion will
return an error.

Metadata consist of key-value pairs : key points to a field, and value is a
string representing a go type. To access nested fields, declare the parent
field as a hardcoded `map[string]interface{}`, then every key declared
inside will be considered as a child of its parent.

## Methods

GCFS provides methods for interacting with your metadata. They wrap the
gocb methods and add some useful checks and syntax.

### Insert

`fileId, err := gcfs.Insert(fileMetadata, fileId)`

Insert metadata and returns the id of the generated document.

*arguments*

| Name | Type | Description |
| :--- | :--- | :--- |
| fileMetadata | map[string]interface{} **(1)** | metadata document. |
| fileId | string | optional file id. |

**(1)** The map has to contain the required metadata, if not autofilled.

```go
package my_package

import (
    "github.com/Alvarios/gcfs"
    "log"
)

func main() {
    // Works with AutoFill = true.
    data := map[string]interface{}{
        "url": "/path/to/my/file",
        "general": gcfs.GeneralMetadata{ // You can also pass a map[string]interface{}, as long as it contains all the required fields.
            Name: "my awesome file",
            Format: "txt",
        },
        "another key": 123456,
    }
    
    // Leave id blank will generate a unique id.
    fileId, err := gcfs.Insert(data, "")

    log.Println(err == (*gcfs.Error)(nil)) // true

    // Fails with missing url error.
    data := map[string]interface{}{
        "general": map[string]interface{}{
            "name": "my awesome file",
            "format": "txt",
        },
        "another key": 123456,
    }
    
    fileId, err = gcfs.Insert(data, "")

    log.Println(err == (*gcfs.Error)(nil)) // false
}
```

*return value*

| Name | Type | Description |
| :--- | :--- | :--- |
| fileId | string | id pointing to the newly created document. |
| err | [Error object](#error-handling) | - |

### InsertF

Shortcut for InsertFlagged. Allows to set custom flag and ignore Global
configuration for a specific action.

`fileId, err := gcfs.InsertF(fileMetadata, fileId, autoProvide)`

Insert metadata and returns the id of the generated tuple.

*arguments*

| Name | Type | Description |
| :--- | :--- | :--- |
| fileMetadata | map[string]interface{} | file metadata object. |
| fileId | string | optional file id. |
| autoProvide | bool | set true to autofill missing values in the document. |

*return value*

| Name | Type | Description |
| :--- | :--- | :--- |
| fileId | string | id pointing to the newly created document. |
| err | [Error object](#error-handling) | - |

### Get

`fileMetadata, err := gcfs.Get(fileId)`

Retrieve metadata from database.

*arguments*

| Name | Type | Description |
| :--- | :--- | :--- |
| fileId | string | id pointing to the document. |

*return value*

| Name | Type | Description |
| :--- | :--- | :--- |
| fileMetadata | map[string]interface{} | file metadata object. |
| err | [Error object](#error-handling) | - |

### Update

`timestamp, err := gcfs.Update(fileId, updateSpecs)`

Perform a partial update of a document.

*arguments*

| Name | Type | Description |
| :--- | :--- | :--- |
| fileId | string | id pointing to the document. |
| fileMetadata | [Update specs](#update-specs) | - |

*return value*

| Name | Type | Description |
| :--- | :--- | :--- |
| timestamp | uint64 | timestamp of the update. |
| err | [Error object](#error-handling) | - |

#### Update specs

GCFS provides an easy, declarative way to perform a partial update of your
metadata, using [gocb MutateIn](https://docs.couchbase.com/go-sdk/2.1/howtos/subdocument-operations.html#mutating) function.

```go
updateSpecs := gcfs.UpdateSpec{
    Remove: []string{"key1", "key2", "key4.key4-2"},
    Upsert: map[string]interface{}{
        "key3": "new value",
        "key4": map[string]interface{}{
            "key4-1": "new value",
        },
    },
    Append: map[string]interface{}{
        "array-key": ["newValue1", "newValue2"],
    },
}
```

##### Remove

A list of keys to remove from the document. Support the short dot syntax
for nested keys.

##### Upsert

Update a list of values. Each value in the given map will replace the
equivalent one in the original document. If a path doesn't exist, it will be
created.

You can also use short dot syntax to access nested keys.

```go
updateSpecs := gcfs.UpdateSpec{
    Upsert: map[string]interface{}{
        "key3": "new value",
        "key4.key4-1": "new value",
    },
}
```

##### Append

Append a list of values to an array key.

> 💡 Tip : you can declare a non array value to append a single value. An
 array value will always be treaten with the `HasMultiple` flag set to true
 (see [gocb specs](https://docs.couchbase.com/go-sdk/2.1/howtos/subdocument-operations.html#array-append-and-prepend)).
 Add a nested array if you want to append an array as an array, and not a list
 of values.

### Delete
 
`err := gcfs.Delete(fileId)`

Delete metadata from database. Only returns an error.

*arguments*

| Name | Type | Description |
| :--- | :--- | :--- |
| fileId | string | Id pointing to the document. |

*return value*

| Name | Type | Description |
| :--- | :--- | :--- |
| err | [Error object](#error-handling) | - |

### AutoProvide (method)

`autoProvidedData, err = gcfs.AutoProvide(fileMetadata)`

Auto fill metadata with default values. More details in the [Metadata section](#metadata-advanced).

### CheckIntegrity

`coreMetadata, err := gcfs.CheckIntegrity(fileMetadata)`

Check integrity of a dataset.

> 💡 Tip : coreMetadata represents the default provided metadata.
 More details in the [Metadata section](#metadata-advanced).

## Error handling
 
GCFS returns a pointer to an error object adapted to web servers. This error
object is retrieved from [kushuh-go-utils](https://www.github.com/Alvarios/kushuh-go-utils) package.

| Key | Type | Description |
| :--- | :--- | :--- |
| Code | int | The http status of the error. |
| Message | string | Describes the nature of the error. |
 
 > 💡 Tip : as a custom error object, `err == nil` wont work to check for errors.
Instead, use the following guard clause : `err == (*gcfsresponses.Error)(nil)`.
 
## Upcoming features

- **Configuration > Metadata**
    - Add a short syntax for nested fields "general.author"
    - Ignore spaces and case for type declaration
- **Methods > InsertF**
    - Add Force flag to insert document without any check
    - Add Strict flag to disable any non required field
    - Merge flag arguments into an unique flag interface
- **Methods**
    - Search method and api with nefts package
- **Documentation**
    - Add an expert section to setup a valid test environment
 
 ## Copyright
 2020 Alvarios - [MIT license](https://github.com/Alvarios/gcfs/blob/master/LICENSE)
