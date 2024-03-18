package mod

import (
	"fmt"
	"os"
	"red-cloud/mod2"
	"red-cloud/utils"
	"strconv"
	"strings"
	"time"
)

// 第一次初始化
func TfInit0(Path string) {

	fmt.Println("cd " + Path + " && bash deploy.sh -init")
	err := utils.Command("cd " + Path + " && bash deploy.sh -init")
	//err := utils.Command("cd " + Path + " && terraform init")
	if err != nil {
		fmt.Println("场景初始化失败,再次尝试!", err)

		// 如果初始化失败就再次尝试一次
		err2 := utils.Command("cd " + Path + " && bash deploy.sh -init")
		if err2 != nil {
			fmt.Println("场景初始化失败,请检查网络连接!", err)
			os.Exit(3)
		}
	}
}

// 复制后的初始化
func TfInit(Path string) {

	fmt.Println("cd " + Path + " && bash deploy.sh -init")
	err := utils.Command("cd " + Path + " && bash deploy.sh -init")
	//err := utils.Command("cd " + Path + " && terraform init")
	if err != nil {
		fmt.Println("场景初始化失败,再次尝试!", err)

		// 如果初始化失败就再次尝试一次
		err2 := utils.Command("cd " + Path + " && bash deploy.sh -init")
		if err2 != nil {
			fmt.Println("场景初始化失败,请检查网络连接!", err)

			// 无法初始化,删除 case 文件夹
			err = os.RemoveAll(Path)
			if err != nil {
				fmt.Println(err)
				os.Exit(3)
			}
			os.Exit(3)
		}
	}
}

func TfApply(Path string) {

	fmt.Println("cd " + Path + " && bash deploy.sh -start")
	err := utils.Command("cd " + Path + " && bash deploy.sh -start")
	if err != nil {
		fmt.Println("场景创建失败!尝试重新创建!")

		// 先关闭
		err2 := utils.Command("cd " + Path + " && bash deploy.sh -stop")
		if err2 != nil {
			fmt.Println("场景销毁,等待重新创建!")
			os.Exit(3)
		}

		// 重新创建
		err3 := utils.Command("cd " + Path + " && bash deploy.sh -start")
		if err3 != nil {
			fmt.Println("场景创建第二次失败!请手动排查问题")
			fmt.Println("path路径: ", Path)
			os.Exit(3)
		}

	}
}

func TfStatus(Path string) {

	fmt.Println("cd " + Path + " && bash deploy.sh -status")
	err := utils.Command("cd " + Path + " && bash deploy.sh -status")
	if err != nil {
		fmt.Println("场景状态查询失败!请手动排查问题")
		fmt.Println("path路径: ", Path)
		os.Exit(3)
	}
}

func TfDestroy(Path string) {

	fmt.Println("cd " + Path + " && bash deploy.sh -stop")
	err := utils.Command("cd " + Path + " && bash deploy.sh -stop")
	if err != nil {
		fmt.Println("场景销毁失败,第二次尝试!", err)

		// 如果初始化失败就再次尝试一次
		err2 := utils.Command("cd " + Path + " && bash deploy.sh -stop")
		if err2 != nil {
			fmt.Println("场景销毁失败,第三次尝试!", err)

			// 第三次
			err3 := utils.Command("cd " + Path + " && bash deploy.sh -stop")
			if err3 != nil {
				fmt.Println("场景销毁多次重试失败!请手动排查问题")
				fmt.Println("path路径: ", Path)
				os.Exit(3)
			}
		}
	}

}

