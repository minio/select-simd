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

Single core

```
BenchmarkSelectCsv-8                   5         271120281 ns/op        4558.57 MB/s           0 B/op          0 allocs/op
```

Multi-core
```
BenchmarkParallel_2cpus_256KB-8            10         125521689 ns/op         9846.27 MB/s     133158 B/op         21 allocs/op
BenchmarkParallel_3cpus_256KB-8            20          89044577 ns/op        13879.79 MB/s     199611 B/op         30 allocs/op
BenchmarkParallel_4cpus_256KB-8            20          80467654 ns/op        15359.22 MB/s     265046 B/op         36 allocs/op
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
