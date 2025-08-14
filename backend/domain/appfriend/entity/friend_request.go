package entity

import "time"

type FriendRequestStatus string

const (
	FriendRequestPending  FriendRequestStatus = "pending"
	FriendRequestAccepted FriendRequestStatus = "accepted"
	FriendRequestDeclined FriendRequestStatus = "declined"
)

type FriendRequest struct {
	ID          uint                `json:"id" gorm:"primaryKey"`
	RequesterID uint                `json:"requester_id" gorm:"not null;index"`
	AddresseeID uint                `json:"addressee_id" gorm:"not null;index"`
	Status      FriendRequestStatus `json:"status" gorm:"type:varchar(16);not null;default:'pending'"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

func (FriendRequest) TableName() string { return "app_friend_requests" }
