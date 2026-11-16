package dao

import (
	"context"
	"errors"
	"zzc/fall-chat/src/cursor"
	"zzc/fall-chat/src/data/executor"
	"zzc/fall-chat/src/model"

	"gorm.io/gorm"
)

type IMMessageDao interface {
	CreateOne(ctx context.Context, m *model.IMMessage) error
	List(ctx context.Context, convID uint64, sort map[string]any, sortDef cursor.SortDefinition, limit int) ([]*model.IMMessage, error)
	PatchList(ctx context.Context, convID uint64, minSeq, maxSeq int64) ([]*model.IMMessage, error)
	FindSeq(ctx context.Context, msgID, convID int64) (int64, error)
}

type imMessageDao struct {
	db executor.Executor
}

func NewIMMessageDao(db executor.Executor) IMMessageDao {
	return &imMessageDao{
		db: db,
	}
}

func (dao *imMessageDao) CreateOne(ctx context.Context, m *model.IMMessage) error {
	return dao.db.Create(m).Error
}

// id,
// public_id,
// conv_id,
// content,
// status,
// send_time,
// msg_type,
// sender_id,
// seq
func (dao *imMessageDao) List(ctx context.Context, convID uint64, sort map[string]any, sortDef cursor.SortDefinition, limit int) ([]*model.IMMessage, error) {
	var rows []*model.IMMessage
	orderby := sortDef.OrderBy(nil)

	qc := dao.db.Model(model.IMMessage{}).
		Where("conv_id = ?", convID).
		Limit(limit).
		Order(orderby).
		Select("id,public_id,conv_id,content,status,send_time,msg_type,sender_id,seq")

	whereClause, args, err := sortDef.Where(sort, nil)
	if err != nil {
		return nil, err
	}

	if args != nil {
		qc.Where(whereClause, args...)
	}

	err = qc.Find(&rows).Error

	return rows, err
}

func (dao *imMessageDao) PatchList(ctx context.Context, convID uint64, minSeq, maxSeq int64) ([]*model.IMMessage, error) {
	var rows []*model.IMMessage
	err := dao.db.Model(model.IMMessage{}).
		Where("conv_id = ?", convID).
		Where("seq > ? and seq < ?", minSeq, maxSeq).
		Order("seq desc").
		Select("id,public_id,conv_id,content,status,send_time,msg_type,sender_id,seq").
		Find(&rows).Error

	return rows, err
}

func (dao *imMessageDao) FindSeq(ctx context.Context, msgID, convID int64) (int64, error) {
	var seq int64
	err := dao.db.Model(model.IMMessage{}).
		Where("id = ? and conv_id = ?", msgID, convID).
		Select("seq").
		Take(&seq).Error
	if err != nil {
		seq = -1
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
	}

	return seq, err
}
