# gator

Gator is an RSS Feed aggregator written in Go.

## Installation and Config

You need PostgreSQL and Go installed on your computer to run gator. For information about the installation visit the [PostgreSQL](https://www.postgresql.org/) and [Go](https://go.dev/) websites.

You need to create a database for gator. You can use `psql` to connect to your PostgreSQL server.

```bash
sudo -u postgres psql
```

And then you can run the query to create the database. I named mine `gator`

```sql
CREATE DATABASE gator;
```

You need to create a `.gatorconfig.json` file in your home directory. This config file should contain the URL of your database (key `db_url`).

The URL may have different structure depending on your OS.

```json
{
    "db_url": "postgres://<user>:<passwd>@localhost:<port>/<db-name>?sslmode=disable"
}
```

After you have completed these steps, you can run the migrations in the `sql/schema` folder by using [`goose`](https://github.com/pressly/goose) (follow the instructions to install it).

Insert the connection string to your database: the same that you have in your config file, without the `?sslmode=disable` part.

```bash
cd sql/schema/
goose postgres "postgres://<user>:<passwd>@localhost:<port>/<db-name>" up
```

Finally, the CLI tool itself can be installed by running the following command in the root of your project:

```bash
go install ./...
```

## Usage

After installation and config you can run the commands by running `gator <command-name>` in your terminal.

List of available commands:

- `register <user-name>`: register new user
- `login <user-name>`: log user in (set as current user)
- `reset`: will delete all users, and their feeds, feed follows, posts, etc.
- `users`: list registered users
- `agg <duration>`: run aggregate loop and fetch last fetched feed with the specified timeout between fetches
- `addfeed <feed-name> <url>`: add new feed (will automatically follow it)
- `feeds`: list registered feeds
- `follow <url>`: follow specified feed by current user
- `following`: list feeds followed by current user
- `unfollow <url>`: unfollow specified feed by current user
- `browse [limit]`: list latest posts from the followed feeds with title, feed name, url, publication date

When you register a new user, they are automatically logged in as the current user. When you add a new feed, it is also automatically followed by the current user. A simple example flow:

```bash
gator register johndoe
gator addfeed ExampleFeed "https://example.com/rss"
gator agg 30s
```

Then in a separate terminal or after exiting the aggregation loop with `Ctrl+C`, you can run

```bash
gator browse 5
```
