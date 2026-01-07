package utils

import (
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	once     sync.Once           // once 用于确保验证器只初始化一次（单例模式）
	validate *validator.Validate // validate 是validator库的验证器实例
)

// initValidator 初始化验证器实例
// 功能：创建validator库的验证器实例
// 实现说明：
//  1. 使用validator.New()创建新的验证器实例
//  2. 该函数通过sync.Once保证只被调用一次
func initValidator() {
	validate = validator.New()
}

// Validate 暴露给业务层的唯一验证入口
// 功能：验证结构体字段是否符合定义的验证规则
// 参数：
//
//	v: 要验证的结构体实例
//
// 返回值：
//
//	error: 验证失败时返回错误信息，验证成功时返回nil
//
// 实现说明：
//  1. 使用sync.Once确保验证器只初始化一次
//  2. 调用validate.Struct()验证结构体字段
//  3. 如果验证失败，将错误信息包装后返回
func Validate(v interface{}) error {
	once.Do(initValidator)
	if err := validate.Struct(v); err != nil {
		// 把 validator.ValidationErrors 简单包装成一条 error，方便打印和处理
		return fmt.Errorf("validate failed: %w", err)
	}
	return nil
}
