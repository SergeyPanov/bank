package gapi

import (
	"context"
	"database/sql"

	db "github.com/SergeyPanov/bank/db/sqlc"
	"github.com/SergeyPanov/bank/pb"
	"github.com/SergeyPanov/bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := s.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user has not been found")
		}
		return nil, status.Errorf(codes.Internal, "failed to retrieve the user: %s", err)
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid password")
	}

	accessToken, accessTokenPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to create an access token: %s", err)
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to create a refresh token: %s", err)
	}

	mtdt := s.extractMetadata(ctx)

	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID: refreshPayload.ID,
		Username: sql.NullString{
			Valid:  true,
			String: user.Username,
		},
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIp,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to create a session: %s", err)
	}

	resp := pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		User:                  convertUser(user),
	}

	return &resp, nil
}
