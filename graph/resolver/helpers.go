package resolver

import (
	"context"
	"errors"
)

var (
	ErrorUnAuthorized = errors.New("unauthorized")
)

const adminRole = "admin"

func GetUserIDFromContext(ctx context.Context) (uint, error) {
	userID := ctx.Value("user_id")
	if userID == nil {
		return 0, ErrorUnAuthorized
	}

	if id, ok := userID.(uint); ok {
		return id, nil
	}

	return 0, ErrorUnAuthorized
}

func GetUserRoleFromContext(ctx context.Context) (string, error) {
	userRole := ctx.Value("user_role")
	if userRole == nil {
		return "", ErrorUnAuthorized
	}

	if role, ok := userRole.(string); ok {
		return role, nil
	}

	return "", ErrorUnAuthorized
}

func IsAdminFromContext(ctx context.Context) bool {
	role, err := GetUserRoleFromContext(ctx)
	if err != nil {
		return false
	}

	return role == adminRole
}

func getPagingNumbers(page, limit *int) (pageNumber, pageLimit int) {
	var p, l = 0, 0

	if page != nil {
		p = *page
	}

	if limit != nil {
		l = *limit
	}

	if p <= 0 {
		p = 1
	}

	if l <= 0 {
		l = 10
	}

	return p, l
}
