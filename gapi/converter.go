package gapi

import (
	db "github.com/SergeyPanov/bank/db/sqlc"
	"github.com/SergeyPanov/bank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(dbUser db.User) *pb.User {
	return &pb.User{
		Username:          dbUser.Username,
		FullName:          dbUser.FullName,
		Email:             dbUser.Email,
		PasswordChangedAt: timestamppb.New(dbUser.PasswordChangedAt),
		CreatedAt:         timestamppb.New(dbUser.CreatedAt),
	}
}
