package service

import (
	"fmt"
	"testing"
)

func TestFavoriteService_IsLike(t *testing.T) {
	fmt.Println(NewFavoriteService().IsLike(1, 5))
}

func TestFavoriteService_Recommend(t *testing.T) {
	fmt.Println(NewFavoriteService().Recommend(9))
}
