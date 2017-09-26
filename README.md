# FreeTAXII/freetaxii-server #

[![Go Report Card](https://goreportcard.com/badge/github.com/freetaxii/freetaxii-server)](https://goreportcard.com/report/github.com/freetaxii/freetaxii-server)

[![GoDoc](https://godoc.org/github.com/freetaxii/freetaxii-server?status.png)](https://godoc.org/github.com/freetaxii/freetaxii-server)

The FreeTAXII-Server is a TAXII 2 Server written in Go (golang)

## Version ##
0.0.1


## Installation ##

This package can be installed with the go get command:

```
go get github.com/freetaxii/freetaxii-server
cd github.com/freetaxii/freetaxii-server
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
	Copyright (c) 2017 Bret Jordan. All rights reserved. 

libtaxii2
	go get github.com/freetaxii/libtaxii2
	Copyright (c) 2017 Bret Jordan. All rights reserved.
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
- [x] API Root Service
- [x] Multiple API Roots
- [x] Collections
- [x] Collection
- [x] Objects (GET)
- [ ] Objects (POST)
- [ ] Objects By ID
- [ ] Manifests
- [ ] Status
- [ ] URL Filtering
- [x] Configuration from a file
- [ ] Configuration from a database
- [ ] Pagination
- [ ] Authentication


## License ##

This is free software, licensed under the Apache License, Version 2.0.


## Copyright ##

Copyright 2017 Bret Jordan, All rights reserved.
