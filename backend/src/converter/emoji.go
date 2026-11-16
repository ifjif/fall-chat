package converter

import (
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/filter"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
)

func ToCreateEmojiPackCommand(req *request.CreateEmojiPack) *command.CreateEmojiPack {
	return &command.CreateEmojiPack{
		Name:      req.Name,
		CoverURL:  req.CoverURL,
		PriceType: req.PriceType,
		Price:     req.Price,
		Source:    req.Source,
	}
}

func ToListEmojiPackQuery(req *request.ListEmojiPack) *query.ListEmojiPack {
	qu := &query.ListEmojiPack{
		Next:    req.Next,
		SortKey: req.SortKey,
	}

	qu.Filter = &filter.EmojiPackFilter{
		PriceType: req.PriceType,
		Status:    req.Status,
	}

	return qu
}

func ToEmojiPackResult(m *model.EmojiPack) *response.EmojiPackResult {
	return &response.EmojiPackResult{
		ID:        m.PublicID,
		Name:      m.Name,
		Source:    m.GetSourceText(),
		PriceType: m.GetPriceTypeText(),
		Price:     m.Price,
		CoverURL:  m.CoverURL,
		Status:    m.GetStatusText(),
		OwnerID:   m.EncodeOwnerID(),
		CreatedAt: m.CreatedAt,
	}
}

func ToCreateEmojiItemCommand(req *request.CreateEmojiItem) *command.CreateEmojiItem {
	return &command.CreateEmojiItem{
		MediaType: req.MediaType,
		MediaURL:  req.MediaURL,
		FontName:  req.FontName,
		ThumbURL:  req.ThumbURL,
		Width:     req.Width,
		Height:    req.Height,
		FileSize:  req.FileSize,
		SortNO:    req.SortNO,
	}
}

func ToCreateEmojiItemCommands(req []*request.CreateEmojiItem) []*command.CreateEmojiItem {
	cmds := make([]*command.CreateEmojiItem, 0, len(req))
	for _, r := range req {
		cmds = append(cmds, ToCreateEmojiItemCommand(r))
	}

	return cmds
}

func ToPublishEmojiCommand(req *request.PublishEmoji) (*command.PublishEmoji, error) {
	packID, err := model.DecodeID(req.ID)
	if err != nil {
		return nil, err
	}
	return &command.PublishEmoji{
		PackID: packID,
		Items:  ToCreateEmojiItemCommands(req.Items),
	}, nil
}

func ToListEmojiItemQuery(req *request.ListEmojiItem) (*query.ListEmojiItem, error) {
	packID, err := model.DecodeID(req.PackID)
	if err != nil {
		return nil, err
	}
	return &query.ListEmojiItem{
		PackID:  packID,
		Next:    req.Next,
		SortKey: req.SortKey,
	}, nil
}

func ToEmojiItemResult(m *model.EmojiItem) *response.EmojiItemResult {
	return &response.EmojiItemResult{
		ID:        m.PublicID,
		PackID:    m.EncodePackID(),
		MediaType: m.GetMediaTypeText(),
		MediaURL:  m.MediaURL,
		FontName:  m.FontName,
		ThumbURL:  m.ThumbURL,
		Width:     m.Width,
		Height:    m.Height,
		FileSize:  m.FileSize,
		SortNO:    m.SortNO,
	}
}

func ToOwnListEmojiPackQuery(req *request.OwnListEmojiPack) *query.OwnListEmojiPack {
	qu := &query.OwnListEmojiPack{
		Next: req.Next,
	}

	ownFilter := &filter.EmojiOwnFilter{
		SourceType: req.SourceType,
	}
	packFilter := &filter.EmojiPackFilter{
		Status: req.Status,
	}
	qu.OwnFilter = ownFilter
	qu.PackFilter = packFilter

	return qu
}

func ToPublishEmojiVerifyCountQuery(req *request.PublishEmojiVerifyCount, userID uint64) (*query.PublishEmojiVerifyCount, error) {
	qu := &query.PublishEmojiVerifyCount{
		UserID:    userID,
		ItemCount: req.ItemCount,
	}
	packID, err := model.DecodeID(req.PackID)
	if err != nil {
		return nil, err
	}

	qu.PackID = packID

	return qu, nil
}
