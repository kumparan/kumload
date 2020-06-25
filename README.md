#kumload
kumload is database migration tool.
We can easily migrate database with defined schema in yaml file.

- Easy to migrate database
- Able to stop or continue migration
- Zero code

#### Limitation
- Now we only support migrate in same database server

## Installation
- Golang
```shell script
go install github.com/kumparan/kumload/cmd/kumload
```
- Executable
```shell script
download from releases, and execute in your machine
```

## Usage
```shell script
kumload --config config.yml
```

```shell script
CLI for database migration

Usage:
  kumload [flags]

Flags:
      --config string   config file path (default "config.yml")
  -h, --help            help for kumload
      --script string   script file path, this is optional
  -v, --version         version for kumload
```

## Configuration
- config.yml
```yaml
level: info
batch:
  interval: 0
  size: 2
database:
  host: 100.100.100.100:212
  username: root
  password: secret
  source:
    name: good_service
    table: foods
    primary_key: id
    order_key: id
  target:
    name: bad_service
    table: drinks
    primary_key: id
mappings:
  id: id
  name: name
  username: username
  email: email
```
## Configuration params
- Batch
```yaml
size:
  size per batch query
interval:
  delay duration per bacth query
```
- Database
```yaml
host:
  database host
username:
  database username
password:
  database password
source:
  source migration
    db:
      source database name
    table:
      source table
    primary_key:
      source field name that used as primary key
target:
  target migration
    db:
      target database name
    table:
      target table
    primary_key:
      target field name that used as primary key
```
- mappings
```yaml
mappings is optional if script provided it will use script as default sql query

<field_from_source>: <field_from_target>
```

## Custom Script
- script.sql
```sql
insert into sessions(
    "id",
    "name",
    "username",
    "email",
) select
    id,
    name,
    username,
    email,
from
    good_service.foods
where
    id not in (
        select
            id
        from drinks
    )
order by created_at asc limit ?
```
