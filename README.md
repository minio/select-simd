# select-simd

This is a package for high performance execution of `S3 Select` statements on large CSV and Parquet objects. It is meant to be integrated into the Minio Object Storage server.

By leveraging SIMD instructions in combination with zero copy behavior, early indications show that parsing and selecting speeds of several GB/sec are possible on a single core on a modern CPU.

Support for JSON objects is planned as well.

## Works in progress

This package is in early development and does not yet cover a wide range of operations.

## Design Goals

`select-simd` has been designed with the following goals in mind:

- zero copy behavior
- low memory footprint
- leverage SIMD instructions
- support for streaming (chunking)
- support large objects (> 4 GB)
- pure Golang (with Golang assembly)

## High Performance

```
BenchmarkSelectCsv-8            2         531865224 ns/op        6971.24 MB/s           0 B/op          0 allocs/op
```

## Architecture

To be described.

```
         CSV           JSON          Parquet

    +-----------+  +-----------+  +------------
    |           |  |           |  |           |
    |  Parsing  |  |  Parsing  |  |  Loading  |
    |           |  |           |  |           |
    +-----------+  +-----------+  +------------

    +-----------------------------------------+
    |                                         |
    |          Evaluation  ("where")          |
    |                                         |
    +-----------------------------------------+

    +-----------------------------------------+
    |                                         |
    |          Processing ("select")          |
    |                                         |
    +-----------------------------------------+
```

## License

Released under the Apache License v2.0. You can find the complete text in the file LICENSE.

## Contributing

Contributions are welcome, please send PRs for any enhancements.
