# dbtool

`dbtool` runs queries directly against a local database.

Ex:

```
go run cmd/dbtool/main.go --data ~/.sifnoded_prod/data/
data directory: /home/martin/.sifnoded_prod/data/
output file: /home/martin/dbtool.data
query: update_client.client_id='07-tendermint-41'
page: 1
per-page: 10
Getting transactions (page 1, perPage 10)...
results: 10 | total: 22636
Writing transactions to /home/martin/dbtool.data...
```

The CLI arguments are:

```
-data string
        data directory (default "~/.sifnoded/data")
  -out string
        output file (default "~/dbtool.data")
  -page int
        page number (default 1)
  -per-page int
        results per page (default 10)
  -query string
        query string (default "update_client.client_id='07-tendermint-41'")
```

