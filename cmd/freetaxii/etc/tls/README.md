# FreeTAXII/freetaxii-server/etc/tls #

The FreeTAXII-Server is a TAXII 2 Server written in Go (golang)


## Installation ##

To generate an RSA key and self-signed certificate for your server you can do the following in the freetaxii/etc/tls directory:

```
openssl req -x509 -nodes -newkey rsa:4096 -keyout server.rsa.key -out server.rsa.crt -days 3650
chmod 600 server.rsa.key
ln -sf server.rsa.key server.key
ln -sf server.rsa.crt server.crt

```

## Configuraiton ##

The following header was added to each of the handlers. This was done per RFC 6797 (https://tools.ietf.org/html/rfc6797)
Additional information can be found here: (https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Strict-Transport-Security)

```
w.Header().Add("Strict-Transport-Security", "max-age=86400; includeSubDomains")
```


## License ##

This is free software, licensed under the Apache License, Version 2.0.


## Copyright ##

Copyright 2017 Bret Jordan, All rights reserved.
