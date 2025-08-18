package types

type SendRequestParams struct {
	ReceiverUsername string `json:"friend_username" validate:"required,min=1,max=20"`
}

type AcceptRequestParams struct {
	FriendshipID string `json:"friendship_id" validate:"required,min=1,max=20"`
	SenderID     string `json:"sender_id" validate:"required,min=1,max=20"`
}

type RemoveFriendParams struct {
	FriendshipID string `json:"friendship_id" validate:"required,min=1,max=20"`
	SenderID     string `json:"sender_id" validate:"required,min=1,max=20"`
	ReceiverID   string `json:"receiver_id" validate:"required,min=1,max=20"`
}
