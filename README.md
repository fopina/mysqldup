# mysqldup
DUPlicate a database within the same mysql server

Every now and then a guy needs to clone a DB to make test a batch of schema changes before running that in the real DB.

Most documentation online implies that even if you want to clone the database inside the same MySQL instance, you still have to use `mysqldump`, create the new DB and import that in the new DB.

That sounds quite inefficient so after checking out this guy's [script](https://gist.github.com/csonuryilmaz/3f8f92fdad007f97986e61ad79aeb514), mysqldu~m~p was born.

## Instalation

Use `go get`:

```
go get github.com/fopina/mysqldup
```

Or download a pre-built binary from [releases](https://github.com/fopina/mysqldup/releases).

## Usage

```
$ ./mysqldup --help
Usage: ./mysqldup [OPTIONS] OLD_DB_NAME NEW_DB_NAME
  -f, --force                  drop NEW_DB_NAME if it already exists
      --help                   this screen
  -h, --hostname string        connect to host (default "127.0.0.1")
  -p, --password               use password when connecting to server (read from tty)
      --password-file string   use password when connecting to server (read from file)
  -P, --port int               port number to use for connection (default 3306)
  -u, --user string            user for login (default "root")
  -V, --version                output version information and exit
```

## Example

```
$ ./mysqldup -f -h remote.mysql.local -u superuser --password-file secret_password.txt prod prod_clone
281 tables to clone
[1 / 281] cloning super_secret_passwords
[2 / 281] cloning all_customer_pii
...
[280 / 281] cloning gdpr_fines
[281 / 281] cloning customer_closure_notice
```

This (real life) example took 5min while `mysqldump | mysql` took 15min.
