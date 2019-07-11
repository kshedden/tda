[![Build Status](https://travis-ci.com/kshedden/tda.svg?branch=master)](https://travis-ci.com/kshedden/tda)
[![Go Report Card](https://goreportcard.com/badge/github.com/kshedden/tda)](https://goreportcard.com/report/github.com/kshedden/tda)
[![codecov](https://codecov.io/gh/kshedden/tda/branch/master/graph/badge.svg)](https://codecov.io/gh/kshedden/tda)
[![GoDoc](https://godoc.org/github.com/kshedden/tda?status.png)](https://godoc.org/github.com/kshedden/tda)

tda : Topological data analysis in Golang
=========================================

The __tda__ package provides support for a few methods from topological data analysis.

Currently, methods for gridded data (images) are provided, including:

* Connected component labeling for binary images

* Object persistence analysis

* Landscape profiles

* Convex hull peels

See the [examples](http://github.com/kshedden/tda/tree/master/examples) directory for some use cases.

Below is a scatterplot of object birth/death times for
[this image](examples/images/HeLa_cells_stained_with_antibody_to_actin_(green)_,_vimentin_(red)_and_DNA_(blue).jpg),
with 90%, 95%, and 99% convex hull
peels plotted in red.  See
[examples/persistence](http://github.com/kshedden/tda/tree/master/examples/persistence)
for the code used to produce this plot.

![Image of persistence diagram](https://github.com/kshedden/tda/blob/master/examples/persistence/persistence.png)

Below is a landscape plot based on the same image used above.
See [examples/landscape](http://github.com/kshedden/tda/tree/master/examples/landscape)
for the code used to produce this plot.

![Image of landscape diagram](https://github.com/kshedden/tda/blob/master/examples/landscape/landscape.png)

Below is an animated PNG showing part of the image above being thresholded at a sequence of values.
See [examples/animate_threshold](http://github.com/kshedden/tda/tree/master/examples/animate_threshold)
for the code used to produce this plot.

![Animation of image thresholding](https://github.com/kshedden/tda/blob/master/examples/images/cells.apng)
