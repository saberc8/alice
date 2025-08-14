package repository

type FriendRepository interface {
	AddRelation(userID, friendID uint) error
	RemoveRelation(userID, friendID uint) error
	ListFriends(userID uint, offset, limit int) ([]uint, int64, error)

	// Friend Requests
	CreateRequest(requesterID, addresseeID uint) error
	AcceptRequest(requestID uint) (requesterID uint, addresseeID uint, err error)
	DeclineRequest(requestID uint) error
	GetPendingRequests(addresseeID uint, offset, limit int) ([]uint, []uint, int64, error)

	// Relationship check
	AreFriends(a, b uint) (bool, error)
}
