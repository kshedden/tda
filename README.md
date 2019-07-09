[![Build Status](https://travis-ci.com/kshedden/tda.svg?branch=master)](https://travis-ci.com/kshedden/tda)
[![Go Report Card](https://goreportcard.com/badge/github.com/kshedden/tda)](https://goreportcard.com/report/github.com/kshedden/tda)
[![codecov](https://codecov.io/gh/kshedden/tda/branch/master/graph/badge.svg)](https://codecov.io/gh/kshedden/tda)
[![GoDoc](https://godoc.org/github.com/kshedden/tda?status.png)](https://godoc.org/github.com/kshedden/tda)

tda : Topological data analysis in Golang
=========================================

The tda package provides support for a few methods from topological data analysis.

Currently, methods for gridded data (images) are supported, including:

* Connected component labeling for binary images

* Object persistence diagrams

* Feature generation based on convex peels

See the [examples] directory for some use cases.

![Image of persistence diagram](github.com/kshedden/tda/examples/persistence/persistence.png)