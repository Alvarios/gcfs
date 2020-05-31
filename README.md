# Go-Couchbase File Server (GCFS)

GCFS allows easy metadata management for file servers.

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

## Basic configuration

No parameter is required. You can run GCFS without any config file, though
you'll have to comply to default values.

*Full example (with default values)*
```json
{
  "database": {
    "address": "127.0.0.1",
    "username": "",
    "password": "",
    "bucket_name": "metadata"
  },
  "global": {
    "debug": false,
    "auto_provide": false
  },
  "server": {
    "port": "8080",
    "logs": {
      "default": "",
      "error": ""
    }
  },
  "routes": {
    "ping": "",
    "ping_database": "",
    "insert": "",
    "delete": "",
    "get": "",
    "update": "",
    "search": ""
  },
  "metadata": {}
}
```

### Database

Couchbase server configuration.

| Key | Default | Description |
| :--- | :--- | :--- |
| address | 127.0.0.1 | Address of Couchbase Cluster. |
| username | - | Credential for the Cluster. |
| password | - | Credential for the Cluster. |
| bucket_name | metadata | Bucket to save metadata to. |

### Global

| Key | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| debug | bool | false | Log debug information. |
| auto_provide | bool | false | On insertion and upsertion, autofill some metadata. Only has effect in methods mode (always true in api mode). |

### Server

Server configuration.

| Key | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| port | string | 8080 | Port to access server. |
| logs | object | {} | Define log files (print to console by default). |
| logs.default | string | - | Path to main log file. |
| logs.error | string | - | Path to error log file. |

### Routes

If any route is set, gcfs will create a pre-configured route for the specified
action. You don't have to set every route, only the ones you need.

