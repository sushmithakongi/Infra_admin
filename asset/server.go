/*Server Management using GOlang.
Modified on 24/6/2022
modified by Sushmitha B Kongi
This file contains funtions related to add asset,list all asset as infra admin(role).*/

package asset

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	dbs "servermanagement.com/infraadmin/database"
)

type asset struct {
	Asset_Id       int    `json:"Asset_Id"`
	Manufacturer   string `json:"Manufacturer"`
	BMC_IP         string `json:"BMC_IP"`
	BMC_User       string `json:"BMC_USER"`
	BMC_Password   string `json:"BMC_Password"`
	Asset_location string `json:"Asset_location"`
	Created_on     string `json:"Created_on"`
	Created_by     string `json:"Created_by"`
	OS_IP          string `json:"OS_IP"`
	OS_User        int    `json:"OS_User"`
	OS_Password    string `json:"OS_Password"`
	Purpose        string `json:"Purpose"`
	Cluster_Id     string `json:"Cluster_ID"`
}

// type platform_profile struct {
// 	Asset_Id   int    `json:"Asset_Id"`
// 	RAM_gb     int    `json:"RAM_gb"`
// 	STORAGE_gb int    `json:"STORAGE_gb"`
// 	Cpu_model  string `json:"Cpu_model"`
// }

var db = dbs.Connect() //database connection using function

//------------------------------------------------add asset(creating asset)---------------------------------------------------------------------

func addAsset(write http.ResponseWriter, request *http.Request) {
	var assets asset
	err := json.NewDecoder(request.Body).Decode(&assets)
	if err != nil {
		json.NewEncoder(write).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": err})
		return
	}
	addStatement := `INSERT INTO asset (Asset_Id,Manufacturer, BMC_IP, BMC_User, BMC_Password, Asset_location,Created_on,Created_by,OS_IP,OS_User,OS_Password,Purpose,Cluster_Id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
	_, err = db.Exec(addStatement, assets.Asset_Id, assets.Manufacturer, assets.BMC_IP, assets.BMC_User, assets.BMC_Password, assets.Asset_location, assets.Created_on, assets.Created_by, assets.OS_IP, assets.OS_User, assets.OS_Password, assets.Purpose, assets.Cluster_Id)
	if err != nil {
		json.NewEncoder(write).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": err})
		return
	}

	json.NewEncoder(write).Encode(map[string]interface{}{"Status Code": "200 OK", "Message": "Recorded sucessfully"})
}

//--------------------------------------------Dashboard:number of reserved and vacant asset---------------------------------------------------------------------

func GetDashboard1(write http.ResponseWriter, request *http.Request) {
	var count, total int
	err := db.QueryRow("SELECT count(*) from asset WHERE reserved is not null").Scan(&count) // exporting table
	if err != nil {
		json.NewEncoder(write).Encode(map[string]interface{}{"Status Code": "400 Bad Request", "Message": err})
	}
	err1 := db.QueryRow("SELECT COUNT(*) FROM asset").Scan(&total) // exporting table
	if err1 != nil {
		json.NewEncoder(write).Encode(map[string]interface{}{"Status Code": "400 Bad Request", "Message": err1})
	}
	vacant := total - count
	json.NewEncoder(write).Encode(map[string]interface{}{"Number of asset reserved": count, "Number of asset vacant": vacant, "Status Code": "200 OK", "Message": "Updated Statistics"})

}

//--------------------------------------------Dashboard:number of reserved and vacant asset in cluster---------------------------------------------------------------------

func GetDashboard2(write http.ResponseWriter, request *http.Request) {
	var count, total int
	err := db.QueryRow("SELECT count(distinct cluster_id) from asset where reserved='t' group by cluster_id").Scan(&count) // exporting table
	if err != nil {
		log.Printf("Failed to open table: %v \n", err)
	}
	err1 := db.QueryRow("select count(distinct cluster_id)from asset").Scan(&total) // exporting table
	if err1 != nil {
		log.Printf("Failed to open table: %v \n", err1)
	}
	json.NewEncoder(write).Encode(map[string]interface{}{"Number of Clusters reserved": count, "Total number of clusters": total, "Status Code": "200 OK", "Message": "Updated Statistics"})

}

//--------------------------------------------Dashboard:number of reserved by location---------------------------------------------------------------------
func GetDashboard3(write http.ResponseWriter, request *http.Request) {
	var reserve, vacant int
	var asset_location string
	rows, err := db.Query(" SELECT asset_location, COUNT(CASE WHEN reserved='t' THEN 1 ELSE NULL END)AS reserved,  COUNT(CASE WHEN reserved='f' or reserved is null THEN 1 ELSE NULL END)AS vacant FROM asset group by asset_location")
	if err != nil {
		log.Fatal(err)
		json.NewEncoder(write).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": err})
	}
	fmt.Printf("Location|Reserved|Vacant\n")
	for rows.Next() {
		err := rows.Scan(&asset_location, &reserve, &vacant)
		if err != nil {
			log.Fatal(err)
			json.NewEncoder(write).Encode(map[string]interface{}{"status": "400 Bad Request", "Message": err})
		}
		fmt.Printf("%v : %v: %v\n ", asset_location, reserve, vacant)
		lct := map[string]interface{}{"Location": asset_location, "Reserved": reserve, "Vacant": vacant}
		json.NewEncoder(write).Encode(map[string]interface{}{"Dashboard": lct, "Status Code": "200 OK", "Message": "Updated Statistics"})

	}
}

//-------------------------------------------------Handle function----------------------------------------------------------
func HandleFunc() {
	http.HandleFunc("/add_asset", addAsset)
	http.HandleFunc("/dashboard1", GetDashboard1)
	http.HandleFunc("/dashboard2", GetDashboard2)
	http.HandleFunc("/dashboard3", GetDashboard3)
	fmt.Println("Assets are hosted \n Add_Asset API : http://localhost:5002/add_asset \n Dasboard1 API :  http://localhost:5002/dashboard1 \n Dasboard2 API :  http://localhost:5002/dashboard2 \n Dasboard3 API :  http://localhost:5002/dashboard3 \n ")
	log.Fatal(http.ListenAndServe(":5002", nil))
}
