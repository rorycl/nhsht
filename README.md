# nhsht

Generate a random salt and hash table for NHS numbers in England.

version 0.0.1 : 29 September 2025 : First release

## About

This program generates a random salt and hash table for all NHS numbers
in England as specified by the number range on the [NHS Number Wikipedia
article](https://en.wikipedia.org/wiki/NHS_number), using the modulo11
calculation specified there. The resulting hash table is saved to a
parquet file, suitable for performing lookups with tools such as
[duckdb](https://duckdb.org/).

The program makes use of goroutines for concurrently calculating 256
hashsums of the salt + nhs number.

Presently the parquet file is only written on completion of the
calculations, requiring about 60GB of memory for generating all ~300
million English NHS numbers, and a similar amount of disk space.

### Usage

```
$ nhsht -h

Usage:
  nhsht 

NHS Number salted hash table generator.

version 0.0.1

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
  -v, --verbose      report progress of number generation

Help Options:
  -h, --help         Show this help message

```

## License

This project is licensed under the [MIT Licence](LICENCE).