func C2Apply(Path string) {

	// 先开c2
	err := utils.Command("cd " + Path + " && bash deploy.sh -step1")
	if err != nil {
		fmt.Println("场景创建失败,自动销毁场景!")
		RedcLog("场景创建失败,自动销毁场景!")
		C2Destroy(Path, strconv.Itoa(Node), Domain)
		// 成功销毁场景后,删除 case 文件夹
		err = os.RemoveAll(Path)
		os.Exit(3)
	}

	// 开rg
	if Node != 0 {
		err = utils.Command("cd " + Path + " && bash deploy.sh -step2 " + strconv.Itoa(Node) + " " + Domain)
		if err != nil {
			fmt.Println("场景创建失败,自动销毁场景!")
			RedcLog("场景创建失败,自动销毁场景!")
			C2Destroy(Path, strconv.Itoa(Node), Domain)
			// 成功销毁场景后,删除 case 文件夹
			err = os.RemoveAll(Path)
			os.Exit(3)
		}
	}

	// 获得本地几个变量
	c2_ip := utils.Command2("cd " + Path + " && cd c2-ecs" + "&& terraform output -json ecs_ip | jq '.' -r")
	c2_pass := utils.Command2("cd " + Path + " && cd c2-ecs" + "&& terraform output -json ecs_password | jq '.' -r")

	cs_port := "8080"
	cs_pass := "q!w@e#raa1dd2ff3gg4"
	cs_domain := Domain
	ssh_ip := c2_ip + ":22"

	// 去掉该死的换行符
	ssh_ip = strings.Replace(ssh_ip, "\n", "", -1)
	c2_pass = strings.Replace(c2_pass, "\n", "", -1)
	c2_ip = strings.Replace(c2_ip, "\n", "", -1)

	time.Sleep(time.Second * 60)

	// ssh上去起teamserver
	if Node != 0 {
		ipsum := utils.Command2("cd " + Path + "&& cd zone-node && cat ipsum.txt")
		ecs_main_ip := utils.Command2("cd " + Path + "&& cd zone-node && cat ecs_main_ip.txt")
		ipsum = strings.Replace(ipsum, "\n", "", -1)
		ecs_main_ip = strings.Replace(ecs_main_ip, "\n", "", -1)
		cscommand := "setsid ./teamserver -new " + cs_port + " " + c2_ip + " " + cs_pass + " " + cs_domain + " " + ipsum + " " + ecs_main_ip + " > /dev/null 2>&1 &"
		fmt.Println("cscommand: ", cscommand)
		err = utils.Gotossh("root", c2_pass, ssh_ip, cscommand)
		if err != nil {
			mod2.PrintOnError(err, "ssh 过程出现报错!自动销毁场景")
			RedcLog("ssh 过程出现报错!自动销毁场景")
			C2Destroy(Path, strconv.Itoa(Node), Domain)
			// 成功销毁场景后,删除 case 文件夹
			err = os.RemoveAll(Path)
			os.Exit(3)
		}
	} else {
		cscommand := "setsid ./teamserver -new " + cs_port + " " + c2_ip + " " + cs_pass + " " + cs_domain + " > /dev/null 2>&1 &"
		fmt.Println("cscommand: ", cscommand)
		err = utils.Gotossh("root", c2_pass, ssh_ip, cscommand)
		if err != nil {
			mod2.PrintOnError(err, "ssh 过程出现报错!自动销毁场景")
			RedcLog("ssh 过程出现报错!自动销毁场景")
			C2Destroy(Path, strconv.Itoa(Node), Domain)
			// 成功销毁场景后,删除 case 文件夹
			err = os.RemoveAll(Path)
			os.Exit(3)
		}
	}

	fmt.Println("ssh结束!")

	err = utils.Command("cd " + Path + " && bash deploy.sh -status")

	if err != nil {
		mod2.PrintOnError(err, "场景创建失败")
		RedcLog("场景创建失败")
		os.Exit(3)
	}

}

func C2Change(Path string) {

	// 重开rg
	fmt.Println("cd " + Path + " && bash deploy.sh -step3 " + strconv.Itoa(Node) + " " + Domain)
	err := utils.Command("cd " + Path + " && bash deploy.sh -step3 " + strconv.Itoa(Node) + " " + Domain)
	if err != nil {
		mod2.PrintOnError(err, "场景更改失败")
		os.Exit(3)
	}

	// 获得本地几个变量
	c2_ip := utils.Command2("cd " + Path + " && cd c2-ecs" + "&& terraform output -json ecs_ip | jq '.' -r")
	c2_pass := utils.Command2("cd " + Path + " && cd c2-ecs" + "&& terraform output -json ecs_password | jq '.' -r")
	ipsum := utils.Command2("cd " + Path + "&& cd zone-node && cat ipsum.txt")
	ecs_main_ip := utils.Command2("cd " + Path + "&& cd zone-node && cat ecs_main_ip.txt")

	cs_port := "8080"
	cs_pass := "q!w@e#raa1dd2ff3gg4"
	cs_domain := "360.com"
	ssh_ip := c2_ip + ":22"

	// 去掉该死的换行符
	ssh_ip = strings.Replace(ssh_ip, "\n", "", -1)
	c2_pass = strings.Replace(c2_pass, "\n", "", -1)
	c2_ip = strings.Replace(c2_ip, "\n", "", -1)
	ipsum = strings.Replace(ipsum, "\n", "", -1)
	ecs_main_ip = strings.Replace(ecs_main_ip, "\n", "", -1)
	cscommand := "setsid ./teamserver -changelistener1 " + cs_port + " " + c2_ip + " " + cs_pass + " " + cs_domain + " " + ipsum + " " + ecs_main_ip + " > /dev/null 2>&1 &"

	// ssh上去起teamserver
	utils.Gotossh("root", c2_pass, ssh_ip, cscommand)

}

