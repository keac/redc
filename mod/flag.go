package mod

import "flag"

var (
	V             bool
	Init          bool
	Cost          bool
	Fc            bool
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
	Base64Command string
	Version       = "v1.3.2(2024/03/18)(wgpsec)"
	Durl          string
	Dtime         int
	Dnum          int
	Dmode         string
)

func init() {

	flag.BoolVar(&V, "version", false, "显示版本号")
	flag.BoolVar(&Init, "init", false, "初始化")
	flag.BoolVar(&Debug, "debug", false, "调试")
	flag.StringVar(&U, "u", "system", "操作者")
	flag.StringVar(&Project, "p", "default", "项目名称")
	flag.BoolVar(&List, "list", false, "查看项目所有场景")
	flag.BoolVar(&Fc, "fc", false, "查询云函数余量")
	flag.BoolVar(&Cost, "cost", false, "查看余额")
	flag.StringVar(&Start, "start", "", "开启case")
	flag.StringVar(&Kill, "kill", "", "强制关闭case")
	flag.StringVar(&Stop, "stop", "", "关闭case")
	flag.StringVar(&Change, "change", "", "更改case状态 (c2场景是切换rg ip,代理池场景是重启代理池)")
	flag.StringVar(&Status, "status", "", "查看case状态")
	flag.StringVar(&Name, "name", "", "查看case状态")
	flag.IntVar(&Node, "node", 10, "机器数量(默认10)")
	flag.StringVar(&Domain, "domain", "www.amazon.com", "CS/dnslog的监听域名")
	flag.StringVar(&Base64Command, "base64command", "", "frp/nps服务端配置(base64传入)")
	flag.StringVar(&Durl, "durl", "", "ddos目标")
	flag.IntVar(&Dtime, "dtime", 600, "ddos持续时间(默认10分钟)")
	flag.IntVar(&Dnum, "dnum", 3500, "ddos线程数(默认3500)")
	flag.StringVar(&Dmode, "dmode", "APACHE", "ddos模式")

}
