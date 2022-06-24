// done by:Vineesha, Sushmitha
package main

import (
	_ "github.com/lib/pq"
	ia "servermanagement.com/infraadmin/asset"
	dbs "servermanagement.com/infraadmin/database"
)

func main() {
	// db.Database()
	ia.HandleFunc()
	dbs.Connect()

}
