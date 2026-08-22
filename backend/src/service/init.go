package service

var GlobalService map[string]Service = map[string]Service{}

const (
	USER_SERVICE_NAME  = "UserService"
	EMAIL_SERVICE_NAME = "EmailService"
	GROUP_SERVICE_NAME = "GroupService"
	FILE_SERVICE_NAME  = "FileService"
)

func Init() {
	doRegister()
	doInject()
}

func register(service Service) {
	GlobalService[service.Name()] = service
}

func doRegister() {
	register(NewUserService(USER_SERVICE_NAME))
	register(NewEmailService(EMAIL_SERVICE_NAME))
	register(NewGroupService(GROUP_SERVICE_NAME))
	register(NewFileService(FILE_SERVICE_NAME))
}

func doInject() {
	for _, service := range GlobalService {
		service.Inject()
	}
}
