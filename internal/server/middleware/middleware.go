package middleware

import (
	"context"
	"errors"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const userIDKey = "user_id"

func UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("missing metadata")
		}

		userIDs := md.Get(userIDKey)
		if len(userIDs) == 0 {
			return nil, errors.New("missing user_id in metadata")
		}

		userID, err := strconv.Atoi(userIDs[0])
		if err != nil {
			return nil, errors.New("invalid user_id format, expected int")
		}

		ctx = context.WithValue(ctx, userIDKey, userID)

		return handler(ctx, req)
	}
}

func GetUserID(ctx context.Context) (int, error) {
	id, ok := ctx.Value(userIDKey).(int)
	if !ok {
		return 0, errors.New("user_id not found in context")
	}
	return id, nil
}
