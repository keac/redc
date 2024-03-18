package mod

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

// LoadConfig 加载配置文件
func LoadConfig(path string) {
	_, err := os.Stat(path)
	if err != nil {
		// 没有配置文件，报错退出，提示进行修改

		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("配置文件创建失败", err)
		} else {
			fmt.Println("已生成状态文件", path, "请自行修改")

			cfg, err := ini.Load(path)
			if err != nil {
				fmt.Printf("Fail to read file: %v", err)
				os.Exit(3)
			}
			cfg.Section("").Key("operator").SetValue("system")
			cfg.Section("").Key("ALICLOUD_ACCESS_KEY").SetValue("changethis")
			cfg.Section("").Key("ALICLOUD_SECRET_KEY").SetValue("changethis")
			cfg.SaveTo(path)
		}
		defer file.Close()
	}

}

// ParseConfig 解析配置文件
func ParseConfig(path string) (string, string) {
	cfg, err := ini.Load(path)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	ALICLOUD_ACCESS_KEY := cfg.Section("").Key("ALICLOUD_ACCESS_KEY").String()
	ALICLOUD_SECRET_KEY := cfg.Section("").Key("ALICLOUD_SECRET_KEY").String()

	return ALICLOUD_ACCESS_KEY, ALICLOUD_SECRET_KEY
}
