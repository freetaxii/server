// Copyright 2017 Bret Jordan, All rights reserved.
//
// Use of this source code is governed by an Apache 2.0 license
// that can be found in the LICENSE file in the root of the source
// tree.

package server

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// --------------------------------------------------
// Load Services from Database
// --------------------------------------------------

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
					tbl_ApiRoots.sHost || tbl_ApiRoots.sPath AS fullPath
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
		// Only one of the string values in the database can be null, so we only
		// need to add special handing for it
		var description sql.NullString
		var resource DiscoveryResourceType
		resource.Api_roots = make([]string, 0)

		// If there is a record returned, then the directory service has been enabled
		this.DiscoveryService.Enabled = true

		err = rows.Scan(&resource.Id, &resource.Location, &resource.Title, &description, &resource.Contact, &resource.Default)

		if description.Valid {
			resource.Description = description.String
		}

		if err != nil {
			log.Printf("error reading from database, %v", err)
		}

		sqlstmt2 := `SELECT 
						sHost || sPath as fullPath
					FROM tbl_ApiRoots
					WHERE bEnabled = 1 AND iDiscoveryID = ?`

		// Pass in the record ID from the record we are working on
		rows2, err := db.Query(sqlstmt2, resource.Id)
		if err != nil {
			log.Printf("error running query, %v", err)
		}
		defer rows2.Close()

		for rows2.Next() {
			var fullUrl string
			err = rows2.Scan(&fullUrl)
			resource.Api_roots = append(resource.Api_roots, fullUrl)
		}

		// Add resources to object
		this.DiscoveryService.Resources = append(this.DiscoveryService.Resources, resource)
	}
}
