package repository

import momententity "alice/domain/moment/entity"

type MomentRepository interface {
	Create(m *momententity.Moment) error
	ListAll(offset, limit int) ([]*momententity.Moment, int64, error)
	ListByUser(userID uint, offset, limit int) ([]*momententity.Moment, int64, error)
	Get(id uint) (*momententity.Moment, error)
	Delete(id uint, userID uint) error
	// Likes
	AddLike(momentID, userID uint) error
	RemoveLike(momentID, userID uint) error
	HasLiked(momentID, userID uint) (bool, error)
	CountLikes(momentID uint) (int64, error)
	// Comments
	AddComment(c *momententity.MomentComment) error
	ListComments(momentID uint, offset, limit int) ([]*momententity.MomentComment, int64, error)
}
