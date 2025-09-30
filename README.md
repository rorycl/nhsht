# nhsht

Generate a random salt and hash table for NHS numbers in England and
save the output to a parquet file.

version 0.0.4 : 30 September 2025 : add bloom filters

## About

This program generates a random 256 bit salt and corresponding hash
table for all NHS numbers in England as specified by the number range on
the [NHS Number Wikipedia article](https://en.wikipedia.org/wiki/NHS_number),
using the modulo11 calculation specified there. The resulting hash table
is saved to a compressed parquet file with bloom filters, suitable for
performing lookups with tools such as [duckdb](https://duckdb.org/).

The program makes use of goroutines for concurrently calculating 256
hashsums of the salt + nhs number. A default of `8 * runtime.NumCPU()`
is used, which may be overridden.

All ~300 million possible English NHS numbers will produce a zstd
compressed parquet file of about 20GB in size. By default the program
buffers to disk to keep RAM usage low. Using in-memory mode may require
significant swap space.

> [!CAUTION]
> Audits of this program prior to its use for sensitive data use cases
> should be undertaken and output files should be held securely.

## Post processing

Lookups with `duckdb` for sets of nhs numbers are reasonably fast. This
is despite records being interleaved by NHS number due to the use of
goroutines in the generation of hashes, which outputs records in a
non-deterministic order. For faster lookups by hash, it may be
worthwhile resorting the parquet file by hash, for example:

```sql
$ duckdb -c "
PRAGMA memory_limit='6GB'; 
PRAGMA threads=4;
COPY (
    SELECT
        *
    FROM
        'input.parquet'
    ORDER BY
        hash
) TO 'output.parquet' (FORMAT parquet);
"
```

## Stats

Tests were run on an old i7-7600U @ 2.80GHz machine (4 core, from 2017)
with 8GB RAM and an nvme disk. Results are indicative only.

```
300 million records generated in 7m57s.

$ time duckdb -c "select nhsno from 'data.parquet' order by random() limit 5;"
┌───────────┐
│   nhsno   │
│   int32   │
├───────────┤
│ xxxxxxx68 │
│ xxxxxxx74 │
│ xxxxxxx60 │
│ xxxxxxx57 │
│ xxxxxxx58 │
└───────────┘
real	0m9.320s

$ time duckdb -c "select nhsno from 'data.parquet' where nhsno IN (
    xxxxxxx68, xxxxxxx74, xxxxxxx60, xxxxxxx57, xxxxxxx58
  );"
┌───────────┐
│   nhsno   │
│   int32   │
├───────────┤
│ xxxxxxx57 │
│ xxxxxxx68 │
│ xxxxxxx60 │
│ xxxxxxx58 │
│ xxxxxxx74 │
└───────────┘
real	0m1.770s

$ time duckdb -c "select * from 'data.parquet' where nhsno IN (
    xxxxxxx68, xxxxxxx74, xxxxxxx60, xxxxxxx57, xxxxxxx58
 );"
┌───────────┬──────────────────────────────────────────────────────────────────┐
│   nhsno   │                               hash                               │
│   int32   │                             varchar                              │
├───────────┼──────────────────────────────────────────────────────────────────┤
│ xxxxxxx57 │ 188216ed34e78fbeaa8f051a81bf0ef24cbc5a9542aef5318b31cbb26cd2445a │
│ xxxxxxx68 │ 39057a683304a8637ce5f15a3602ada3d0e0c6fd5de7fb2ac078c88b9d76da3c │
│ xxxxxxx60 │ c7fff4d724a2b847ef9db4d4bb86c7c1ed64e0d2af72d4b81f2ec74811a3fe18 │
│ xxxxxxx58 │ 82ee7961a3debf695a4c46d17814ab256d92c9b0f2333346f946af3ceb78729d │
│ xxxxxxx74 │ da50ce5d246de97bf383e7f4d67e200d976e4861e497cbf2b7137f6d89d89585 │
└───────────┴──────────────────────────────────────────────────────────────────┘
real	0m23.151s

$ time duckdb -c "select * from 'data.parquet' where hash IN (
    '188216ed34e78fbeaa8f051a81bf0ef24cbc5a9542aef5318b31cbb26cd2445a',
    '39057a683304a8637ce5f15a3602ada3d0e0c6fd5de7fb2ac078c88b9d76da3c',
    'c7fff4d724a2b847ef9db4d4bb86c7c1ed64e0d2af72d4b81f2ec74811a3fe18',
    '82ee7961a3debf695a4c46d17814ab256d92c9b0f2333346f946af3ceb78729d',
    'da50ce5d246de97bf383e7f4d67e200d976e4861e497cbf2b7137f6d89d89585'
  );"
┌───────────┬──────────────────────────────────────────────────────────────────┐
│   nhsno   │                               hash                               │
│   int32   │                             varchar                              │
├───────────┼──────────────────────────────────────────────────────────────────┤
│ xxxxxxx57 │ 188216ed34e78fbeaa8f051a81bf0ef24cbc5a9542aef5318b31cbb26cd2445a │
│ xxxxxxx68 │ 39057a683304a8637ce5f15a3602ada3d0e0c6fd5de7fb2ac078c88b9d76da3c │
│ xxxxxxx60 │ c7fff4d724a2b847ef9db4d4bb86c7c1ed64e0d2af72d4b81f2ec74811a3fe18 │
│ xxxxxxx58 │ 82ee7961a3debf695a4c46d17814ab256d92c9b0f2333346f946af3ceb78729d │
│ xxxxxxx74 │ da50ce5d246de97bf383e7f4d67e200d976e4861e497cbf2b7137f6d89d89585 │
└───────────┴──────────────────────────────────────────────────────────────────┘
real	0m25.346s

```

## Usage

```
$ nhsht -h

Usage:
  nhsht 

NHS Number salted hash table generator.

version 0.0.4

This program: 

* generates a random 256bit salt and saves this, hex encoded, to the
  specified file.

* generates all of the currently possible NHS numbers for England and
  saves each number with a salted sha256 value in the specified parquet
  file which is zstd compressed. Not that a full hash table parquet file
  is likely to be around 20GB in size.

Use the "records" flag to avoid generating all possible 300 million NHS
numbers for England.

e.g.
  nhsht -s salt -p hashes.parquet
to generate only 20 records:
  nhsht -s salt -p hashes.parquet -r 20
to show progress:
  nhsht -s salt -p hashes.parquet -v

Application Options:
  -s, --saltfile=    file to save hex encoded salt (required)
  -p, --parquetfile= hash table parquet file path (required)
  -r, --records=     only generate this number of records
  -g, --goroutines=  number of goroutines (default 8 * numcpu)
  -m, --memory       use RAM memory, don't buffer to disk
  -v, --verbose      report progress of number generation

Help Options:
  -h, --help         Show this help message

```

## License

This project is licensed under the [MIT Licence](LICENCE) and provides
no guarantee of its fitness for purpose for any use.
