# redc

## 编译

```
goreleaser --snapshot --skip-publish --rm-dist
```

## 安装依赖工具

**mac**
```bash
brew install aliyun-cli
brew install terraform
brew install jq
```

**linux**
```bash
# terraform
mkdir -p /tmp/terraform && cd /tmp/terraform && wget -O terraform_1.6.6_linux_amd64.zip 'https://releases.hashicorp.com/terraform/1.6.6/terraform_1.6.6_linux_amd64.zip'
unzip terraform_1.6.6_linux_amd64.zip
mv --force terraform /usr/local/bin/terraform > /dev/null 2>&1 && chmod +x /usr/local/bin/terraform
rm -rf /tmp/terraform

# 如果是 arm64 机器
# mkdir -p /tmp/terraform && cd /tmp/terraform && wget -O terraform_1.6.6_linux_arm64.zip 'https://releases.hashicorp.com/terraform/1.6.6/terraform_1.6.6_linux_arm64.zip'
# unzip terraform_1.6.6_linux_arm64.zip
# mv --force terraform /usr/local/bin/terraform > /dev/null 2>&1 && chmod +x /usr/local/bin/terraform
# rm -rf /tmp/terraform

cd /tmp
terraform -version

# aliyun
mkdir -p /tmp/aliyuncli && cd /tmp/aliyuncli && wget -O aliyun-cli-linux-latest-amd64.tgz 'https://aliyuncli.alicdn.com/aliyun-cli-linux-latest-amd64.tgz?file=aliyun-cli-linux-latest-amd64.tgz'
tar -xzvf aliyun-cli-linux-latest-amd64.tgz
mv --force aliyun /usr/local/bin/aliyun > /dev/null 2>&1 && chmod +x /usr/local/bin/aliyun
rm -rf /tmp/aliyuncli

apt install jq || yum install jq

# aws
https://docs.aws.amazon.com/zh_cn/cli/latest/userguide/getting-started-install.html
aws configure
```

**windows**
```
https://github.com/aliyun/aliyun-cli/releases/download/v3.0.121/aliyun-cli-windows-3.0.121-amd64.zip
https://releases.hashicorp.com/terraform/1.2.3/terraform_1.2.3_windows_amd64.zip
```

## 配置

```bash
aliyun configure set --profile cloud-tool --mode AK --region cn-beijing --access-key-id xxxxxxxxxxxxxx --access-key-secret xxxxxxxxxxxxxx
```

配置tf插件缓存路径
```bash
echo 'plugin_cache_dir = "$HOME/.terraform.d/plugin-cache"' > ~/.terraformrc
```

使用前需初始化redc，将自动下载tf模块依赖
```
./redc -init
```

## 思路

1. 先创建新项目
2. 指定项目下要创建场景会从场景库复制一份场景文件夹到项目文件夹下
3. 不同项目下创建同一场景互不干扰
4. 同一项目下创建同一场景互不干扰
5. 多用户操作互不干扰(本地有做鉴权,但这个实际上要在平台上去做)

- redc 配置文件 (.redc.ini)
- 项目1 (./project1)
    - 场景1 (./project1/[uuid1])
        - main.tf
        - version.tf
        - output.tf
    - 场景2 (./project1/[uuid2])
        - main.tf
        - version.tf
        - output.tf
    - 项目状态文件 (project.ini)
- 项目2 (./project2)
    - 场景1 (./project2/[uuid1])
        - main.tf
        - version.tf
        - output.tf
    - 场景2 (./project2/[uuid2])
        - ...
    - 项目状态文件 (project.ini)
- 项目3 (./project3)
    - ...

## 交互

```bash
# 项目 test 开启 awvs 场景
redc -project test -start awvs -u zhangsan
.........
.........
项目uuid:xxxxxxxxx

# 查看 test 项目中指定场景的状态
redc -project test -status [uuid] -u zhangsan

# 关闭 test 项目中指定场景
redc -project test -stop [uuid] -u zhangsan

# 查看 test 项目的所有场景
redc -project test -list -u zhangsan
uuid        type    createtime      operator
xxxxxxxxx   awvs    2022.02.22      system
bbbbbbbbb   file    2022.02.22      system
```

场景名称
```
awvs
file
chat
c2
nessus
proxy
pupy
```

---

## 设计规划

tf 分成 2 类场景
- 基础场景
- 复杂场景 (由基础场景修改而来)

redc 考虑是给予平台使用，在平台上由多项目、多用户进行操作,同时兼顾单机版需求

由于 tf 的局限性，使用时和其文件夹结构脱不开关联，在多用户的情况下需要用 Backend 同步状态锁，融入到平台虽然可以用 Consul 解决多用户操作的问题，但多项目下要使用依然无法解决

无法让多项目用1个文件夹场景，如果多复制几个文件夹太过笨重。。。这些都不够高效率
