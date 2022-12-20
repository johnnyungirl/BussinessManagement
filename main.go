package main

import (
	"BussinessManagement/model"
	"BussinessManagement/route"
)

// @title           Bussiness Management APIs
// @version         1.0
// @description     This is document for APIs.
// @termsOfService  http://swagger.io/terms/

// @contact.name   ITD Support
// @contact.url    http://www.itd.io/support
// @contact.email  support@itd.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey  Bearer
// @in                          header
// @name                        Authorization
// @description					Description for what is this security definition being used

func main() {

	db, _ := model.DBConnection()
	route.SetupRoutes(db)
}
