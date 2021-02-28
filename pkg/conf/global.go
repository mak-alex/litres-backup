package conf

var (
	// DBConn - Database connection object
	//DBConn *gorm.DB

	// GlobalConfig - Config global parameters
	GlobalConfig Conf
	FilterBook   BookFilter
)

// Get - ...
func Get() Conf {
	return GlobalConfig
}
