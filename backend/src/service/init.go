package service

var GlobalService map[string]Service = map[string]Service{}

const (
	USER_SERVICE_NAME                = "UserService"
	EMAIL_SERVICE_NAME               = "EmailService"
	GROUP_SERVICE_NAME               = "GroupService"
	FILE_SERVICE_NAME                = "FileService"
	WS_SERVICE_NAME                  = "WsService"
	CONVERSATION_SERVICE_NAME        = "ConversationService"
	MESSAGE_SEQ_SERVICE_NAME         = "MessageSeqService"
	CONVERSATION_MEMBER_SERVICE_NAME = "ConversationMemberService"
	SEARCH_SERVICE_NAME              = "SearchService"
	FRIENDSHIP_SERVICE_NAME          = "FriendshipService"
	FRIENDSHIPREQUEST_SERVICE_NAME   = "FriendshipRequestService"
	GROUPREQUEST_SERVICE_NAME        = "GroupRequestService"
	GROUPMEMBER_SERVICE_NAME         = "GroupMemberService"
	IMMESSAGE_SERVICE_NAME           = "IMMessageService"
	EMOJI_SERVICE_NAME               = "EmojiService"
	RTC_SERVICE_NAME                 = "RTCService"
)

func Init() {
	doRegister()
	doInject()
}

func register(service Service) {
	GlobalService[service.Name()] = service
}

func doInject() {
	for _, service := range GlobalService {
		service.Inject()
	}
}

func doRegister() {
	register(NewUser(USER_SERVICE_NAME))
	register(NewEmail(EMAIL_SERVICE_NAME))
	register(NewGroup(GROUP_SERVICE_NAME))
	register(NewFile(FILE_SERVICE_NAME))
	register(NewWs(WS_SERVICE_NAME))
	register(NewConversation(CONVERSATION_SERVICE_NAME))
	register(NewMessageSeq(MESSAGE_SEQ_SERVICE_NAME))
	register(NewConversationMember(CONVERSATION_MEMBER_SERVICE_NAME))
	register(NewSearch(SEARCH_SERVICE_NAME))
	register(NewFriendship(FRIENDSHIP_SERVICE_NAME))
	register(NewFriendshipRequest(FRIENDSHIPREQUEST_SERVICE_NAME))
	register(NewGroupRequest(GROUPREQUEST_SERVICE_NAME))
	register(NewGroupMember(GROUPMEMBER_SERVICE_NAME))
	register(NewIMMessage(IMMESSAGE_SERVICE_NAME))
	register(NewEmoji(EMOJI_SERVICE_NAME))
	register(NewRTC(RTC_SERVICE_NAME))
}
