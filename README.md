# POE Marketwatch
Tool to analyze the Path of Exile market.

Currently supports:
- Estimate currency exchange prices
- Estimate unique prices
- Estimate divination card prices
- Report on most profitable divination cards

## Setup Instructions
Running this requires a postgresql database to connect to.

### Database Setup
Create the database: `$> createdb dbname`
Load the schema to the database: `$> psql dbname -f ./schema.sql`
Insert the inital page: `psql> INSERT INTO api_pages (id) VALUES ('<pageId>);`
Configure the application to connect to the database...

### Run the application
Run `$> go build` to product the executable
Run `$>DB_NAME=<database> DB_USER=<username> DB_HOST=<host> DB_PASSWORD=<password> LEAGUE=<league> ./poe-marketwatch`
After some time, enough data will accumulate to crunch some numbers...
If you want to exit, `Ctrl+C` will safely complete the current page and exit.

### Crunch some numbers
First, refresh the materialized views in postgres. This should always be done before checking data:
Run `psql> REFRESH MATERIALIZED VIEW mat_chaos_currency_rates; REFRESH MATERIALIZED VIEW mat_div_unique_pairs;`
See the profitable divination cards: `psql> SELECT * FROM divination_card_profits`