| Route | Method | URL Params | Body | Return | Description |
| :--- | :--- | :--- | :--- | :--- | :--- |
| ping | Get | - | - | "pong" | Ping route. |
| ping_database | Get | - | - | [Couchbase health check data](https://docs.couchbase.com/go-sdk/2.1/concept-docs/health-check.html). | Ping route for Couchbase server. |
| insert | Post | - | [Metadata](#metadata) | id: file id. | Insert a file metadata to the server. |
| delete | Delete | id: file id. | - | - | Delete a file metadata from the server. |
| get | Get | id: file id. | - | [Metadata](#metadata) | Fetch file metadata from server. |
| update | Post | id: file id. | json object | [Metadata](#metadata) | Partial update of a metadata. |
| search | Post | - | [Search data](https://github.com/Alvarios/nefts-go#options) | [Metadata List](https://github.com/Alvarios/nefts-go#results) | Retrieve a set of metadata from server. |

### Metadata (params)

An additional set of metadata to check on insertion. See the [Metadata](#metadata)
section for more details.

## Api mode

Api mode let you define pre-configured routes for your application. Just
provide a route with the correct url parameters when needed, and GCFS will
handle all the work for you.

> 💡 Tip : url parameter in go is given between curly braces.
 `/url/path/{url_parameter}/etc.`

## Methods mode

Methods mode is the default mode of gcfs. You can always use it, even when
api mode is set up.

Methods mode provides you some functions to interact directly with Couchbase.
Thus, you have full control of your api behavior.

### Insert

`fileId, err := gcfs_methods.Insert(fileMetadata)`

Insert metadata and returns the id of the generated tuple.

*arguments*

| Name | Type | Description |
| :--- | :--- | :--- |
| fileMetadata | [Metadata](#metadata) | Video metadata object. |

*return value*

| Name | Type | Description |
| :--- | :--- | :--- |
| fileId | string | Id pointing to the newly created document. |
| err | [Error object](#error-handling) | - |

### Delete

`err := gcfs_methods.Delete(fileId)`

Delete metadata from database. Only returns an error.

*arguments*

| Name | Type | Description |
| :--- | :--- | :--- |
| fileId | string | Id pointing to the document. |

*return value*

| Name | Type | Description |
| :--- | :--- | :--- |
| err | [Error object](#error-handling) | - |

### Get

`fileMetadata, err := gcfs_methods.Get(fileId)`

Retrieve metadata from database.

*arguments*

| Name | Type | Description |
| :--- | :--- | :--- |
| fileId | string | Id pointing to the document. |

*return value*

| Name | Type | Description |
| :--- | :--- | :--- |
| fileMetadata | [Metadata](#metadata) | Video metadata object. |
| err | [Error object](#error-handling) | - |

### Search

`queryResults, err := gcfs_methods.Search(start, end, options)`

Search function using [nefts](https://github.com/Alvarios/nefts-go).

*arguments*

| Name | Type | Description |
| :--- | :--- | :--- |
| start | int64 | Index based pagination for the first element to retrieve. |
| end | int64 | Index based pagination for the last element to retrieve. |
| options | [nefts options](https://github.com/Alvarios/nefts-go#options) | More details on [nefts](https://github.com/Alvarios/nefts-go#options) page. **(1)** |

**(1)** Config.Cluster and Config.Bucket options are overridden by gcfs configuration.

*return value*

| Name | Type | Description |
| :--- | :--- | :--- |
| queryResults | [nefts queryResults](https://github.com/Alvarios/nefts-go#results) | A list of video metadata. |
| err | [Error object](#error-handling) | - |

### Update

`timestamp, err := methods.Update(fileId, fileMetadata)`

Partial update of a document.

*arguments*

| Name | Type | Description |
| :--- | :--- | :--- |
| fileId | string | Id pointing to the document. |
| fileMetadata | [Partial Metadata](#partial-metadata) | Video partial metadata object. |

*return value*

| Name | Type | Description |
| :--- | :--- | :--- |
| timestamp | uint64 | Timestamp for the update. |
| err | [Error object](#error-handling) | - |

#### Partial metadata

GCFS provides an easy, declarative way to perform a partial update of your
metadata, using [gocb MutateIn](https://docs.couchbase.com/go-sdk/2.1/howtos/subdocument-operations.html#mutating) function.

```go
partialMetadata := gcfs_methods.UpdateSpec{
    Remove: []string{"key1", "key2", "key4.key4-2"},
    Upsert: map[string]interface{}{
        "key3": "new value",
        "key4": map[string]interface{}{
            "key4-1": "new value"
        }
    },
    Append: map[string]interface{}{
        "array-key": ["newValue1", "newValue2"]
    }
}
```

##### Remove

A list of keys to remove from the document. Use the dot syntax for nested keys.

##### Upsert

Update a list of values. Each value in the given map will replace the
equivalent one in the original document. If a path doesn't exist, it will be
created.

##### Append

Append a list of values to an array key.

> 💡 Tip : you can declare a non array value to append a single value. An
 array value will always be treaten with the `HasMultiple` flag set to true
 (see [gocb specs](https://docs.couchbase.com/go-sdk/2.1/howtos/subdocument-operations.html#array-append-and-prepend)).
 Add a nested array if you want to append an array as an array, and not a list
 of values.

### AutoProvide

`autoProvided, err = gcfs.metadata.AutoProvide(fileMetadata)`

Auto fill metadata with default values. More details in the [Metadata](#metadata) section.

### CheckIntegrity

`coreMetadata, err := gcfs.metadata.CheckIntegrity(fileMetadata)`

Check integrity of a dataset.

> 💡 Tip : coreMetadata represents the default and required metadata.
 More details in the [Metadata](#metadata) section.

## Metadata

### Core metadata

Metadata is a configurable JSON object that will link file to the
client. Some metadata is required for video creation.

Required metadata is provided below. You don't need to set them in config file, but you have to provide the ones
marked as required (you can provide the others, or let gcfs fill them by default).

| Key | Type | Default | Required | Description |
| :--- | :--- | :--- | :--- | :--- |
| url | url | - | true | Link to the distant file. |
| id | string | auto-generated | - | Id to use for the file. |
| general | Interface | {} | true | Critical metadata. |
| general.name | String | - | true | Name of the file. |
| general.size | Number | - | true | Byte size of the file. |
| general.creation_time | Date | Current Date<br/><br/>Nil if creation_time already present | - | Date of the first upload on distant server. |
| general.modification_time | Date | Current Date<br/><br/>Nil if no creation_time | - | Date of the last upload on distant server. |
| general.format | String | - | true | File format. |

### Custom metadata

Additionally, you can add a set of custom REQUIRED metadata. All data sent by
your client will be uploaded to database. This parameter just serve to enforce
some data to be present on insertion.

Custom metadata are set of key-value pairs. Key is the metadata key, and value
its required type. Use `interface` for undefined type.

```json
{
  "metadata": {
    "labels": "[]string",
    "stats": {
      "download_count": "uint64",
      "status": "string"
    }
  }
}
```

> 💡 Tip : types are string representing a go type. Custom types from your
 application are supported, however you'll need to prefix it with your
 package name, and run your code in go module mode.

## Error handling
 
 NEFTS returns a pointer to an error object adapted to web servers.
 
 | Key | Type | Description |
 | :--- | :--- | :--- |
 | Code | int | The http status of the error. |
 | Message | string | Describes the nature of the error. |
 
 ## Copyright
 2020 Alvarios - [MIT license](https://github.com/Alvarios/gcfs/blob/master/LICENSE)
