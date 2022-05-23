package service

import (
	"fmt"
	"testing"
)

func TestFavoriteService_IsLike(t *testing.T) {
	fmt.Println(NewFavoriteService().IsLike(1, 5))
}
