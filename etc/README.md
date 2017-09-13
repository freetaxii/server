# FreeTAXII/freetaxii-server/etc/freetaxii.conf #

The FreeTAXII-Server is a TAXII 2 Server written in Go (golang)


## Configuration File Directives ##

The freetaxii.conf configuration file uses a JSON encoded configuration file with the following directives:


### Global Directives ###
- system
- logging
- discoveryservice
- apirootservice
- discoveryresources
- apirootresources
- collectionresources

### system directives ###

#### protocol ####
The protocol being used for this server. The only two options are http and https

#### listen ####
The IP address and port number that this server listens on. Example 127.0.0.1:8080

#### prefix ####
The installation prefix for the server. Example /opt/freetaxii

#### dbconfig ####
A boolean flag to tell the server if the server configuration comes from this text file or a database.

#### dbtype ####
The type of database that contains the server configuration information. Currently the only option is sqlite3

#### dbfile ####
The location of the database file that contains the server configuration

#### htmldir ####
The html template directory

#### tlskey ####
The name of the TLS private key that is located in etc/tls/

#### tlscrt ####
The name of the TLS public certificate that is located in etc/tls/

### logging directives ###

#### enabled ####
A boolean flag to enable logging

#### loglevel ####
An integer with a value of 1-5 that specifies the current logging level. The log levels are currently defined as:

- Log Level 1 = basic system logging information, sent to STDOUT unless Enabled = true then it is logged to a file
- Log Level 2 =
- Log Level 3 = detailed debugging information and code troubleshooting
- Log Level 4 =
- Log Level 5 = RAW packet/message decode and output

#### logfile ####
The location of the log file. Example: log/freetaxii.log

## License ##

This is free software, licensed under the Apache License, Version 2.0.


## Copyright ##

Copyright 2017 Bret Jordan, All rights reserved.
