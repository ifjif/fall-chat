package converter

import (
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/query"
)

func ToConvMembListBrief4ConvQuery(req *request.ConvMembListBrief4ConvRequest) *query.ConvMembListBrief4ConvQuery {
	return &query.ConvMembListBrief4ConvQuery{
		ConvID: req.ConvID,
	}
}

func ToConvMembListBriefBatchQuery(req *request.ConvMembListBriefBatchRequest) *query.ConvMembListBriefBatchQuery {
	return &query.ConvMembListBriefBatchQuery{
		ConvID:    req.ConvID,
		MemberIDs: req.MemberIDs,
	}
}
