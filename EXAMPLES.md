# EXAMPLES

How to use **bigboy**

## Project Structure

├── bigboy
├── bigboy.json
└── examples
    ├── basic
    │   ├── extract.sql
    │   └── target.json
    ├── nest
    │   ├── extract.sql
    │   ├── nest.sql
    │   └── target.json
    ├── params
    │   ├── extract.sql
    │   └── target.json
    ├── prefetch
    │   ├── extract.sql
    │   ├── prefetch.sql
    │   └── target.json
    └── split
        ├── extract.sql
        └── target.json

Note: *bigboy* is a binary executable. Can be downloaded from the list of [releases](https://github.com/igor-starostenko/bigboy/releases/latest)

All examples listed above are using the publically available [Rfam MySQL Database](https://docs.rfam.org/en/rfam-help-page/database.html)

**bigboy.json**
```json
{
  "connections": {
    "mysql_rfam": {
      "driver": "mysql",
      "server": "tcp(mysql-rfam-public.ebi.ac.uk:4497)",
      "database": "Rfam",
      "port": 4497,
      "user": "rfamro",
      "password": null
    }
  }
}
```

### Basic

Minimal target setup

**extract.sql**
```sql
SELECT ncbi_id, scientific_name, created
FROM genome
```

**target.json**
```json
{
  "connection": "mysql_rfam"
}
```

Run `bigboy examples/basic` in terminal
Output is saved with `json` extension under `out/examples/basic` with current date filename.
First 5 lines:
```json
{"created":"2017-06-06 15:11:02","ncbi_id":"12892","scientific_name":"Potato spindle tuber viroid"}
{"created":"2017-06-06 15:11:07","ncbi_id":"12901","scientific_name":"Columnea latent viroid"}
{"created":"2017-06-06 15:11:12","ncbi_id":"53194","scientific_name":"Tomato apical stunt viroid-S"}
{"created":"2017-06-06 15:11:17","ncbi_id":"12885","scientific_name":"Tomato apical stunt viroid"}
{"created":"2017-06-06 15:11:21","ncbi_id":"32618","scientific_name":"Cucumber yellows virus"}
...
```