func C2Destroy(Path string, Command1 string, Domain string) {

	fmt.Println("cd " + Path + " && bash deploy.sh -stop " + Command1 + " " + Domain)
	err := utils.Command("cd " + Path + " && bash deploy.sh -stop " + Command1 + " " + Domain)
	if err != nil {
		fmt.Println("场景销毁失败,第一次尝试!", err)
		RedcLog("场景销毁失败,第一次尝试!")

		// 如果初始化失败就再次尝试一次
		err2 := utils.Command("cd " + Path + " && bash deploy.sh -stop " + Command1 + " " + Domain)
		if err2 != nil {
			fmt.Println("场景销毁失败,第二次尝试!", err)
			RedcLog("场景销毁失败,第二次尝试!")

			// 第三次
			err3 := utils.Command("cd " + Path + " && bash deploy.sh -stop " + Command1 + " " + Domain)
			if err3 != nil {
				fmt.Println("场景销毁失败!")
				RedcLog("场景销毁失败")
				os.Exit(3)
			}
		}
	}

}

func AwsProxyApply(Path string) {

	fmt.Println("cd " + Path + " && bash deploy.sh -start " + strconv.Itoa(Node))
	err := utils.Command("cd " + Path + " && bash deploy.sh -start " + strconv.Itoa(Node))
	if err != nil {
		fmt.Println("场景创建失败!")
		RedcLog("场景创建失败!")
		os.Exit(3)
	}

}

func AwsProxyDestroy(Path string, Command1 string) {

	fmt.Println("cd " + Path + " && bash deploy.sh -stop " + Command1)
	err := utils.Command("cd " + Path + " && bash deploy.sh -stop " + Command1)
	if err != nil {
		fmt.Println("场景销毁失败,第一次尝试!", err)

		// 如果初始化失败就再次尝试一次
		err2 := utils.Command("cd " + Path + " && bash deploy.sh -stop " + Command1)
		if err2 != nil {
			fmt.Println("场景销毁失败,第二次尝试!", err)

			// 第三次
			err3 := utils.Command("cd " + Path + " && bash deploy.sh -stop " + Command1)
			if err3 != nil {
				fmt.Println("场景销毁失败!")
				RedcLog("场景销毁失败!")
				os.Exit(3)
			}
		}
	}

}

func DDOSApply(Path string) {

	fmt.Println("cd " + Path + " && bash deploy.sh -start " + strconv.Itoa(Node) + " " + Durl + " " + strconv.Itoa(Dnum) + " " + strconv.Itoa(Dtime) + " " + Dmode)
	err := utils.Command("cd " + Path + " && bash deploy.sh -start " + strconv.Itoa(Node) + " " + Durl + " " + strconv.Itoa(Dnum) + " " + strconv.Itoa(Dtime) + " " + Dmode)
	if err != nil {
		fmt.Println("场景创建失败!")
		RedcLog("场景创建失败!")
		os.Exit(3)
	}

}

func DDOSDestroy(Path string, Command1 string) {

	fmt.Println("cd " + Path + " && bash deploy.sh -stop " + Command1 + "1 1 1 1")
	err := utils.Command("cd " + Path + " && bash deploy.sh -stop " + Command1 + "1 1 1 1")
	if err != nil {
		fmt.Println("场景销毁失败,第一次尝试!", err)

		// 如果初始化失败就再次尝试一次
		err2 := utils.Command("cd " + Path + " && bash deploy.sh -stop " + Command1 + "1 1 1 1")
		if err2 != nil {
			fmt.Println("场景销毁失败,第二次尝试!", err)

			// 第三次
			err3 := utils.Command("cd " + Path + " && bash deploy.sh -stop " + Command1 + "1 1 1 1")
			if err3 != nil {
				fmt.Println("场景销毁失败!")
				RedcLog("场景销毁失败!")
				os.Exit(3)
			}
		}
	}

}

