package auth

import (
	"context"
	"fmt"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	fmt.Println(GenerateToken("hello", 001))
}
func TestFromContext(t *testing.T) {
	// 创建一个空的context
	ctx := context.Background()

	// 创建一个CurrentUser类型的值
	currentUser := &CurrentUser{
		UserID: 001,
	}

	// 将CurrentUser类型的值存储到context中
	ctx = context.WithValue(ctx, currentUserKey, currentUser)

	// 调用FromContext函数
	result := FromContext(ctx)

	// 检查返回的结果是否与预期值相等
	if result != currentUser {
		t.Errorf("Expected: %v, Got: %v", currentUser, result)
	}
}
