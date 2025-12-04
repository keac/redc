package mod

import "flag"

var (
	V             bool
	Init          bool
	List          bool
	Debug         bool
	U             string
	Name          string
	Project       string
	Start         string
	Change        string
	Stop          string
	Kill          string
	Status        string
	Node          int
	Domain        string
	Domain2       string
	Base64Command string
	Version       = "v1.0.0(2025/12/04)"
	C2Port        string
	C2Pass        string
)

func init() {

	flag.BoolVar(&V, "version", false, "显示版本号")
	flag.BoolVar(&Init, "init", false, "初始化")
	flag.BoolVar(&Debug, "debug", false, "调试")
	flag.StringVar(&U, "u", "system", "操作者")
	flag.StringVar(&Project, "p", "default", "项目名称")
	flag.BoolVar(&List, "list", false, "查看项目所有场景")
	flag.StringVar(&Start, "start", "", "开启case")
	flag.StringVar(&Kill, "kill", "", "强制关闭case")
	flag.StringVar(&Stop, "stop", "", "关闭case")
	flag.StringVar(&Change, "change", "", "更改case状态 (c2场景是切换rg ip,代理池场景是重启代理池,asm场景是重启执行器)")
	flag.StringVar(&Status, "status", "", "查看case状态")
	flag.StringVar(&Name, "name", "", "查看case状态")
	flag.IntVar(&Node, "node", 10, "机器数量(默认10)")
	flag.StringVar(&Domain, "domain", "www.amazon.com", "CS/dnslog的监听域名")
	flag.StringVar(&Domain2, "domain2", "www.amazon.com", "asmnode 备节点")
	flag.StringVar(&Base64Command, "base64command", "", "frp/nps服务端配置(base64传入)")
	flag.StringVar(&C2Port, "c2port", "8080", "c2前置场景的默认端口")
	flag.StringVar(&C2Pass, "c2pass", "changeme", "c2前置场景的默认密码")
}
