package converter

import (
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
)

func ToConvMembListBrief4ConvQuery(req *request.ConvMembListBrief4Conv) (*query.ConvMembListBrief4Conv, error) {
	convID, err := model.DecodeID(req.ConvID)
	if err != nil {
		return nil, err
	}
	return &query.ConvMembListBrief4Conv{
		ConvID: convID,
	}, nil
}

func ToConvMembListBriefBatchQuery(req *request.ConvMembListBriefBatch) (*query.ConvMembListBriefBatch, error) {
	convID, err := model.DecodeID(req.ConvID)
	if err != nil {
		return nil, err
	}

	memberIDs := make([]int64, 0, len(req.MemberIDs))

	for _, pID := range req.MemberIDs {
		id, err := model.DecodeID(pID)
		if err != nil {
			return nil, err
		}

		memberIDs = append(memberIDs, id)
	}
	return &query.ConvMembListBriefBatch{
		ConvID:     convID,
		MemberIDs:  memberIDs,
		MemberPIDs: req.MemberIDs,
	}, nil
}

func ToConvMembCommitReadCommand(req *request.ConvMembCommitRead, userID int64) (*command.ConvMembCommitRead, error) {
	cmd := &command.ConvMembCommitRead{
		UserID: userID,
		MsgSeq: req.MsgSeq,
	}

	msgID, err := model.DecodeID(req.MsgID)
	if err != nil {
		return nil, err
	}

	convID, err := model.DecodeID(req.ConvID)
	if err != nil {
		return nil, err
	}

	cmd.ConvID = convID
	cmd.MsgID = msgID

	return cmd, nil
}
