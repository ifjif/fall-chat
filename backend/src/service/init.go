package service

var GlobalService map[string]Service = map[string]Service{}

const (
	USER_SERVICE_NAME                = "UserService"
	EMAIL_SERVICE_NAME               = "EmailService"
	GROUP_SERVICE_NAME               = "GroupService"
	FILE_SERVICE_NAME                = "FileService"
	SOCIAL_SERVICE_NAME              = "SocialService"
	WS_SERVICE_NAME                  = "WsService"
	CONVERSATION_SERVICE_NAME        = "ConversationService"
	MESSAGE_SEQ_SERVICE_NAME         = "MessageSeqService"
	CONVERSATION_MEMBER_SERVICE_NAME = "ConversationMemberService"
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
	register(NewSocialService(SOCIAL_SERVICE_NAME))
	register(NewWsService(WS_SERVICE_NAME))
	register(NewConversationService(CONVERSATION_SERVICE_NAME))
	register(NewMessageSeqService(MESSAGE_SEQ_SERVICE_NAME))
	register(NewConversationMemberService(CONVERSATION_MEMBER_SERVICE_NAME))
}

func doInject() {
	for _, service := range GlobalService {
		service.Inject()
	}
}
