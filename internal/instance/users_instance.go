package instance

var usersInstance = &[]map[string]interface{}{}

func GetUsersInstance() []map[string]interface{} {
	return *usersInstance
}

func SetUsersInstance(instance []map[string]interface{}) {
	usersInstance = &instance
}
