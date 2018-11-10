# miniDB

[![Build Status](https://travis-ci.org/xumc/miniDB.svg?branch=master)](https://travis-ci.org/xumc/miniDB)

[![codecov](https://codecov.io/gh/xumc/miniDB/branch/master/graph/badge.svg)](https://codecov.io/gh/xumc/miniDB)


miniDB is s practice project which foucus on implementing key technical points in relational database. It's not a prodution ready project.

### roadmap
1. store records in to file.
2. DDL SQL.
3. DML SQL.
4. tcp interface.
5. master slave arch.

### Store

We plan to only support integer, bool and varchar types. In Store component, our goal is storing records into the file. Each file represents a table. the folder of the table files is database name.


### Query
