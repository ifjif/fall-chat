package converter

import (
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/query"
)

func ToConvMembListBrief4ConvQuery(req *request.ConvMembListBrief4Conv) *query.ConvMembListBrief4Conv {
	return &query.ConvMembListBrief4Conv{
		ConvID: req.ConvID,
	}
}

func ToConvMembListBriefBatchQuery(req *request.ConvMembListBriefBatch) *query.ConvMembListBriefBatch {
	return &query.ConvMembListBriefBatch{
		ConvID:    req.ConvID,
		MemberIDs: req.MemberIDs,
	}
}
