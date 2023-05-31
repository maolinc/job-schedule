package model

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

var (
	_ JobRecordModel = (*defaultJobRecordModel)(nil)
)

type (
	JobRecord struct {
		Id        int64      `gorm:"id;primary_key"` //
		JobId     int64      `gorm:"job_id"`         //任务id
		StartTime *time.Time `gorm:"start_time"`     //开始时间
		EndTime   *time.Time `gorm:"end_time"`       //结束时间
		Result    string     `gorm:"result"`         //执行结果完整信息json格式
		Status    string     `gorm:"status"`         //结果状态，ok | error
		UseMilli  int64      `gorm:"use_milli"`      //耗时,毫秒
		ExecType  string     `gorm:"exec_type"`      //耗时,毫秒
	}

	// JobRecord query cond
	JobRecordQuery struct {
		SearchBase
		JobRecord
	}

	JobRecordModel interface {
		// Trans Transaction
		Trans(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) (err error)) (err error)
		// Builder Custom assembly conditions
		Builder(ctx context.Context, db ...*gorm.DB) (b *gorm.DB)
		Insert(ctx context.Context, data *JobRecord, db ...*gorm.DB) (err error)
		Update(ctx context.Context, data *JobRecord, db ...*gorm.DB) (err error)
		// Delete When a tombstone field is set, it is tombstone. Otherwise, it is physically deleted
		Delete(ctx context.Context, id int64, db ...*gorm.DB) (err error)
		// ForceDelete Physical deletion
		ForceDelete(ctx context.Context, id int64, db ...*gorm.DB) (err error)
		Count(ctx context.Context, cond *JobRecordQuery) (total int64, err error)
		FindOne(ctx context.Context, id int64) (data *JobRecord, err error)
		// FindByPage Contains total information
		FindByPage(ctx context.Context, cond *JobRecordQuery) (total int64, list []*JobRecord, err error)
		// FindListByPage Normal pagination
		FindListByPage(ctx context.Context, cond *JobRecordQuery) (list []*JobRecord, err error)
		// FindListByCursor Cursor is required based on cursor paging, Only the primary key is of type int, and other types can be expanded by themselves
		FindListByCursor(ctx context.Context, cond *JobRecordQuery) (list []*JobRecord, err error)
		FindAll(ctx context.Context, cond *JobRecordQuery) (list []*JobRecord, err error)
		// ---------------Write your other interfaces below---------------
	}

	defaultJobRecordModel struct {
		*customConn
		table string
	}
)

func NewJobRecordModel(db *gorm.DB) JobRecordModel {
	return &defaultJobRecordModel{
		customConn: newCustomConnNoCache(db),
		table:      "job_record",
	}
}

func (m *JobRecord) TableName() string {
	return "`job_record`"
}

func (m *defaultJobRecordModel) conn(ctx context.Context, db ...*gorm.DB) *gorm.DB {
	if len(db) == 0 {
		return m.db.Model(&JobRecord{}).Session(&gorm.Session{Context: ctx})
	}
	return db[0]
}

func (m *defaultJobRecordModel) Builder(ctx context.Context, db ...*gorm.DB) (b *gorm.DB) {
	return m.conn(ctx, db...)
}

func (m *defaultJobRecordModel) Trans(ctx context.Context, fn func(ctx context.Context, db *gorm.DB) (err error)) (err error) {
	return m.conn(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, tx)
	})
}

func (m *defaultJobRecordModel) Insert(ctx context.Context, data *JobRecord, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Create(data).Error
	})
}

func (m *defaultJobRecordModel) Update(ctx context.Context, data *JobRecord, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Model(&data).Updates(data).Error
	})
}

func (m *defaultJobRecordModel) Delete(ctx context.Context, id int64, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Where(" id = ?", id).Delete(JobRecord{}).Error
	})
}

func (m *defaultJobRecordModel) ForceDelete(ctx context.Context, id int64, db ...*gorm.DB) (err error) {
	return m.Exec(ctx, func() error {
		return m.conn(ctx, db...).Unscoped().Where(" id = ?", id).Delete(JobRecord{}).Error
	})
}

func (m *defaultJobRecordModel) Count(ctx context.Context, cond *JobRecordQuery) (total int64, err error) {
	err = m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
	).Where(cond.JobRecord).Count(&total).Error
	return total, err
}

func (m *defaultJobRecordModel) FindOne(ctx context.Context, id int64) (data *JobRecord, err error) {
	err = m.QueryRow(ctx, &data, func(v interface{}) error {
		tx := m.conn(ctx).Where(" id = ?", id).Find(v)
		if tx.RowsAffected == 0 {
			return sql.ErrNoRows
		}
		return tx.Error
	})
	switch err {
	case nil:
		return data, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *defaultJobRecordModel) FindByPage(ctx context.Context, cond *JobRecordQuery) (total int64, list []*JobRecord, err error) {
	conn := m.conn(ctx).Scopes(
		orderScope(cond.OrderSort...),
		searchPlusScope(cond.SearchPlus, m.table),
	).Where(cond.JobRecord)

	total, list, err = pageHandler[*JobRecord](conn, cond.PageCurrent, cond.PageSize)
	return total, list, err
}

func (m *defaultJobRecordModel) FindListByPage(ctx context.Context, cond *JobRecordQuery) (list []*JobRecord, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		orderScope(cond.OrderSort...),
		pageScope(cond.PageCurrent, cond.PageSize),
	)
	err = conn.Where(cond.JobRecord).Find(&list).Error
	return list, err
}

func (m *defaultJobRecordModel) FindListByCursor(ctx context.Context, cond *JobRecordQuery) (list []*JobRecord, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		cursorScope(cond.Cursor, cond.CursorAsc, cond.PageSize, "id"),
		orderScope(cond.OrderSort...),
	)
	err = conn.Where(cond.JobRecord).Find(&list).Error
	return list, err
}

func (m *defaultJobRecordModel) FindAll(ctx context.Context, cond *JobRecordQuery) (list []*JobRecord, err error) {
	conn := m.conn(ctx).Scopes(
		searchPlusScope(cond.SearchPlus, m.table),
		orderScope(cond.OrderSort...),
	)
	err = conn.Where(cond.JobRecord).Find(&list).Error
	return list, err
}
