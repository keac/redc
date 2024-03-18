package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/beevik/ntp"
	"os"
	"os/exec"
	redc "red-cloud/mod"
	"red-cloud/mod2"
	"red-cloud/utils"
	"time"
)

var ProjectPath = "./redc-taskresult"

func main() {

	// ntp校验
	//CheckStatus()
	flag.Parse()

	if redc.Debug {
		DebugFunc()
		os.Exit(0)
	}

	// -version 显示版本号
	if redc.V {
		fmt.Println(redc.Version)
		os.Exit(0)
	}

	// 解析配置(暂时不需要这一步)
	// redc.LoadConfig(configPath)

	// -init 初始化
	if redc.Init {
		redc.RedcLog("进行初始化")
		fmt.Println("初始化中")
		// 先删除文件夹
		err := os.RemoveAll("redc-templates")
		mod2.PrintOnError(err, "初始化过程中删除模板文件夹失败")
		// 释放templates资源
		utils.ReleaseDir("redc-templates")

		// 遍历 redc-templates 文件夹,不包括子目录
		_, dirs := utils.GetFilesAndDirs("./redc-templates")
		for _, v := range dirs {
			redc.TfInit0(v)
		}

		// 遍历 redc-templates 文件夹,包括子目录 (现已被替代)
		/*dirs := utils.ChechDirMain("./redc-templates")
		for _, v := range dirs {
			err := utils.CheckFileName(v, "tf")
			if err {
				fmt.Println(v)
				redc.TfInit(v)
			}
		}*/
		os.Exit(0)
	}

	// 解析项目名称
	redc.ProjectParse(ProjectPath+"/"+redc.Project, redc.Project, redc.U)

	// list 操作查看项目里所有 case
	if redc.List {
		redc.CaseList(ProjectPath + "/" + redc.Project)
	}

	if redc.Cost {
		redc.RedcLog("查看余额")
		fmt.Print("阿里云当前余额: ")
		err := utils.Command("aliyun bssopenapi QueryAccountBalance --region cn-beijing | jq -r .Data.AvailableAmount")
		if err != nil {
			fmt.Println("查询阿里云当前余额失败!", err)
		}

		fmt.Print("华为云当前余额: ")
		err2 := utils.Command("hcloud BSS ShowCustomerAccountBalances --cli-region=\"cn-north-1\" | jq .account_balances | jq '.[1] | .amount'")
		if err2 != nil {
			fmt.Println("查询华为云当前余额失败!", err2)
		}

		fmt.Print("腾讯云当前余额: ")
		err3 := utils.Command("tccli billing DescribeAccountBalance --cli-unfold-argument | jq '.Balance | tonumber/100'")
		if err3 != nil {
			fmt.Println("查询腾讯云当前余额失败!", err3)
		}
	}

	if redc.Fc {
		redc.RedcLog("查看云函数余量")
		// https://next.api.aliyun.com/api/BssOpenApi/2017-12-14/QueryResourcePackageInstances?tab=CLI
		fmt.Print("阿里云Fc当前余量: \n")
		err := utils.CommandUTF("aliyun bssopenapi QueryResourcePackageInstances --region cn-beijing | jq .Data.Instances.Instance | jq -r '.[] | \"\\(.Remark): \\(.RemainingAmount) \\(.TotalAmountUnit)\"'")
		if err != nil {
			fmt.Println("查询阿里云SCF当前余量失败!", err)
		}

		fmt.Print("\n腾讯云scf目前不支持查询余量,请到控制台查看. \n")
	}

	// start 操作,去调用 case 创建方法
	if redc.Start != "" {
		redc.RedcLog("start " + redc.Start)
		if redc.Start == "pte" {
			redc.Start = "pte_arm"
		}
		//fmt.Println("step1")
		redc.CaseCreate(ProjectPath+"/"+redc.Project, redc.Start, redc.U, redc.Name)
	}

	// stop 操作,去调用 case 删除方法
	if redc.Stop != "" {
		redc.RedcLog("stop " + redc.Stop)
		redc.CheckUser(ProjectPath+"/"+redc.Project, redc.Stop)
		redc.CaseStop(ProjectPath+"/"+redc.Project, redc.Stop)
	}
	if redc.Kill != "" {
		redc.CheckUser(ProjectPath+"/"+redc.Project, redc.Kill)
		redc.CaseKill(ProjectPath+"/"+redc.Project, redc.Kill)
	}

	// change 操作,去调用 case 更改方法
	if redc.Change != "" {
		redc.RedcLog("change " + redc.Change)
		redc.CheckUser(ProjectPath+"/"+redc.Project, redc.Change)
		redc.CaseChange(ProjectPath+"/"+redc.Project, redc.Change)
	}

	// status 操作,去调用 case 状态方法
	if redc.Status != "" {
		redc.RedcLog("status" + redc.Status)
		redc.CheckUser(ProjectPath+"/"+redc.Project, redc.Status)
		redc.CaseStatus(ProjectPath+"/"+redc.Project, redc.Status)
	}

}

func DebugFunc() {

}

// CheckStatus 有效期,过期后调用自删除
func CheckStatus() {
	now := time.Now()

	// 连接超时时间
	timeout := 1 * time.Second

	// 尝试连接3次
	var response *ntp.Response
	var err error
	for i := 0; i < 4; i++ {
		response, err = ntp.QueryWithOptions("ntp1.aliyun.com", ntp.QueryOptions{Timeout: timeout})
		if err != nil {
			//fmt.Printf("第 %d 次连接失败： %s\n", i+1, err)
			continue
		}
		break
	}

	if err != nil {
		mod2.ExitOnError(err, "连接 NTP 服务器失败")
	}

	// 获取当前时间
	now = time.Now()

	// 计算偏移量
	offset := response.ClockOffset

	// 校正时间
	corrected := now.Add(offset)

	// 指定过期时间
	expireTime := time.Date(2024, 6, 10, 0, 0, 0, 0, time.Local)

	// 比较当前时间和过期时间
	if corrected.After(expireTime) {
		fmt.Println("当前时间：", now)
		fmt.Println("已过期")
		NoFile()
		os.Exit(1)
	}
}

// NoFile linux落地删、进程隐藏
func NoFile() {
	exePath, _ := os.Executable()
	cmd := exec.Command("sh", "-c", "rm -f "+exePath)
	cmd.Start()
	cmd = exec.Command("sh", "-c", "rm -f nohup.out")
	cmd.Start()
}
