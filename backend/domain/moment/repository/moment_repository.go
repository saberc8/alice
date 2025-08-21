package repository

import momententity "alice/domain/moment/entity"

type MomentRepository interface {
	Create(m *momententity.Moment) error
	ListAll(offset, limit int) ([]*momententity.Moment, int64, error)
	ListByUser(userID uint, offset, limit int) ([]*momententity.Moment, int64, error)
	Get(id uint) (*momententity.Moment, error)
	Delete(id uint, userID uint) error
}
