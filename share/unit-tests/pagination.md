# FreeTAXII/libstix2/db/tests/pagination.md #

This documents the various unit tests that should be performed on your implementation of TAXII pagination 

## Pagination ##

These tests will verify that the server is sending the right response based on the client request. 

### Setup ###

For this test the server will have 10 records in the collection and the server is setup to only send 5 records at a time.

### Test 1 ###

Client requests objects, but does not indicate any pagination or subset of the data
1. 	Server responds with all records
2.	Server responds with objects 1-5/10

### Test 2 ###
Client requests objects 1-2
1. 	Server responds with objects 1-2/10

### Test 3 ###

Client requests objects 4-5
1. 	Server responds with objects 4-5/10

### Test 4 ###

Client requests objects 4-10
1. 	Server responds with objects 4-8/10

### Test 5 ###

Client requests objects 1-20
1. 	Server responds with an error


## License ##

This is free software, licensed under the Apache License, Version 2.0.


## Copyright ##

Copyright 2016 Bret Jordan, All rights reserved.