func AliyunProxyApply(Path string) {

	fmt.Println("cd " + Path + " && bash deploy.sh -start " + strconv.Itoa(Node))
	err := utils.Command("cd " + Path + " && bash deploy.sh -start " + strconv.Itoa(Node))
	if err != nil {
		fmt.Println("场景创建失败!")
		RedcLog("场景创建失败!")
		os.Exit(3)
	}

}

func AliyunProxyChange(Path string) {

	// 重开proxy
	fmt.Println("cd " + Path + " && bash deploy.sh -change " + strconv.Itoa(Node))
	err := utils.Command("cd " + Path + " && bash deploy.sh -change " + strconv.Itoa(Node))
	if err != nil {
		fmt.Println("场景更改失败!")
		RedcLog("场景更改失败!")
		os.Exit(3)
	}

}

func AliyunProxyDestroy(Path string, Command1 string) {

	fmt.Println("cd " + Path + " && bash deploy.sh -stop " + Command1)
	err := utils.Command("cd " + Path + " && bash deploy.sh -stop " + Command1)
	if err != nil {
		fmt.Println("场景销毁失败,第一次尝试!", err)

		// 如果初始化失败就再次尝试一次
		err2 := utils.Command("cd " + Path + " && bash deploy.sh -stop " + Command1)
		if err2 != nil {
			fmt.Println("场景销毁失败,第二次尝试!", err)

			// 第三次
			err3 := utils.Command("cd " + Path + " && bash deploy.sh -stop " + Command1)
			if err3 != nil {
				fmt.Println("场景销毁失败!")
				RedcLog("场景销毁失败!")
				os.Exit(3)
			}
		}
	}

}

func DnslogApply(Path string) {

	fmt.Println("cd " + Path + " && bash deploy.sh -start " + Domain)
	err := utils.Command("cd " + Path + " && bash deploy.sh -start " + Domain)
	if err != nil {
		fmt.Println("场景创建失败!")
		RedcLog("场景创建失败!")
		os.Exit(3)
	}

}

func DnslogDestroy(Path string, Domain string) {

	fmt.Println("cd " + Path + " && bash deploy.sh -stop " + Domain)
	err := utils.Command("cd " + Path + " && bash deploy.sh -stop " + Domain)
	if err != nil {
		fmt.Println("场景销毁失败,第一次尝试!", err)

		// 如果初始化失败就再次尝试一次
		err2 := utils.Command("cd " + Path + " && bash deploy.sh -stop " + Domain)
		if err2 != nil {
			fmt.Println("场景销毁失败,第二次尝试!", err)

			// 第三次
			err3 := utils.Command("cd " + Path + " && bash deploy.sh -stop " + Domain)
			if err3 != nil {
				fmt.Println("场景销毁失败!")
				RedcLog("场景销毁失败!")
				os.Exit(3)
			}
		}
	}

}

func Base64Apply(Path string) {

	fmt.Println("cd " + Path + " && bash deploy.sh -start " + Base64Command)
	err := utils.Command("cd " + Path + " && bash deploy.sh -start " + Base64Command)
	if err != nil {
		fmt.Println("场景创建失败!")
		RedcLog("场景创建失败!")
		os.Exit(3)
	}

}

func Base64Destroy(Path string, Base64Command string) {

	fmt.Println("cd " + Path + " && bash deploy.sh -stop " + Base64Command)
	err := utils.Command("cd " + Path + " && bash deploy.sh -stop " + Base64Command)
	if err != nil {
		fmt.Println("场景销毁失败,第一次尝试!", err)

		// 如果初始化失败就再次尝试一次
		err2 := utils.Command("cd " + Path + " && bash deploy.sh -stop " + Base64Command)
		if err2 != nil {
			fmt.Println("场景销毁失败,第二次尝试!", err)

			// 第三次
			err3 := utils.Command("cd " + Path + " && bash deploy.sh -stop " + Base64Command)
			if err3 != nil {
				fmt.Println("场景销毁失败!")
				RedcLog("场景销毁失败!")
				os.Exit(3)
			}
		}
	}

}
