package instance

var usersInstance = &[]map[string]interface{}{}

func GetUsersInstance() []map[string]interface{} {
	return *usersInstance
}

func SetUsersInstance(instance []map[string]interface{}) {
	usersInstance = &instance
}

func SetUserDataByStudentId(studentId string, key string, value interface{}) {
	for _, user := range *usersInstance {
		if user["Student Id"] == studentId {
			user[key] = value
		}
	}
}

func SetUserDataByDiscordId(discordId string, key string, value interface{}) {
	for _, user := range *usersInstance {
		if user["Discord ID"] == discordId {
			user[key] = value
		}
	}
}
