# FreeTAXII/freetaxii-server/html #

The FreeTAXII-Server is a TAXII 2 Server written in Go (golang)


## HTML Template Files Variables ##

Each of the HTML template files can be customized as needed per organizational 
requirements and supports the variables as defined below. The freetaxii.conf file 
allows template files to be defined globally and also redefined at a specific 
instance. For example, one can define an HTML template for all of the API roots
that are found on the server and can also override that template for a specific 
API root. This would allow each API root to have its own template and branding.


## HTML Configuration ##

```
"html" : {
    "enabled"           : true,
    "templatedir"       : "templates/html/",
    "templatefiles"     : {
        "discovery"     : "discoveryResource.html",
        "apiroot"       : "apirootResource.html",
        "collections"   : "collectionsResource.html",
        "collection"    : "collectionResource.html",
        "objects"       : "objectsResource.html",
        "versions"      : "versionsResource.html",
        "manifest"      : "manifestResource.html"
    }
}
```


## License ##

This is free software, licensed under the Apache License, Version 2.0.


## Copyright ##

Copyright 2015-2018 Bret Jordan, All rights reserved.