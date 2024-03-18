package mod

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

func ProjectParse(ProjectPath string, ProjectName string, User string) {
	// 确认项目文件夹是否存在,不存在就创建
	_, err := os.Stat(ProjectPath)
	if err != nil {
		// 创建项目目录
		err := os.MkdirAll(ProjectPath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			os.Exit(3)
		} else {
			fmt.Println("已创建项目目录", ProjectPath)
		}
		// 创建项目状态文件
		filePath := ProjectPath + "/project.ini"
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("项目状态文件创建失败", err)
		} else {
			fmt.Println("已创建项目状态文件", filePath)
		}
		defer file.Close()

		/*
			// 写入项目创建时间,创建者
			cfg, err := ini.Load(filePath)
			if err != nil {
				fmt.Printf("Fail to read file: %v", err)
				os.Exit(3)
			}
			cfg.Section("Global").Key("ProjectName").SetValue(ProjectName)
			cfg.Section("Global").Key("ProjectPath").SetValue(ProjectPath)
			currentTime := time.Now().Format("2006-01-02 15:04:05")
			cfg.Section("Global").Key("CreateTime").SetValue(currentTime)
			cfg.Section("Global").Key("Operator").SetValue(User)
			cfg.SaveTo(filePath)
		*/

	}

}

func ProjectConfigParse(path string) {
	cfg, err := ini.Load(path)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(3)
	}
	fmt.Println("项目名称:", cfg.Section("Global").Key("ProjectName").String())
	fmt.Println("项目路径:", cfg.Section("Global").Key("ProjectPath").String())
	fmt.Println("创建时间:", cfg.Section("Global").Key("CreateTime").String())
	fmt.Println("创建人员:", cfg.Section("Global").Key("Operator").String())
}
