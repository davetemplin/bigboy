# EXAMPLES

How to use **bigboy**

## Project Structure

```
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
```

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
LIMIT 5
```

**target.json**
```json
{
  "connection": "mysql_rfam"
}
```

Run `bigboy examples/basic` in terminal
Output is saved with `json` extension under `out/examples/basic` with current date filename.
Output:
```json
{"created":"2017-06-06 15:11:02","ncbi_id":"12892","scientific_name":"Potato spindle tuber viroid"}
{"created":"2017-06-06 15:11:07","ncbi_id":"12901","scientific_name":"Columnea latent viroid"}
{"created":"2017-06-06 15:11:12","ncbi_id":"53194","scientific_name":"Tomato apical stunt viroid-S"}
{"created":"2017-06-06 15:11:17","ncbi_id":"12885","scientific_name":"Tomato apical stunt viroid"}
{"created":"2017-06-06 15:11:21","ncbi_id":"32618","scientific_name":"Cucumber yellows virus"}
```

### Prefetch

Allows to run SQL extract in parallel

**extract.sql**
```sql
SELECT ncbi_id, scientific_name, created
FROM genome
WHERE ncbi_id IN (%s)
LIMIT 5
```

**prefetch.sql**
```sql
SELECT ncbi_id
FROM genome
WHERE created >= '2020-01-01'
```

**target.json**
```json
{
  "connection": "mysql_rfam",
  "prefetch": true
}
```

Run `bigboy examples/prefetch` in terminal
Output is saved with `json` extension under `out/examples/prefetch` with current date filename.
Output:
```json
{"created":"2020-09-11 00:21:31","ncbi_id":"11053","scientific_name":"Dengue virus 1"}
{"created":"2020-09-11 00:21:36","ncbi_id":"11053","scientific_name":"Dengue virus 1"}
{"created":"2020-09-11 00:22:17","ncbi_id":"11053","scientific_name":"Dengue virus 1"}
{"created":"2020-09-11 00:22:30","ncbi_id":"11053","scientific_name":"Dengue virus 1"}
{"created":"2020-09-11 00:24:18","ncbi_id":"11053","scientific_name":"Dengue virus 1"}
```

### Nest

Combines two query results into one
Note: the parent column has to be numeric. Currently can't perform nesting by string

#### nesting objects

**extract.sql**
```sql
SELECT rfam_acc, CONVERT(SUBSTRING(rfam_acc, 3, 5),UNSIGNED INTEGER) id
FROM family
WHERE rfam_acc IN ("RF00511", "RF00516", "RF01099", "RF01173", "RF01214")
```

**nest.sql**
```sql
SELECT CONVERT(SUBSTRING(rfam_acc, 3, 5),UNSIGNED INTEGER) _parent, author_id, name
FROM family_author
LEFT JOIN author USING(author_id)
WHERE CONVERT(SUBSTRING(rfam_acc, 3, 5),UNSIGNED INTEGER) IN (%s)
ORDER BY 1, 2;
```

**target.json**
```json
{
  "connection": "mysql_rfam",
  "nest": [
    {
      "connection": "mysql_rfam",
      "childKey": "authors",
      "parentKey": "id"
    }
  ]
}
```

Run `bigboy examples/nest` in terminal
Output is saved with `json` extension under `out/examples/nest` with current date filename.
Output:
```json
{"authors":[{"author_id":"39","name":"Moxon SJ"}],"id":"511","rfam_acc":"RF00511"}
{"authors":[{"author_id":"5","name":"Barrick JE"},{"author_id":"9","name":"Breaker RR"}],"id":"516","rfam_acc":"RF00516"}
{"authors":[{"author_id":"38","name":"Moore B"},{"author_id":"61","name":"Wilkinson A"}],"id":"1099","rfam_acc":"RF01099"}
{"authors":[{"author_id":"61","name":"Wilkinson A"},{"author_id":"62","name":"Eberhardt R"}],"id":"1173","rfam_acc":"RF01173"}
{"authors":[{"author_id":"11","name":"Burge SW"},{"author_id":"61","name":"Wilkinson A"}],"id":"1214","rfam_acc":"RF01214"}
```

#### nesting values

If you modify the **nest.sql** to use `_value` column like this:
```sql
SELECT CONVERT(SUBSTRING(rfam_acc, 3, 5),UNSIGNED INTEGER) _parent, name _value
FROM family_author
LEFT JOIN author USING(author_id)
WHERE CONVERT(SUBSTRING(rfam_acc, 3, 5),UNSIGNED INTEGER) IN (%s)
ORDER BY 1;
```

The output childKey property would have an array of strings rather than objects:
```json
{"authors":["Moxon SJ"],"id":"511","rfam_acc":"RF00511"}
{"authors":["Barrick JE","Breaker RR"],"id":"516","rfam_acc":"RF00516"}
{"authors":["Wilkinson A","Moore B"],"id":"1099","rfam_acc":"RF01099"}
{"authors":["Wilkinson A","Eberhardt R"],"id":"1173","rfam_acc":"RF01173"}
{"authors":["Burge SW","Wilkinson A"],"id":"1214","rfam_acc":"RF01214"}
```

## params

Customizing the query using CLI arguments or default params

**extract.sql**
```sql
SELECT ncbi_id, scientific_name, created, updated
FROM genome
WHERE updated BETWEEN ? AND ?
LIMIT 5;
```

**target.json**
```json
{
  "connection": "mysql_rfam",
  "timezone": "UTC",
  "params": [
    {
      "name": "updated",
      "type": "date"
    },
    {
      "name": "updated",
      "type": "date",
      "default": "today"
    }
  ]
}
```

Run `bigboy examples/params -1y` in terminal
Output is saved with `json` extension under `out/examples/params` with current date filename.
Output:
```json
{"created":"2017-06-06 15:11:02","ncbi_id":12892,"scientific_name":"Potato spindle tuber viroid","updated":"2020-04-23 11:46:08"}
{"created":"2017-06-06 15:11:07","ncbi_id":12901,"scientific_name":"Columnea latent viroid","updated":"2020-04-23 11:46:08"}
{"created":"2017-06-06 15:11:12","ncbi_id":53194,"scientific_name":"Tomato apical stunt viroid-S","updated":"2020-04-23 11:47:10"}
{"created":"2017-06-06 15:11:17","ncbi_id":12885,"scientific_name":"Tomato apical stunt viroid","updated":"2020-04-23 11:46:08"}
{"created":"2017-06-06 15:11:21","ncbi_id":32618,"scientific_name":"Cucumber yellows virus","updated":"2020-04-23 11:46:28"}
```

## split

Allows to split output into multiple files (currently by date)

**extract.sql**
```sql
SELECT upid, description, created
FROM genome
WHERE date BETWEEN '2020-01-01' AND '2020-01-07'
```

**target.json**
```json
{
  "connection": "mysql_rfam",
  "split": {
    "by": "date",
    "layout": "2006-01-02 15:04:05",
    "value": "created"
  }
}
```

Run `bigboy examples/split` in terminal
Output is saved in multiple files with `json` extension under `out/examples/split` each has records created in a corresponding date.

```
13548 rows written to 4 files in 1 seconds
```
```
out
└── examples
    └── split
        ├── 2018-03-12.json
        ├── 2018-03-13.json
        ├── 2018-03-19.json
        └── 2018-03-27.json
```

Output of the `out/examples/split/2018-03-27.json`:
```json
{"created":"2018-03-27 16:25:52","description":"Mesorhizobium huakuii 7653R genome","upid":"UP000027109"}
{"created":"2018-03-27 16:28:19","description":"Yersinia pseudotuberculosis IP 32953, complete genome.","upid":"UP000031850"}
```
