package user

import (
	"context"
	"errors"
	proto "github.com/Tairascii/google-docs-protos/gen/go/user"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrOnRequest = errors.New("error on request")
	ErrInvalidId = errors.New("invalid user id")
	ErrNotFound  = errors.New("user not found")
)

type UserService interface {
	IdByEmail(ctx context.Context, email string) (uuid.UUID, error)
}

type Service struct {
	client proto.UserClient
}

func NewUserService(grpcClient *grpc.ClientConn) UserService {
	c := proto.NewUserClient(grpcClient)
	return &Service{
		client: c,
	}
}

func (s *Service) IdByEmail(ctx context.Context, email string) (uuid.UUID, error) {
	resp, err := s.client.IdByEmail(ctx, &proto.IdByEmailRequest{Email: email})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return uuid.Nil, errors.Join(ErrOnRequest, err)
		}

		if st.Code() == codes.NotFound {
			return uuid.Nil, errors.Join(ErrNotFound, err)
		}

		return uuid.Nil, errors.Join(ErrOnRequest, err)
	}

	idParsed, err := uuid.Parse(resp.Id)
	if err != nil {
		return uuid.Nil, errors.Join(ErrInvalidId, err)
	}

	return idParsed, nil
}
