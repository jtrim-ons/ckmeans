# ckmeans

This is a fairly direct translation of the [ckmeans implementation in Simple Statistics](https://github.com/simple-statistics/simple-statistics/blob/master/src/ckmeans.js).

The Simple Statistics implementation is based on a version of the algorithm developed by Song, Wang and Zhong in a [series of papers](https://cran.r-project.org/web/packages/Ckmeans.1d.dp/citation.html) and released as an R package.

See also Bill Mill's [Python translation](https://github.com/llimllib/ckmeans).

## Input

The command line program `cmd/ckmeans/ckmeans.go` reads from standard input.  It expects input to look like this:

```
3
1.2
1.3
2.3
3.4
```

where the 3 on the first line is the number of clusters to create, and the remaining numbers are the array of data.

## Tests

I haven't added unit tests yet.  The `compare-to-*` directories provide Python and JS scripts with the same interface that should give the same results.  The `compare-to-simple-statistics` directory has a script to create a random dataset.

## Doubts

The Simple Statistics implementation has ` || 0` in a few places.  Bill Mill's Python implementation converts to `int` in the corresponding places.  I don't know why.  To add to my confusion, [the R implementation](https://github.com/cran/Ckmeans.1d.dp/blob/c0552fcdbf0c5aef12e7f36173eb09689ba9bf79/src/fill_log_linear.cpp#L58-L59) has the corresponding lines commented out and has slightly different code (which may produce different results? I don't know!)

As observed by Mill in a code comment, if all the values in the input are equal then the algorithm just returns one cluster.
