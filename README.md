# Gator
This tool is a blog aggregator that scrapes RSS feeds within the CLI to view.

## Setup
To setup gator, PostgreSQL and Go should be installed.

Ensure that you have a configuration file in your home directory. By default it is .gatorconfig.json. However if you choose to modify the file name/location, you may do so in config.go.

The config file's format should be a simple json in this format:
```json 
{
  "db_url":"postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name":"john"
}
```
**current_username** does not have to be filled at the start, but you will need to create a database and provide a connection string to it for **db_url**

sql schema migrations are provided in the sql/schema directory.
Using the command 
``goose postgres "postgres://postgres:postgres@localhost:5432/gator" up-to 4``
with these files will set up the necessary tables.

## Commands

### login \<username>
logs in a user with a specified name

### register \<username>
registers and logs in a new user

### reset
wipes all rows from the database

### users
lists all users, with the current logged in user shows with (current) next to it

### agg \<time interval>
gathers and saves posts for all rss feeds the currently logged in user is following

### addfeed \<feedname> \<url>
adds a feed with the given url and gives it the name provided

### feeds
lists all the feeds with their title, url and owners

### follow \<url>
follows the specified feed url for the currently logged in user

### following
lists all the feeds followed by the current user

### unfollow \<url>
unfollows a feed with the provided url

### browse \<optional limit>
lists the most **limit** most recent feeds that the current user is following. If no limit is provided, 2 is the default
