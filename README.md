## Installation
go install github.com/Zoisit/blog-aggregator/

## Set-up
Create a config file in your home directory with the content:
{
  "db_url": "postgres://<username>:<opt. password>@localhost:5432/<database>?sslmode=disable"
}
under the name '.gatorconfig.json'

Replace the values with your database connection string.

## Usage

Start the aggregator:

```bash
gator agg 30s
```

There are a few other commands you'll need as well:
- gator register <username> - Add new user and login as that user
- gator login <username> - Log in as the given user, if that user is registeres
- gatpr addfeed <url> - Add new feed to database and follow that feed
- gator users - List all users
- gator feeds - List all feeds
- gator browse <limit> - List the <limit> most recent posts
- gator following - List all feeds the logged-in user follows
- gator follow <url> - Follow a feed that already exists in the database
- gator unfollow <url> - Unfollow a feed that already exists in the database
- gator agg <time> - Gets new posts for every feed starting from the oldest. <time> is the duration between getting the different feeds
