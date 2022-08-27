package gapi

import (
	"context"
	"database/sql"

	db "github.com/Rexkizzy22/micro-bank/db/sqlc"
	"github.com/Rexkizzy22/micro-bank/pb"
	"github.com/Rexkizzy22/micro-bank/util"
	"github.com/Rexkizzy22/micro-bank/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	if violations := validateExistingUser(req); violations != nil {
		return nil, inValidArgumentError(violations)
	}

	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to find user")
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Error(codes.NotFound, "incorrect password")
	}

	accessToken, accessPayload, err := server.token.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create access token")
	}

	refreshToken, refreshPayload, err := server.token.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create refresh token")
	}

	metad := server.getMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     refreshPayload.Username,
		RefreshToken: refreshToken,
		UserAgent:    metad.UserAgent,
		ClientIP:     metad.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create session")
	}

	resp := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		User:                  converter(user),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return resp, nil
}

func validateExistingUser(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := validation.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
