package service

import (
	"fmt"
	_ "github.com/RaymondCode/simple-demo/dal"
	"testing"
)

func TestVideoService_GetVideoById(t *testing.T) {
	fmt.Println(NewVideoService().GetVideoById(2))
}
