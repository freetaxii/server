# FreeTAXII/server #

[![Go Report Card](https://goreportcard.com/badge/github.com/freetaxii/server)](https://goreportcard.com/report/github.com/freetaxii/server)  [![GoDoc](https://godoc.org/github.com/freetaxii/server?status.png)](https://godoc.org/github.com/freetaxii/server)

The FreeTAXII Server is a TAXII 2 Server written in Go (golang)

## Version ##
0.3.1


## Installation ##

This package can be installed with the go get command:

```
go get github.com/freetaxii/server/cmd/freetaxii
cd github.com/freetaxii/server/cmd/freetaxii
go build freetaxii.go
```

## Dependencies ##

This software uses the following external libraries:
```
getopt
	go get github.com/pborman/getopt
	Copyright (c) 2017 Google Inc. All rights reserved.

gorilla/mux
	go get github.com/gorilla/mux
	Copyright (c) 2012 Rodrigo Moraes

libstix2
	go get github.com/freetaxii/libstix2
	Copyright (c) 2015-2018 Bret Jordan. All rights reserved. 

```

This software uses the following builtin libraries:
```
crypto/tls, encoding/json, errors, fmt, html/template, io/ioutil, log, net/http, os, strings, time
	Copyright 2009 The Go Authors
```

## Features ##

Below is a list of major features and which ones have been implemented:

- [x] TLS 1.2
- [x] Discovery Service
  - [x] Multiple Discovery Services
- [x] API Root Service
  - [x] Multiple API Roots Services
- [x] Endpoints
  - [x] Discovery
  - [x] API Root
  - [x] Collections
  - [x] Collection
  - [x] Objects (GET)
  - [ ] Objects (POST)
  - [x] Objects By ID
  - [x] Manifest
  - [x] Versions
  - [ ] Status
- [x] URL Filtering
  - [x] added_after
  - [x] limit
  - [x] match[id]
  - [x] match[type]
  - [x] match[version]
  - [x] match[spec_version]
- [x] Configuration
  - [x] From a file
  - [ ] From a database
- [x] Pagination
- [ ] Authentication


## License ##

This is free software, licensed under the Apache License, Version 2.0.


## Copyright ##

Copyright 2015-2018 Bret Jordan, All rights reserved.
