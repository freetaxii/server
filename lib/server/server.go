// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package server

import (
	"database/sql"
	"github.com/freetaxii/freetaxii-server/lib/systemconfig"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// ----------------------------------------------------------------------
// Define Server Type
// ----------------------------------------------------------------------

type ServerType struct {
	SysConfig        *systemconfig.SystemConfigType
	DiscoveryService DiscoveryServiceType
}

type DiscoveryServiceType struct {
	Enabled   bool
	Resources []DiscoveryResourceType
}

type DiscoveryResourceType struct {
	Title       string
	Description string
	Location    string
	Contact     string
	Default     string
	Api_roots   []string
}

// --------------------------------------------------
// Load Services from Database
// --------------------------------------------------

// Init - This method will setup the struct with all of the subobjects need
func (this *ServerType) Init() {
	if this.DiscoveryService.Resources == nil {
		a := make([]DiscoveryResourceType, 0)
		this.DiscoveryService.Resources = a
	}
}

// LoadDiscoveryService - This method will load the configuration for the Discovery
// service and store it in the DiscoveryServiceType. All information will be pulled
// from the SQL database.
func (this *ServerType) LoadDiscoveryServicesConfig() {

	if this.SysConfig.Logging.LogLevel >= 1 {
		log.Println("DEBUG-1: Loading Discovery Services Configuration")
	}

	// Open connection to database
	filename := this.SysConfig.System.DbFileFullPath
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatalf("Unable to open file %s due to error %v", filename, err)
	}
	defer db.Close()

	// Read in services for the discovery server.
	sqlstmt := `SELECT
					tbl_DiscoveryServices.aID,
					tbl_DiscoveryServices.sLocation,
					tbl_DiscoveryServices.sTitle,
					tbl_DiscoveryServices.sDescription,
					tbl_DiscoveryServices.sContact,
					tbl_ApiRoots.sURL
				FROM
					tbl_DiscoveryServices
				JOIN tbl_ApiRoots
					ON tbl_DiscoveryServices."iDefault" = tbl_ApiRoots.aID
				WHERE
					tbl_DiscoveryServices.bEnabled = 1`
	rows, err := db.Query(sqlstmt)
	if err != nil {
		log.Printf("error running query, %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var location string
		var title string
		var description string
		var contact string
		var default_apiroot string

		// If there is a record returned, then the directory service has been enabled
		this.DiscoveryService.Enabled = true

		err = rows.Scan(&id, &location, &title, &description, &contact, &default_apiroot)

		if err != nil {
			log.Printf("error reading from database, %v", err)
		}

		var resource DiscoveryResourceType
		resource.Location = location
		resource.Title = title
		resource.Description = description
		resource.Contact = contact
		resource.Default = default_apiroot
		resource.Api_roots = make([]string, 0)

		sqlstmt2 := `SELECT sURL
					FROM tbl_ApiRoots
					WHERE bEnabled = 1 AND iDiscoveryID = ?`

		// Pass in the record ID from the record we are working on
		rows2, err := db.Query(sqlstmt2, id)
		if err != nil {
			log.Printf("error running query, %v", err)
		}
		defer rows2.Close()

		for rows2.Next() {
			var url string
			err = rows2.Scan(&url)
			resource.Api_roots = append(resource.Api_roots, url)
		}

		// Add resources to object
		this.DiscoveryService.Resources = append(this.DiscoveryService.Resources, resource)
	}
}
