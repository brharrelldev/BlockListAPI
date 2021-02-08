# BlockListAPI

To compile and run do the following

```bash
git clone git@github.com:brharrelldev/BlockListAPI.git
```

Then build

```bash
make build-all
```


The build-all target will both regnerate a schema, and compile the go code

To run the cli type:

```bash
bin/blocklist-cli <flags>
```

There are 5 flags, to make it easier you can export them app as environment variables and run the root command with no args:

```bash
bin/blocklist-cli
```


The 5 environment variables you can export are:

```bash
export BLOCKLIST_USER=<blocklist_user>
export BLOCKLIST_PASS=<bloocklist_pass>
export DB_CACHE=<badger db cache for token storage>
export DB_PATH=<path to backend database>
export PORT=<port service runs on>
```

Docker Image can be



## Dependencies

Below is a list of depdencies

|Dependencies|Reason for usage|
|------------|----------------|
|badger    | kv store used to store token|
|gqlgen    | used to build graphql server in Go|
|mux       | for routing and auth middleware|
| mattn    | used as driver for sqlite3 db




