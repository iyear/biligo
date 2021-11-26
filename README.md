![](https://socialify.git.ci/iyear/biligo/image?description=1&font=Raleway&forks=1&issues=1&logo=https%3A%2F%2Fraw.githubusercontent.com%2Fiyear%2Fbiligo%2Fv0%2Flogo%2Fbilibili.png&owner=1&pattern=Circuit%20Board&stargazers=1&theme=Light)

## 简介
![](https://img.shields.io/github/go-mod/go-version/iyear/biligo?style=flat-square)
![](https://img.shields.io/badge/license-GPL-lightgrey.svg?style=flat-square)
![](https://img.shields.io/github/v/release/iyear/biligo?color=red&style=flat-square)
![](https://img.shields.io/github/last-commit/iyear/biligo?style=flat-square)
[![Go Reference](https://pkg.go.dev/badge/github.com/iyear/biligo.svg)](https://pkg.go.dev/github.com/iyear/biligo)

**v0版本不保证对外函数、结构的不变性，请勿大规模用于生产环境**

`BiliBili API` 的 `Golang` 实现，目前已经实现了 100+ API，还在进一步更新中

### 特性

- 良好的设计，支持自定义 `client` 与 `UA`
- 完善的单元测试，易懂的函数命名，极少的第三方库依赖
- 代码、结构体注释完善，无需文档开箱即用
- 其他功能性代码，例如 `AV/BV`互转，`GetVideoZone()`获取分区信息...
- 配套工具 [biligo-live](https://github.com/iyear/biligo-live) 封装直播 `WebSocket` 协议
### 说明

- 该项目永远不会编写直接涉及滥用的接口
- 该项目仅供学习，请勿用于商业用途。任何使用该项目造成的后果由开发者自行承担

### 参考

> 感谢以下所有项目的贡献者

https://github.com/SocialSisterYi/bilibili-API-collect

https://github.com/MoyuScript/bilibili-api

## 快速开始

### 安装
请让 `biligo` 永远保持在最新版本

```shell
go get -u github.com/iyear/biligo
```

```go
import "github.com/iyear/biligo"
```

### 使用

```go
package main

import (
	"fmt"
	bg "github.com/iyear/biligo"
	"log"
)

func main() {
	b, err := bg.NewBiliClient(&bg.BiliSetting{
		// 均来自Web登录的Cookie
		Auth: &bg.CookieAuth{
			// DedeUserID
			DedeUserID: "YOUR_DedeUserID",
			// SESSDATA
			SESSDATA: "YOUR_SESSDATA",
			// bili_jct
			BiliJCT: "YOUR_BiliJCT",
			// DedeUserID__ckMd5
			DedeUserIDCkMd5: "YOUR_DedeUserIdCkMd5",
		},
		// DEBUG 模式将输出请求和响应
		DebugMode: true,
		// Client: myClient,
		// UserAgent: "My UA",
	})

	if err != nil {
		log.Fatal("failed to make new bili client; error: ", err)
		return
	}

	fmt.Printf("mid: %d, uname: %s,userID: %s,rank: %s\n", b.Me.MID, b.Me.UName, b.Me.UserID, b.Me.Rank)
	fmt.Printf("birthday: %s,sex: %s\n", b.Me.Birthday, b.Me.Sex)
	fmt.Printf("sign: %s\n", b.Me.Sign)
}
```

### 例子

同目录下的 `example` 文件夹

### 说明

共有两种 `Client`

一种为 `BiliClient`，必须提供有效的 `Auth` 信息，公共接口不在此提供

另一种为 `CommClient` ，只提供公共接口

一些接口同时在两种 `Client` 中，因为其中某些字段在登录与非登录状态下有所区别

## 约定

> 除了维护API列表、 `example` 目录、测试文件，不会有任何其他的专用文档。请利用好注释和测试。

> 常规的函数命名方式
>
> 由于API众多，函数名显得冗长。只能尽可能简化，但不会失去原有语义，这是无文档的基础

除无对应分类的API外，其余对应分类API均以 `分类名+动作+描述` 命名 

无对应分类的API以 `动作+描述` 命名 例如: `VideoGetTags` `VideoIsFavoured` `VideoSetFavour`

- `Video` - `视频`
- `Audio` - `音频`
- `Fav` - `收藏夹`
- `Chan` - `频道`
- `Space` - `个人空间`
- `Danmaku` - `弹幕`
- `Emote` - `表情`
- `Charge` - `充电`
- `Dyna` - `动态`
- `Live` - `直播`
- `Followings` - `关注` 


> 结构体编写规范

嵌套可分离可不分离，结构体和数组用指针。

具体看 `types.go` 即可

> 只要以下名词出现在该库的任何地方，可能将不做任何解释，请提交PR时尽可能遵守以下约定

- `aid` - `稿件av` (库中所有参数均使用 `aid` ， `bvid` 请开发者自行使用 `BV2AV()`)
- `cid` - `稿件分P的ID` (单P视频中`aid` ≠ `cid`) 或 `频道ID`
- `mid` - `B站ID`
- `mlid` - `收藏夹ID`
- `auid` - `音频ID`
- `dmid` - `弹幕ID`
- `dyid` - `动态ID`
- `dfid` - `定时发布动态ID`
- `oid` - 根据上下文不同分别指代以上不同的ID，具体看注释
- `ID` 而不是 `Id`
- `URL` 而不是 `Url`

> 关于大接口

对于评论、搜索这种超大的响应，可能不会做对应的接口，一方面这种接口经常变更，另一方面大多数的响应字段是无用的，所以请开发者自行使用 `Raw() or RawParse() + gjson` 获取需要的数据

示例：

```go
package main

import (
	"fmt"
	bg "github.com/iyear/biligo"
	"github.com/tidwall/gjson"
	"log"
)

func GetComments() {
	resp, err := bg.NewCommClient(&bg.CommSetting{}).RawParse(
		bg.BiliApiURL,
		"x/v2/reply/main",
		"GET",
		map[string]string{
			"next": "0",
			"type": "1",
			"oid":  "591018280",
			"mode": "3",
			"plat": "1",
		},
	)
	if err != nil {
		log.Println(err)
		return
	}
	for i, r := range gjson.Get(string(resp.Data), "replies.#.content.message").Array() {
		fmt.Printf("[%d] %s\n\n", i+1, r.String())
	}
}
```

## 维护与贡献
- 这个库的日常维护基本都是体力活
- 安全水平有限，部分新接口我无法分析
- 欢迎 `typo` 和 完善注释、文档的PR
- 不接受无脑无贡献的催更
- 尽量遵守约定提交PR
- 由于项目特殊性，随时可能删库，随时可能弃坑
## 已实现API

[v0](https://github.com/iyear/biligo/blob/v0/API.md)

## 留言
### 1
这个库是在暑假开始写的，其实写到一半已经不想写了，`90%`的代码都是在重复着查接口、写响应结构体、测试的过程，只是看到 `python` 有一套还在更新的 `SDK`，想着还是为 `Golang` 也写一套

后来某天突然想把这个项目给弄完 ，所以目前动态与直播的部分基本没写

### 2

`map[string]string` 而不是 `map[string]interface{}` 因为在这个项目里 `interface{}` 的无约束很容易在一些特定的参数下忘记转换。

## LICENSE

GPL
