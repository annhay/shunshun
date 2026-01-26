package initialization

import (
	"fmt"
	"log"
	"shunshun/internal/pkg/global"

	"github.com/spf13/viper"
)

// ViperInit 初始化Viper配置
// 
// 功能:
//   1. 加载配置文件
//   2. 解析配置到全局变量
//   3. 打印配置信息
// 
// 配置文件路径: ../../configs/config.yaml
func ViperInit() {
	// 设置配置文件路径
	viper.SetConfigFile("../../configs/config.yaml")
	
	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}
	
	// 解析配置到全局变量
	err = viper.Unmarshal(&global.AppConf)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}
	
	// 打印配置信息
	log.Println("viper 动态配置完成:", global.AppConf)
}
