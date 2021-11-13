package biligo

import (
	"encoding/json"
	"io"
)

type Response struct {
	Code    int             `json:"code,omitempty"`
	Message string          `json:"message,omitempty"`
	TTL     int             `json:"ttl,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}
type Account struct {
	MID      int64  `json:"mid"`       // 我的mid
	UName    string `json:"uname"`     // 我的昵称
	UserID   string `json:"userid"`    // 我的用户名
	Sign     string `json:"sign"`      // 我的签名
	Birthday string `json:"birthday"`  // 我的生日 YYYY-MM-DD
	Sex      string `json:"sex"`       // 我的性别 男 女 保密
	NickFree bool   `json:"nick_free"` // 是否未设置昵称 false：设置过昵称 true：未设置昵称
	Rank     string `json:"rank"`      // 我的会员等级
}
type NavStat struct {
	Following    int `json:"following"`     // 关注数
	Follower     int `json:"follower"`      // 粉丝数
	DynamicCount int `json:"dynamic_count"` // 发布动态数
}
type VideoRecommendInfo struct {
	videoBase
}

// ExpRewardStat 完成为true 未完成为false
type ExpRewardStat struct {
	Login        bool `json:"login"`         // 每日登录 5经验
	Watch        bool `json:"watch"`         // 每日观看 5经验
	Coins        int  `json:"coins"`         // 每日投币所奖励的经验 上限50经验 该值更新存在延迟 想要无延迟请使用 GetExpCoinRewardStat
	Share        bool `json:"share"`         // 每日分享 5经验
	Email        bool `json:"email"`         // 绑定邮箱
	Tel          bool `json:"tel"`           // 绑定手机号 首次完成100经验
	SafeQuestion bool `json:"safe_question"` // 设置密保问题
	IdentifyCard bool `json:"identify_card"` // 实名认证 50经验
}
type VipStat struct {
	MID     int `json:"mid"`      // 用户MID
	VipType int `json:"vip_type"` // 大会员类型	0:无 1:月度 2:年度
	// 大会员状态
	//
	// 1:正常
	//
	// 2:由于IP地址更换过于频繁,服务被冻结
	//
	// 3:你的大会员账号风险过高，大会员功能已被锁定
	VipStatus  int   `json:"vip_status"`
	VipDueDate int64 `json:"vip_due_date"` // 大会员到期时间 时间戳(东八区) 毫秒
	VipPayType int   `json:"vip_pay_type"` // 是否已购买大会员 0:未购买 1:已购买
	ThemeType  int   `json:"theme_type"`   // 0 作用尚不明确
}
type RealNameInfo struct {
	Status   int    `json:"status"`   // 认证状态 1:已认证 3:未认证
	Remark   string `json:"remark"`   // 驳回信息 默认为空
	Realname string `json:"realname"` // 实名姓名 星号隐藏完全信息
	Card     string `json:"card"`     // 证件号码 星号隐藏部分信息
	// 证件类型代码
	//
	// 0:身份证
	//
	// 2:港澳居民来往内地通行证
	//
	// 3:台湾居民来往大陆通行证
	//
	// 4:护照(中国签发)
	//
	// 5:外国人永久居留证
	//
	// 6:其他国家或地区身份证明
	CardType int `json:"card_type"`
}
type CoinLog struct {
	Time   string  `json:"time"`   // 变化时间 YYYY-MM-DD HH:MM:SS
	Delta  float64 `json:"delta"`  // 变化量 正值为收入，负值为支出
	Reason string  `json:"reason"` // 变化说明
}
type RelationStat struct {
	MID       int64 `json:"mid"`       // 目标用户mid
	Following int   `json:"following"` // 关注数
	Whisper   int   `json:"whisper"`   // 悄悄关注数 需要登录(Cookie或APP) 未登录或非自己恒为0
	Black     int   `json:"black"`     // 黑名单数 需要登录(Cookie或APP) 未登录或非自己恒为0
	Follower  int   `json:"follower"`  // 粉丝数
}
type UpStat struct {
	Archive *UpStatArchive `json:"archive"` // 视频播放量
	Article *UpStatArticle `json:"article"` // 专栏阅读量
	Likes   int64          `json:"likes"`   // 获赞次数
}
type SpaceVideo struct {
	videoBase
	Reason     string `json:"reason"`      // 置顶视频备注
	InterVideo bool   `json:"inter_video"` // 是否为合作视频
}
type ChanVideo struct {
	List *ChanVideoList `json:"list"` // 频道信息
	Page *ChanVideoPage `json:"page"` // 页面信息
}
type ChanVideoPage struct {
	Count int // 总计视频数
	Num   int // 当前页码(可以用于请求下一页的数据计算)
	Size  int // 每页项数(可以用于请求下一页的数据计算)
}
type ChanVideoList struct {
	Archives []*ChanVideoInfo `json:"archives"` // 包含的视频列表
	CID      int64            `json:"cid"`      // 频道id
	Count    int              `json:"count"`    // 频道内含视频数
	Cover    string           `json:"cover"`    // 封面图片url
	Intro    string           `json:"intro"`    // 简介 无则为空
	MID      int64            `json:"mid"`      // 创建用户mid
	Mtime    int64            `json:"mtime"`    // 创建时间 时间戳
	Name     string           `json:"name"`     // 标题
}
type ChanVideoInfo struct {
	videoBase
	InterVideo bool `json:"inter_video"` // 是否为合作视频
}
type UpStatArchive struct {
	View int64 `json:"view"` // 视频播放量
}
type UpStatArticle struct {
	View int64 `json:"view"` // 专栏阅读量
}
type SpaceGame struct {
	Website string `json:"website"` // 游戏主页链接ur
	Image   string `json:"image"`   // 游戏图片url
	Name    string `json:"name"`    // 游戏名
}
type AccountSafetyStat struct {
	AccountInfo  *AccountSafetyInfo  `json:"account_info"` // 账号绑定信息
	AccountSafe  *AccountSafetySafe  `json:"account_safe"` // 密码安全信息
	AccountSNS   *AccountSafetySNS   `json:"account_sns"`  // 互联登录绑定信息
	AccountOther *AccountSafetyOther `json:"account_other"`
}
type AccountSafetyInfo struct {
	HideTel           string `json:"hide_tel"`           // 绑定的手机号	星号隐藏部分信息
	HideMail          string `json:"hide_mail"`          // 绑定的邮箱 星号隐藏部分信息
	BindTel           bool   `json:"bind_tel"`           // 是否绑定手机号
	BindMail          bool   `json:"bind_mail"`          // 是否绑定邮箱
	TelVerify         bool   `json:"tel_verify"`         // 是否验证手机号
	MailVerify        bool   `json:"mail_verify"`        // 是否验证邮箱
	UnneededCheck     bool   `json:"unneeded_check"`     // 是否未设置密码 注意:true为未设置，false为已设置
	RealnameCertified bool   `json:"realname_certified"` // 是否实名认证 文档中未更新此项
}
type AccountSafetySafe struct {
	Score    int  `json:"score"`     // 当前密码强度 0-100
	PwdLevel int  `json:"pwd_level"` // 当前密码强度等级 1:弱 2:中 3:强
	Security bool `json:"security"`  // 当前密码是否安全
}
type AccountSafetySNS struct {
	WeiboBind  int `json:"weibo_bind"`  // 是否绑定微博 0:未绑定 1:已绑定
	QQBind     int `json:"qq_bind"`     // 是否绑定qq	0:未绑定 1:已绑定
	WechatBind int `json:"wechat_bind"` // 是否绑定微信	0:未绑定 1:已绑定 文档中未更新此项
}
type AccountSafetyOther struct {
	SkipVerify bool `json:"skipVerify"` // 恒为false 作用尚不明确
}
type DanmakuResp struct {
	Danmaku []*Danmaku `json:"danmaku"`
}
type Danmaku struct {
	ID       uint64 `json:"id"`       // 弹幕dmid
	Progress int64  `json:"progress"` // 弹幕出现位置(单位ms)
	// 弹幕类型
	//
	// 1 2 3：普通弹幕
	//
	// 4：底部弹幕
	//
	// 5：顶部弹幕
	//
	// 6：逆向弹幕
	//
	// 7：高级弹幕
	//
	// 8：代码弹幕
	//
	// 9：BAS弹幕（仅限于特殊弹幕专包）
	Mode     int    `json:"mode"`
	FontSize int    `json:"font_size"` // 弹幕字号 18：小 25：标准 36：大
	Color    int    `json:"color"`     // 弹幕颜色 十进制RGB888值
	MidHash  string `json:"mid_hash"`  // 发送着mid hash 用于屏蔽用户和查看用户发送的所有弹幕 也可反查用户id
	Content  string `json:"content"`   // 弹幕正文 utf-8编码
	Ctime    int64  `json:"ctime"`     // 发送时间 时间戳
	Weight   int    `json:"weight"`    // 权重 区间:[1,10] 用于智能屏蔽，根据弹幕语义及长度通过AI识别得出 值越大权重越高
	Action   string `json:"action"`    // 动作 作用尚不明确
	Pool     int    `json:"pool"`      // 弹幕池 0：普通池 1：字幕池 2：特殊池（代码/BAS弹幕）
	IDStr    string `json:"id_str"`    // 弹幕dmid的字符串形式
	Attr     int    `json:"attr"`      // 弹幕属性位(bin求AND) bit0:保护 bit1:直播 bit2:高赞
}
type SpaceVideoSearchResult struct {
	List           *SpaceVideoSearchList           `json:"list"`            // 列表信息
	Page           *SpaceVideoSearchPage           `json:"page"`            // 页面信息
	EpisodicButton *SpaceVideoSearchEpisodicButton `json:"episodic_button"` // “播放全部“按钮
}
type SpaceVideoSearchList struct {
	Tlist map[string]*SpaceVideoSearchTList `json:"tlist"` // 投稿视频分区索引 key为tid字符串形式，value为详细信息
	Vlist []*SpaceVideoSearchVList          `json:"vlist"` // 投稿视频列表
}
type SpaceVideoSearchTList struct {
	TID   int    `json:"tid"`   // 该分区tid
	Count int    `json:"count"` // 投稿至该分区的视频数
	Name  string `json:"name"`  // 该分区名称
}
type SpaceVideoSearchPage struct {
	Count int `json:"count"` // 总计稿件数
	PN    int `json:"pn"`    // 当前页码(可以用于请求下一页的数据计算)
	PS    int `json:"ps"`    // 当前每页项数(可以用于请求下一页的数据计算)
}
type SpaceVideoSearchEpisodicButton struct {
	Text string `json:"text"` // 按钮文字
	Uri  string `json:"uri"`  // 全部播放页url(经测试页面空白...)
}
type SpaceVideoSearchVList struct {
	AID          int64  `json:"aid"`            // 稿件avid
	Author       string `json:"author"`         // 视频UP主 不一定为目标用户（合作视频）
	BVID         string `json:"bvid"`           // 稿件bvid
	Comment      int    `json:"comment"`        // 视频评论数
	Copyright    string `json:"copyright"`      // 空 作用尚不明确
	Created      int64  `json:"created"`        // 投稿时间 时间戳
	Description  string `json:"description"`    // 视频简介
	HideClick    bool   `json:"hide_click"`     // 恒为false 作用尚不明确
	IsPay        int    `json:"is_pay"`         // 恒为0 作用尚不明确
	IsUnionVideo int    `json:"is_union_video"` // 是否为合作视频 0：否 1：是
	Length       string `json:"length"`         // 视频长度 MM:SS
	MID          int64  `json:"mid"`            // 视频UP主mid 不一定为目标用户（合作视频）
	Pic          string `json:"pic"`            // 视频封面
	Play         int64  `json:"play"`           // 视频播放次数
	Review       int    `json:"review"`         // 恒为0 作用尚不明确
	Subtitle     string `json:"subtitle"`       // 恒为空 作用尚不明确
	Title        string `json:"title"`          // 视频标题
	TypeID       int    `json:"typeid"`         // 视频分区tid
	VideoReview  int    `json:"video_review"`   // 视频弹幕数
}
type ChannelList struct {
	Count int         `json:"count"` // 总计频道数
	List  []*ChanInfo `json:"list"`  // 频道列表
}
type FavoritesList struct {
	Count int        `json:"count"` // 创建的收藏夹数
	List  []*FavInfo `json:"list"`  // 收藏夹列表
}
type FavInfo struct {
	ID         int64  `json:"id"`          // 收藏夹mlid
	FID        int64  `json:"fid"`         // 原始收藏夹mlid 去除lmid最后两位
	MID        int64  `json:"mid"`         // 创建用户mid
	Attr       int    `json:"attr"`        // 收藏夹属性位配置
	Title      string `json:"title"`       // 收藏夹标题
	FavState   int    `json:"fav_state"`   // 0 作用尚不明确
	MediaCount int    `json:"media_count"` // 收藏夹总计视频数
}
type FavDetail struct {
	ID        int64           `json:"id"`         // 收藏夹mlid（完整id） 收藏夹原始id+创建者mid尾号2位
	FID       int64           `json:"fid"`        // 收藏夹原始id
	MID       int64           `json:"mid"`        // 创建者mid
	Attr      int             `json:"attr"`       // 属性位（？）
	Title     string          `json:"title"`      // 收藏夹标题
	Cover     string          `json:"cover"`      // 收藏夹封面图片url
	Upper     *FavDetailUpper `json:"upper"`      // 创建者信息
	CoverType int             `json:"cover_type"` // 封面图类别（？）
	CntInfo   *FavDetailCnt   `json:"cnt_info"`   // 收藏夹状态数
	Type      int             `json:"type"`       // 类型（？） 一般是11
	Intro     string          `json:"intro"`      // 备注
	Ctime     int64           `json:"ctime"`      // 创建时间	时间戳
	Mtime     int64           `json:"mtime"`      // 收藏时间	时间戳
	State     int             `json:"state"`      // 状态（？） 一般为0
	// 收藏夹收藏状态
	//
	// 已收藏收藏夹：1
	//
	// 未收藏收藏夹：0
	//
	// 需要登录
	FavState int `json:"fav_state"`
	// 点赞状态
	//
	// 已点赞：1
	//
	// 未点赞：0
	//
	// 需要登录
	LikeState  int `json:"like_state"`
	MediaCount int `json:"media_count"` // 收藏夹内容数量
}
type FavDetailUpper struct {
	MID      int64  `json:"mid"`      // 创建者mid
	Name     string `json:"name"`     // 创建者昵称
	Face     string `json:"face"`     // 创建者头像url
	Followed bool   `json:"followed"` // 是否已关注创建者 需登录
	// 会员类别
	//
	// 0：无
	//
	// 1：月大会员
	//
	// 2：年度及以上大会员
	VipType int `json:"vip_type"`
	// 会员开通状态(B站程序员有点粗心，打成statue了)
	//
	// 0：无
	//
	// 1：有
	VipStatue int `json:"vip_statue"`
}
type FavDetailCnt struct {
	Collect int   `json:"collect"`  // 收藏数
	Play    int64 `json:"play"`     // 收藏夹播放数
	ThumbUp int   `json:"thumb_up"` // 收藏夹点赞数
	Share   int   `json:"share"`    // 收藏夹分享数
}
type SpaceVideoCoin struct {
	videoBase
	Coins      int    `json:"coins"`       // 投币数量
	Time       int64  `json:"time"`        // 投币时间 时间戳
	IP         string `json:"ip"`          // 空
	InterVideo bool   `json:"inter_video"` // 是否为合作视频
}

// FavRes 收藏夹内容id
type FavRes struct {
	// 内容id
	//
	// 视频稿件：视频稿件avid
	//
	// 音频：音频auid
	//
	// 视频合集：视频合集id
	ID int64 `json:"id"`
	// 内容类型
	// 2：视频稿件
	//
	// 12：音频
	//
	// 21：视频合集
	Type int `json:"type"`
}
type FavResDetail struct {
	Info   *FavDetail           `json:"info"`   // 收藏夹元数据
	Medias []*FavResDetailMedia `json:"medias"` // 收藏夹内容
}
type FavResDetailMedia struct {
	// 内容id
	//
	// 视频稿件：视频稿件avid
	//
	// 音频：音频auid
	//
	// 视频合集：视频合集id
	ID int64 `json:"id"`
	// 内容类型
	//
	// 2：视频稿件
	//
	// 12：音频
	//
	// 21：视频合集
	Type     int                   `json:"type"`
	Title    string                `json:"title"`    // 标题
	Cover    string                `json:"cover"`    // 封面url
	Intro    string                `json:"intro"`    // 简介
	Page     int                   `json:"page"`     // 视频分P数
	Duration int64                 `json:"duration"` // 音频/视频时长
	Upper    *VideoOwner           `json:"upper"`    // UP主信息
	Attr     int                   `json:"attr"`     // 属性位
	CntInfo  *FavResDetailMediaCnt `json:"cnt_info"` // 状态数
	Link     string                `json:"link"`     // 跳转uri
	Ctime    int64                 `json:"ctime"`    // 投稿时间 时间戳
	Pubtime  int64                 `json:"pubtime"`  // 发布时间	时间戳
	FavTime  int64                 `json:"fav_time"` // 收藏时间 时间戳
	BVID     string                `json:"bvid"`     // 视频稿件bvid
}
type FavResDetailMediaCnt struct {
	Collect int   `json:"collect"` // 收藏数
	Play    int64 `json:"play"`    // 播放数
	Danmaku int   `json:"danmaku"` // 弹幕数
}

// ChanInfo 原频道仍能使用，视频列表为新版频道，还未实现相关接口
type ChanInfo struct {
	CID   int64  `json:"cid"`   // 频道id
	Count int    `json:"count"` // 频道内含视频数
	Cover string `json:"cover"` // 封面图片url
	Intro string `json:"intro"` // 简介 无则为空
	MID   int64  `json:"mid"`   // 创建用户mid
	Mtime int64  `json:"mtime"` // 创建时间	时间戳
	Name  string `json:"name"`  // 标题
}
type NavInfo struct {
	EmailVerified      int                    `json:"email_verified"`       // 是否验证邮箱地址 0:未验证 1:已验证
	Face               string                 `json:"face"`                 // 用户头像url
	LevelInfo          *NavInfoLevel          `json:"level_info"`           // 等级信息
	MID                int64                  `json:"mid"`                  // 用户mid
	MobileVerified     int                    `json:"mobile_verified"`      // 是否验证手机号 0:未验证 1:已验证
	Money              float64                `json:"money"`                // 拥有硬币数
	Moral              int                    `json:"moral"`                // 当前节操值 上限为70
	Official           *NavInfoOfficial       `json:"official"`             // 认证信息
	OfficialVerify     *NavInfoOfficialVerify `json:"officialVerify"`       // 认证信息2
	Pendant            *NavInfoPendant        `json:"pendant"`              // 头像框信息
	Scores             int                    `json:"scores"`               // 0 作用尚不明确
	Uname              string                 `json:"uname"`                // 用户昵称
	VipDueDate         int64                  `json:"vipDueDate"`           // 会员到期时间 毫秒 时间戳(东八区)
	VipStatus          int                    `json:"vipStatus"`            // 会员开通状态 0:无 1:有
	VipType            int                    `json:"vipType"`              // 会员类型 0:无 1:月度大会员 2:年度及以上大会员
	VipPayType         int                    `json:"vip_pay_type"`         // 会员开通状态	0:无 1:有
	VipThemeType       int                    `json:"vip_theme_type"`       // 0 作用尚不明确
	VipLabel           *NavInfoVipLabel       `json:"vip_label"`            // 会员标签
	VipAvatarSubscript int                    `json:"vip_avatar_subscript"` // 是否显示会员图标 0:不显示 1:显示
	VipNicknameColor   string                 `json:"vip_nickname_color"`   // 会员昵称颜色	颜色码 如#FFFFFF
	Wallet             *NavInfoWallet         `json:"wallet"`               // B币钱包信息
	HasShop            bool                   `json:"has_shop"`             // 是否拥有推广商品 false:无 true:有
	ShopURL            string                 `json:"shop_url"`             // 商品推广页面url
	AllowanceCount     int                    `json:"allowance_count"`      // 0 作用尚不明确
	AnswerStatus       int                    `json:"answer_status"`        // 0 作用尚不明确
}
type NavInfoLevel struct {
	CurrentLevel int `json:"current_level"` // 当前等级
	CurrentMin   int `json:"current_min"`   // 当前等级经验最低值
	CurrentExp   int `json:"current_exp"`   // 当前经验
	NextExp      int `json:"next_exp"`      // 升级下一等级需达到的经验
}
type NavInfoOfficial struct {
	// 认证类型
	//
	// 0:无
	//
	// 1 2 7:个人认证
	//
	// 3 4 5 6:机构认证
	Role  int    `json:"role"`
	Title string `json:"title"` // 认证信息 无为空
	Desc  string `json:"desc"`  // 认证备注 无为空
	Type  int    `json:"type"`  // 是否认证 -1:无 0:认证
}
type NavInfoOfficialVerify struct {
	Type int    `json:"type"` // 是否认证 -1:无 0:认证
	Desc string `json:"desc"` // 认证信息 无为空
}
type NavInfoPendant struct {
	PID    int64  // 挂件id
	Name   string // 挂件名称
	Image  string // 挂件图片url
	Expire int    // 0 作用尚不明确
}
type NavInfoVipLabel struct {
	Path string `json:"path"` // 空 作用尚不明确
	Text string `json:"text"` // 会员名称
	// 会员标签
	//
	// vip:大会员
	//
	// annual_vip:年度大会员
	//
	// ten_annual_vip:十年大会员
	//
	// hundred_annual_vip:百年大会员
	LabelTheme string `json:"label_theme"`
}
type NavInfoWallet struct {
	MID           int64   `json:"mid"`             // 登录用户mid
	BcoinBalance  float64 `json:"bcoin_balance"`   // 拥有B币数
	CouponBalance float64 `json:"coupon_balance"`  // 每月奖励B币数
	CouponDueTime int     `json:"coupon_due_time"` // 0 作用尚不明确
}
type MsgUnRead struct {
	At     int `json:"at"`      // 未读at数
	Chat   int `json:"chat"`    // 恒为0 作用尚不明确
	Like   int `json:"like"`    // 未读点赞数
	Reply  int `json:"reply"`   // 未读回复数
	SysMsg int `json:"sys_msg"` // 未读系统通知数
	Up     int `json:"up"`      // UP主助手信息数
}
type GeoInfo struct {
	Addr        string  `json:"addr"`         // 公网IP地址
	Country     string  `json:"country"`      // 国家/地区名
	Province    string  `json:"province"`     // 省/州 非必须存在项
	City        string  `json:"city"`         // 城市 非必须存在项
	Isp         string  `json:"isp"`          // 运营商名
	Latitude    float64 `json:"latitude"`     // 纬度
	Longitude   float64 `json:"longitude"`    // 经度
	ZoneID      int64   `json:"zone_id"`      // ip数据库id
	CountryCode int     `json:"country_code"` // 国家/地区代码
}
type VideoZone struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Desc string `json:"desc"`
}
type VideoSingleStat struct {
	AID        int64  `json:"aid"`        // 稿件avid
	BVID       string `json:"bvid"`       // 稿件bvid
	View       int64  `json:"view"`       // 播放次数
	Danmaku    int    `json:"danmaku"`    // 弹幕条数
	Reply      int    `json:"reply"`      // 评论条数
	Favorite   int    `json:"favorite"`   // 收藏人数
	Coin       int    `json:"coin"`       // 投币枚数
	Share      int    `json:"share"`      // 分享次数
	NowRank    int    `json:"now_rank"`   // 作用尚不明确
	HisRank    int    `json:"his_rank"`   // 历史最高排行
	Like       int    `json:"like"`       // 获赞次数
	Dislike    int    `json:"dislike"`    // 点踩次数,恒为0
	NoReprint  int    `json:"no_reprint"` // 禁止转载标志 0：无 1：禁止
	Copyright  int    `json:"copyright"`  // 版权标志 1：自制 2：转载
	ArgueMsg   string `json:"argue_msg"`  // 警告信息,默认为空
	Evaluation string `json:"evaluation"` // 视频评分,默认为空
}
type videoBase struct {
	AID       int64           `json:"aid"`       // 稿件avid
	BVID      string          `json:"bvid"`      // 稿件bvid
	Videos    int             `json:"videos"`    // 稿件分P总数
	TID       int             `json:"tid"`       // 分区TID
	Tname     string          `json:"tname"`     // 子分区名称
	Copyright int             `json:"copyright"` // 视频类型 (1:原创 2:转载)
	Pic       string          `json:"pic"`       // 稿件封面图片URL
	Title     string          `json:"title"`     // 稿件标题
	Pubdate   int64           `json:"pubdate"`   // 稿件发布时间 时间戳 时区为东八区(也就是说转换过来就是中国发布时间)
	Ctime     int64           `json:"ctime"`     // 用户投稿时间 时间戳 时区为东八区(也就是说转换过来就是中国发布时间)
	Desc      string          `json:"desc"`      // 视频简介
	State     int             `json:"state"`     // 视频状态
	Duration  int64           `json:"duration"`  // 稿件总时长(所有分P) 单位为秒
	Rights    *VideoRights    `json:"rights"`    // 视频属性标志
	Owner     *VideoOwner     `json:"owner"`     // 视频UP主信息
	Stat      *VideoStat      `json:"stat"`      // 视频状态数
	Dynamic   string          `json:"dynamic"`   // 视频同步发布的的动态的文字内容
	CID       int64           `json:"cid"`       // 视频1P cid
	Dimension *VideoDimension `json:"dimension"` // 视频1P分辨率
}
type VideoInfo struct {
	videoBase
	NoCache     bool           `json:"no_cache"`     // 恒为true 作用尚不明确
	Pages       []*VideoPage   `json:"pages"`        // 视频分P列表
	Subtitle    *VideoSubtitle `json:"subtitle"`     // 视频CC字幕信息
	Staff       []*VideoStaff  `json:"staff"`        // 合作成员列表 (非合作视频无此项)
	UserGarb    *VideoUserGarb `json:"user_garb"`    // 用户装扮信息
	DescV2      []*VideoDesc   `json:"desc_v2"`      // 新版视频简介
	Forward     int64          `json:"forward"`      // 撞车视频跳转avid (仅撞车视频存在此字段)
	MissionID   int64          `json:"mission_id"`   // 稿件参与的活动ID
	RedirectURL string         `json:"redirect_url"` // 重定向URL  仅番剧或影视视频存在此字段,用于番剧&影视的av/bv->ep
}
type VideoDesc struct {
	RawText string `json:"raw_text"` // 简介内容
	Type    int    `json:"type"`     // 未知
	BizID   int64  `json:"biz_id"`   // 未知
}
type VideoRights struct {
	BP            int `json:"bp"`              // 恒为0
	Elec          int `json:"elec"`            // 是否支持充电
	Download      int `json:"download"`        // 是否允许下载
	Movie         int `json:"movie"`           // 是否为电影
	Pay           int `json:"pay"`             // 是否PGC付费
	HD5           int `json:"hd5"`             // 是否有高码率
	NoReprint     int `json:"no_reprint"`      // 是否禁止转载
	Autoplay      int `json:"autoplay"`        // 是否自动播放
	UGCPay        int `json:"ugc_pay"`         //	是否UGC付费
	IsSteinGate   int `json:"is_stein_gate"`   // 是否为互动视频
	IsCooperation int `json:"is_cooperation"`  // 是否为联合投稿
	UGCPayPreview int `json:"ugc_pay_preview"` // 恒为0
	NoBackground  int `json:"no_background"`   // 恒为0
}
type VideoOwner struct {
	MID  int64  `json:"mid"`  // UP主mid
	Name string `json:"name"` // UP主昵称
	Face string `json:"face"` // UP主头像URL直链
}
type VideoStat struct {
	AID        int64  `json:"aid"`        // 稿件avid
	View       int64  `json:"view"`       // 播放次数
	Danmaku    int    `json:"danmaku"`    // 弹幕条数
	Reply      int    `json:"reply"`      // 评论条数
	Favorite   int    `json:"favorite"`   // 收藏人数
	Coin       int    `json:"coin"`       // 投币枚数
	Share      int    `json:"share"`      // 分享次数
	NowRank    int    `json:"now_rank"`   // 作用尚不明确
	HisRank    int    `json:"his_rank"`   // 历史最高排行
	Like       int    `json:"like"`       // 获赞次数
	Dislike    int    `json:"dislike"`    // 点踩次数,恒为0
	ArgueMsg   string `json:"argue_msg"`  // 警告信息,默认为空
	Evaluation string `json:"evaluation"` // 视频评分,默认为空
}
type VideoDimension struct {
	Width  int `json:"width"`  // 当前分P宽度
	Height int `json:"height"` // 当前分P高度
	Rotate int `json:"rotate"` // 是否将宽高对换 0：正常 1：对换
}
type VideoPage struct {
	CID       int64           `json:"cid"`       // 当前分P的CID
	Page      int             `json:"page"`      // 当前分P 在Pages中的id
	From      string          `json:"from"`      // 视频来源 vupload：普通上传(B站) hunan：芒果TV qq：腾讯
	Part      string          `json:"part"`      // 当前分P标题
	Duration  int64           `json:"duration"`  // 当前分P持续时间 单位为秒
	VID       string          `json:"vid"`       // 站外视频VID 仅站外视频有效
	Weblink   string          `json:"weblink"`   // 站外视频跳转URL 仅站外视频有效
	Dimension *VideoDimension `json:"dimension"` // 当前分P分辨率
}
type VideoSubtitle struct {
	AllowSubmit bool                 `json:"allow_submit"` // 是否允许提交字幕
	List        []*VideoSubtitleList `json:"list"`         // 字幕列表
}
type VideoSubtitleList struct {
	ID          int64                `json:"id"`           // 字幕ID
	Lan         string               `json:"lan"`          // 字幕语言
	LanDoc      string               `json:"lan_doc"`      // 字幕语言名称
	IsLock      bool                 `json:"is_lock"`      // 是否锁定
	AuthorMID   int64                `json:"author_mid"`   // 字幕上传者MID
	SubtitleURL string               `json:"subtitle_url"` // JSON格式字幕文件URL
	Author      *VideoSubtitleAuthor `json:"author"`       // 字幕上传者信息
}
type VideoSubtitleAuthor struct {
	MID           int64  `json:"mid"`             // 字幕上传者MID
	Name          string `json:"name"`            // 字幕上传者昵称
	Sex           string `json:"sex"`             // 字幕上传者性别 (男 女 保密)
	Face          string `json:"face"`            // 字幕上传者头像URL
	Sign          string `json:"sign"`            // 字幕上传者个性签名
	Rank          int    `json:"rank"`            // 恒为10000 作用尚不明确
	Birthday      int    `json:"birthday"`        // 恒为0 作用尚不明确
	IsFakeAccount int    `json:"is_fake_account"` // 恒为0 作用尚不明确
	IsDeleted     int    `json:"is_deleted"`      // 恒为0 作用尚不明确
}
type VideoStaff struct {
	MID      int64               `json:"mid"`      // 成员MID
	Title    string              `json:"title"`    // 成员名称
	Name     string              `json:"name"`     // 成员昵称
	Face     string              `json:"face"`     // 成员头像URL
	VIP      *VideoStaffVIP      `json:"vip"`      // 成员大会员状态
	Official *VideoStaffOfficial `json:"official"` //	成员认证信息
	Follower int                 `json:"follower"` // 成员粉丝数
}
type VideoStaffVIP struct {
	Type      int `json:"type"`       // 成员会员类型 (0:无 1:月会员 2:年会员)
	Status    int `json:"status"`     // 会员状态 (0:无 1:有)
	ThemeType int `json:"theme_type"` // 恒为0
}
type VideoStaffOfficial struct {
	// 成员认证级别
	//
	// 0为无
	//
	// 1 2 7为个人认证
	//
	// 3 4 5 6为机构认证
	Role  int    `json:"role"`
	Title string `json:"title"` // 成员认证名,Role无时该值为空
	Desc  string `json:"desc"`  // 成员认证备注,Role无时该值为空
	Type  int    `json:"type"`  // 成员认证类型 (-1:无 0:有)
}
type VideoUserGarb struct {
	URLImageAniCut string // 一串URL,未知
}
type VideoTag struct {
	TagID        int64          `json:"tag_id"`        // TAG ID
	TagName      string         `json:"tag_name"`      // TAG名称
	Cover        string         `json:"cover"`         // TAG图片URL
	HeadCover    string         `json:"head_cover"`    // TAG页面头图URL
	Content      string         `json:"content"`       // TAG介绍
	ShortContent string         `json:"short_content"` // TAG简介
	Type         int            `json:"type"`          // 未知
	State        int            `json:"state"`         // 恒为0
	Ctime        int64          `json:"ctime"`         // 创建时间 时间戳(已经为东八区)
	Count        *VideoTagCount `json:"count"`         // 状态数
	IsAtten      int            `json:"is_atten"`      // 是否关注 0:未关注 1:已关注 需要登录(Cookie) 未登录为0
	Likes        int            `json:"likes"`         // 恒为0 作用尚不明确
	Hates        int            `json:"hates"`         // 恒为0 作用尚不明确
	Attribute    int            `json:"attribute"`     // 恒为0 作用尚不明确
	Liked        int            `json:"liked"`         // 是否已经点赞 0：未点赞 1：已点赞 需要登录(Cookie) 未登录为0
	Hated        int            `json:"hated"`         // 是否已经点踩 0：未点踩 1：已点踩 需要登录(Cookie) 未登录为0
}
type VideoTagCount struct {
	View  int `json:"view"`  // 恒为0 作用尚不明确
	Use   int `json:"use"`   // 被使用次数
	Atten int `json:"atten"` // TAG关注数
}
type VideoPlayURLResult struct {
	From              string                `json:"from"`               // local 作用尚不明确
	Result            string                `json:"result"`             // suee 作用尚不明确
	Message           string                `json:"message"`            // 空 作用尚不明确
	Quality           int                   `json:"quality"`            // 当前的视频分辨率代码
	Format            string                `json:"format"`             // 当前请求的视频格式
	TimeLength        int64                 `json:"timelength"`         // 视频长度 单位为毫秒 不同分辨率/格式可能有略微差异
	AcceptFormat      string                `json:"accept_format"`      // 视频支持的全部格式 每项用,分隔
	AcceptDescription []string              `json:"accept_description"` // 视频支持的分辨率列表
	AcceptQuality     []int                 `json:"accept_quality"`     // 视频支持的分辨率代码列表
	VideoCodecid      int                   `json:"video_codecid"`      // 7 作用尚不明确
	SeekParam         string                `json:"seek_param"`         // start 作用尚不明确
	SeekType          string                `json:"seek_type"`          // ??? 作用尚不明确
	DURL              []*VideoPlayDURL      `json:"durl"`               // 视频分段	注：仅flv/mp4存在此项
	Dash              *VideoPlayURLDash     `json:"dash"`               // dash音视频流信息	注：仅dash存在此项
	SupportFormats    []*VideoPlayURLFormat `json:"support_formats"`    // 支持的分辨率的详细信息
}
type VideoPlayDURL struct {
	Order     int      `json:"order"`      // 视频分段序号 某些视频会分为多个片段（从1顺序增长）
	Length    int64    `json:"length"`     // 视频长度 单位为毫秒
	Size      int64    `json:"size"`       // 视频大小 单位为Byte
	Ahead     string   `json:"ahead"`      // 空 作用尚不明确
	Vhead     string   `json:"vhead"`      // 空 作用尚不明确
	URL       string   `json:"url"`        // 视频流url 注：url内容存在转义符 有效时间为120min
	BackupURL []string `json:"backup_url"` // 备用视频流
}
type VideoPlayURLFormat struct {
	Quality        int    `json:"quality"`         // 清晰度标识
	Format         string `json:"format"`          // 视频格式
	NewDescription string `json:"new_description"` // 新版清晰度描述
	DisplayDesc    string `json:"display_desc"`    // 显示名称
	Superscript    string `json:"superscript"`     // 角标？
}
type VideoPlayURLDash struct {
	Duration      int64                    `json:"duration"`        // 作用尚不明确
	MinBufferTime float64                  `json:"min_buffer_time"` // 1.5 作用尚不明确
	Video         []*VideoPlayURLDashMedia `json:"video"`           // 视频流信息
	Audio         []*VideoPlayURLDashMedia `json:"audio"`           // 音频流信息
}
type VideoPlayURLDashMedia struct {
	ID           int                       `json:"id"`             // 音视频清晰度代码
	BaseURL      string                    `json:"base_url"`       // 默认视频/音频流url 有效时间为120min
	BackupURL    []string                  `json:"backup_url"`     // 备用视频/音频流url
	Bandwidth    int64                     `json:"bandwidth"`      // 视频/音频所需最低带宽
	MimeType     string                    `json:"mime_type"`      // 视频/音频格式类型
	Codecs       string                    `json:"codecs"`         // 编码/音频类型
	Width        int                       `json:"width"`          // 视频宽度	单位为像素 仅视频有效
	Height       int                       `json:"height"`         // 视频高度 单位为像素 仅视频有效
	FrameRate    string                    `json:"frame_rate"`     // 视频帧率 仅视频有效
	Sar          string                    `json:"sar"`            //1:1	作用尚不明确
	StartWithSap int                       `json:"start_with_sap"` // 1	作用尚不明确
	SegmentBase  *VideoPlayURLDashMediaSeg `json:"segment_base"`   // ??? 作用尚不明确
	Codecid      int                       `json:"codecid"`        // 7 作用尚不明确
}
type VideoPlayURLDashMediaSeg struct {
	Initialization string `json:"initialization"` // ??? 作用尚不明确
	IndexRange     string `json:"index_range"`    // ??? 作用尚不明确
}

// VideoShot
//
// 快照的截取时间根据视频画面变化程度决定，各视频不相同
//
// 截取时间表的时间根据视频画面变化程度决定，各每个视频不相同
//
// 截取时间表的时间和快照一一对应，并按照从左到右 从上到下的顺序排布
type VideoShot struct {
	// bin格式截取时间表URL
	//
	// bin数据格式: https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/video/snapshot.md#bin%E6%A0%BC%E5%BC%8F%E6%88%AA%E5%8F%96%E6%97%B6%E9%97%B4%E8%A1%A8
	Pvdata   string   `json:"pvdata"`
	ImgXLen  int      `json:"img_x_len"`  // 每行图片数 一般为10
	ImgYLen  int      `json:"img_y_len"`  // 每列图片数 一般为10
	ImgXSize int      `json:"img_x_size"` // 每张图片长 一般为160
	ImgYSize int      `json:"img_y_size"` // 每张图片宽 一般为90
	Image    []string `json:"image"`      // 图片拼版URL 第一张拼版占满时延续第二张
	Index    []int    `json:"index"`      // json数组格式截取时间表 单位为秒
}
type SearchAll struct {
	SEID           string
	Page           int
	PageSize       int
	NumResults     int
	NumPages       int
	SuggestKeyword string
	RqtType        string
	CostTime       *SearchCostTime
	// ExpList 作用尚不明确
	EggHit         int
	PageInfo       *SearchPage
	TopTlist       *SearchTopTlist
	ShowColumn     int
	ShowModuleList []string
	Result         *SearchResult
}
type SearchCostTime struct {
	ParamsCheck         string `json:"params_check"`
	IllegalHandler      string `json:"illegal_handler"`
	AsResponseFormat    string `json:"as_response_format"`
	AsRequest           string `json:"as_request"`
	SaveCache           string `json:"save_cache"`
	DeserializeResponse string `json:"deserialize_response"`
	AsRequestFormat     string `json:"as_request_format"`
	Total               string `json:"total"`
	MainHandler         string `json:"main_handler"`
}
type SearchPage struct {
	PGC           *SearchPageInfo `json:"pgc"`            //
	LiveRoom      *SearchPageInfo `json:"live_room"`      // 直播数
	Photo         *SearchPageInfo `json:"photo"`          // 相簿数
	Topic         *SearchPageInfo `json:"topic"`          // 话题数
	Video         *SearchPageInfo `json:"video"`          // 视频数
	User          *SearchPageInfo `json:"user"`           //
	BiliUser      *SearchPageInfo `json:"bili_user"`      // 用户数
	MediaFT       *SearchPageInfo `json:"media_ft"`       // 电影数
	Article       *SearchPageInfo `json:"article"`        // 专栏数
	MediaBangumi  *SearchPageInfo `json:"media_bangumi"`  // 番剧数
	Special       *SearchPageInfo `json:"special"`        //
	OperationCard *SearchPageInfo `json:"operation_card"` //
	UpUser        *SearchPageInfo `json:"upuser"`         //
	Movie         *SearchPageInfo `json:"movie"`          //
	LiveAll       *SearchPageInfo `json:"live_all"`       //
	TV            *SearchPageInfo `json:"tv"`             //
	Live          *SearchPageInfo `json:"live"`           // 直播间数
	Bangumi       *SearchPageInfo `json:"bangumi"`        //
	Activity      *SearchPageInfo `json:"activity"`       // 活动数
	LiveMaster    *SearchPageInfo `json:"live_master"`    //
	LiveUser      *SearchPageInfo `json:"live_user"`      // 主播数
}
type SearchPageInfo struct {
	NumResults int `json:"numResults"` // 总计数量
	Total      int `json:"total"`      // 总计数量
	Pages      int `json:"pages"`      // 分页数量
}
type SearchTopTlist struct {
	PGC           int `json:"pgc"`            //
	LiveRoom      int `json:"live_room"`      // 直播数
	Photo         int `json:"photo"`          // 相簿数
	Topic         int `json:"topic"`          // 话题数
	Video         int `json:"video"`          // 视频数
	User          int `json:"user"`           //
	BiliUser      int `json:"bili_user"`      // 用户数
	MediaFT       int `json:"media_ft"`       // 电影数
	Article       int `json:"article"`        // 专栏数
	MediaBangumi  int `json:"media_bangumi"`  // 番剧数
	Special       int `json:"special"`        //
	Card          int `json:"card"`           //
	OperationCard int `json:"operation_card"` //
	UpUser        int `json:"upuser"`         //
	Movie         int `json:"movie"`          //
	LiveAll       int `json:"live_all"`       //
	TV            int `json:"tv"`             //
	Live          int `json:"live"`           // 直播间数 	//
	Bangumi       int `json:"bangumi"`        //
	Activity      int `json:"activity"`       // 活动数
	LiveMaster    int `json:"live_master"`    //
	LiveUser      int `json:"live_user"`      // 主播数
}

// SearchResult 在原搜索接口进行魔改，方便使用
type SearchResult struct {
	ResultType string
	Video      []*SearchResultVideo
	Media      []*SearchResultMedia
}
type SearchResultVideo struct {
	Type         string   `json:"type"`           // 结果类型 固定为video
	ID           int64    `json:"id"`             // 稿件avid
	Author       string   `json:"author"`         // UP主昵称
	MID          int64    `json:"mid"`            // UP主mid
	TypeID       string   `json:"typeid"`         // 视频分区tid
	TypeName     string   `json:"typename"`       // 视频子分区名
	ArcURL       string   `json:"arcurl"`         // 视频重定向URL
	AID          int64    `json:"aid"`            // 稿件avid
	BVID         string   `json:"bvid"`           // 稿件bvid
	Title        string   `json:"title"`          // 视频标题 关键字用xml标签<em class="keyword">标注
	Description  string   `json:"description"`    // 视频简介
	ArcRank      string   `json:"arcrank"`        // 恒为0 作用尚不明确
	Pic          string   `json:"pic"`            // 视频封面url
	Play         int64    `json:"play"`           // 视频播放量
	VideoReview  int      `json:"video_review"`   // 视频弹幕量
	Favorites    int      `json:"favorites"`      // 视频收藏数
	Tag          string   `json:"tag"`            // 视频TAG 每项TAG用,分隔
	Review       int      `json:"review"`         // 视频评论数
	PubDate      int64    `json:"pubdate"`        // 视频投稿时间 时间戳(东八区)
	SendDate     int64    `json:"senddate"`       // 视频发布时间 时间戳(东八区)
	Duration     string   `json:"duration"`       // 视频时长 格式: HH:MM
	BadgePay     bool     `json:"badgepay"`       // 恒为false 作用尚不明确
	HitColumns   []string `json:"hit_columns"`    // 关键字匹配类型
	ViewType     string   `json:"view_type"`      // 空 作用尚不明确
	IsPay        int      `json:"is_pay"`         // 空 作用尚不明确
	IsUnionVideo int      `json:"is_union_video"` // 是否为合作视频 0:否 1:是
	RankScore    int64    `json:"rank_score"`     // 结果排序量化值
	// RecTags      string NULL
	// NewRecTags   []string 空数组
}
type SearchResultMedia struct {
	Type           string                          `json:"type"`             // 结果类型 (media_bangumi:番剧 media_ft:影视)
	MediaID        int64                           `json:"media_id"`         // 剧集mdid
	SeasonID       int64                           `json:"season_id"`        // 剧集ssid
	Title          string                          `json:"title"`            // 剧集标题 关键字用xml标签<em class="keyword">标注
	OrgTitle       string                          `json:"org_title"`        // 剧集原名 关键字用xml标签<em class="keyword">标注 可为空
	Cover          string                          `json:"cover"`            // 剧集封面url
	MediaType      int                             `json:"media_type"`       // 剧集类型 (1:番剧 2:电影 3:纪录片 4:国创 5:电视剧 7:综艺)
	Areas          string                          `json:"areas"`            // 地区
	Styles         string                          `json:"styles"`           // 风格
	CV             string                          `json:"cv"`               // 声优
	Staff          string                          `json:"staff"`            // 制作组
	PlayState      int                             `json:"play_state"`       // 恒为0 作用尚不明确
	GotoURL        string                          `json:"goto_url"`         // 剧集重定向url
	Desc           string                          `json:"desc"`             // 简介
	Corner         int                             `json:"corner"`           // 角标有无 2：无 13：有
	PubTime        int64                           `json:"pub_time"`         // 开播时间 时间戳(东八区)
	MediaMode      int                             `json:"media_mode"`       // 恒为2 作用尚不明确
	IsAvid         bool                            `json:"is_avid"`          // 恒为false 作用尚不明确
	FixPubTimeStr  string                          `json:"fix_pub_time_str"` // 开播时间重写信息 优先级高于pubtime 可为空
	MediaScore     *SearchResultMediaScore         `json:"media_score"`      // 评分信息	有效时：obj 无效时：null
	HitColumns     []string                        `json:"hit_columns"`      // 关键字匹配类型 有效时：array 无效时：null
	AllNetName     string                          `json:"all_net_name"`     // 空 作用尚不明确
	AllNetIcon     string                          `json:"all_net_icon"`     // 空 作用尚不明确
	AllNetURL      string                          `json:"all_net_url"`      // 空 作用尚不明确
	AngleTitle     string                          `json:"angle_title"`      // 角标内容
	AngleColor     int                             `json:"angle_color"`      // 角标颜色 (0:红色 2:橙色)
	DisplayInfo    []*SearchResultMediaDisplayInfo `json:"display_info"`     // 剧集标志信息
	HitEpids       string                          `json:"hit_epids"`        // 关键字匹配分集标题的分集epid 多个用,分隔
	PgcSeasonID    int64                           `json:"pgc_season_id"`    // 剧集ssid
	SeasonType     int                             `json:"season_type"`      // 剧集类型 (1:番剧 2:电影 3:纪录片 4:国创 5:电视剧 7:综艺)
	SeasonTypeName string                          `json:"season_type_name"` // 剧集类型文字
	SelectionStyle string                          `json:"selection_style"`  // 分集选择按钮风格 horizontal:横排式 grid:按钮式
	EpSize         int                             `json:"ep_size"`          // 结果匹配的分集数
	URL            string                          `json:"url"`              // 剧集重定向url
	ButtonText     string                          `json:"button_text"`      // 观看按钮文字
	IsFollow       int                             `json:"is_follow"`        // 是否追番 需要登录(SESSDATA) 未登录则恒为0 (0:否 1:是)
	IsSelection    int                             `json:"is_selection"`     // 恒为1 作用尚不明确
	Eps            []*SearchResultEp               `json:"eps"`              // 结果匹配的分集信息
	Badges         []*SearchResultEpBadge          `json:"badges"`           // 剧集标志信息
}
type SearchResultMediaScore struct {
	UserCount int     `json:"user_count"` // 总计评分人数
	Score     float64 `json:"score"`      // 评分
}
type SearchResultMediaDisplayInfo struct {
	BgColorNight     string `json:"bg_color_night"`     // 夜间背景颜色 颜色码 例如:#BB5B76
	Text             string `json:"text"`               // 剧集标志 颜色码 例如:#BB5B76
	BorderColor      string `json:"border_color"`       // 背景颜色 颜色码 例如:#BB5B76
	BgStyle          int    `json:"bg_style"`           // 恒为1
	TextColor        string `json:"text_color"`         // 文字颜色 颜色码 例如:#BB5B76
	BgColor          string `json:"bg_color"`           // 背景颜色 颜色码 例如:#BB5B76
	TextColorNight   string `json:"text_color_night"`   // 夜间文字颜色 颜色码 例如:#BB5B76
	BorderColorNight string `json:"border_color_night"` // 夜间背景颜色 颜色码 例如:#BB5B76
}
type SearchResultEp struct {
	ID          int64                `json:"id"`           // 分集epid
	Cover       string               `json:"cover"`        // 分集封面url
	Title       string               `json:"title"`        // 完整标题
	URL         string               `json:"url"`          // 分集重定向url
	ReleaseDate string               `json:"release_date"` // 空
	Badges      *SearchResultEpBadge `json:"badges"`       // 分集标志
	IndexTitle  string               `json:"index_title"`  // 短标题
	LongTitle   string               `json:"long_title"`   // 单集标题
}
type SearchResultEpBadge struct {
	Text             string `json:"text"`               // 剧集标志	颜色码 例如:#BB5B76
	TextColor        string `json:"text_color"`         // 文字颜色	颜色码 例如:#BB5B76
	TextColorNight   string `json:"text_color_night"`   // 夜间文字颜色	颜色码 例如:#BB5B76
	BgColor          string `json:"bg_color"`           // 背景颜色	颜色码 例如:#BB5B76
	BgColorNight     string `json:"bg_color_night"`     // 夜间背景颜色	颜色码 例如:#BB5B76
	BorderColor      string `json:"border_color"`       // 空
	BorderColorNight string `json:"border_color_night"` // 空
	BgStyle          int    `json:"bg_style"`           // 恒为1
}
type DanmakuPostResult struct {
	Action  string `json:"action"`   // 空 作用尚不明确
	Dmid    uint64 `json:"dmid"`     // 弹幕dmid
	DmidStr string `json:"dmid_str"` // 弹幕dmid的字符串形式
	Visible bool   `json:"visible"`  // 作用尚不明确
}
type DanmakuCommandPostResult struct {
	// 指令
	//
	// UP主头像弹幕:#UP#
	//
	// 关联视频弹幕:#LINK#
	//
	// 视频内嵌引导关注按钮:#ATTENTION#
	Command  string          `json:"command"`
	Content  string          `json:"content"`  // 弹幕内容
	Extra    json.RawMessage `json:"extra"`    // JSON序列，具体请参考 https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/danmaku/action.md#%E5%8F%91%E9%80%81%E4%BA%92%E5%8A%A8%E5%BC%B9%E5%B9%95
	ID       uint64          `json:"id"`       // 弹幕dmid
	IDStr    string          `json:"idStr"`    // 弹幕dmid的字符串形式
	MID      int64           `json:"mid"`      // 用户mid
	OID      int64           `json:"oid"`      // 视频cid
	Progress int64           `json:"progress"` // 弹幕出现在视频内的时间
	// 互动弹幕类型
	//
	// 1:UP主头像弹幕
	//
	// 2:关联视频弹幕
	//
	// 5:视频内嵌引导关注按钮
	Type int `json:"type"`
}
type DanmakuGetLikesResult struct {
	Likes int `json:"likes"` // 点赞数
	// 当前账户是否点赞
	//
	// 0:未点赞
	// 1:已点赞
	//
	// 需要登录(Cookie或APP)
	// 未登录恒为0
	UserLike int    `json:"user_like"`
	IDStr    string `json:"id_str"` // 弹幕dmid的字符串形式
}

// DanmakuConfig 未启用的就传入空
type DanmakuConfig struct {
	DmSwitch     bool    `json:"dm_switch"`    // 弹幕开关
	BlockScroll  bool    `json:"blockscroll"`  // 屏蔽类型-滚动
	BlockTop     bool    `json:"blocktop"`     // 屏蔽类型-顶部
	BlockBottom  bool    `json:"blockbottom"`  // 屏蔽类型-底部
	BlockColor   bool    `json:"blockcolor"`   // 屏蔽类型-彩色
	BlockSpecial bool    `json:"blockspecial"` // 屏蔽类型-特殊
	AISwitch     bool    `json:"ai_switch"`    // 是否打开智能云屏蔽
	AILevel      int     `json:"ai_level"`     // 智能云屏蔽等级
	PreventShade bool    `json:"preventshade"` // 防挡弹幕（底部15%）
	DmMask       bool    `json:"dmask"`        // 智能防挡弹幕（人像蒙版）
	Opacity      float64 `json:"opacity"`      // 弹幕不透明度 区间：[0-1]
	// 弹幕显示区域
	//
	// 100：不重叠
	//
	// 75：3/4屏
	//
	// 50：半瓶
	//
	// 25：1/4屏
	//
	// 0：不限
	DmArea     int     `json:"dmarea"`
	SpeedPlus  float64 `json:"speedplus"`  // 弹幕速度 区间：[0.4-1.6]
	FontSize   float64 `json:"fontsize"`   // 字体大小 区间：[0.4-1.6]
	ScreenSync bool    `json:"screensync"` // 跟随屏幕缩放比例
	SpeedSync  bool    `json:"speedsync"`  // 根据播放倍速调整速度
	FontFamily string  `json:"fontfamily"` // 字体类型 未启用
	Bold       bool    `json:"bold"`       // 粗体 未启用
	FontBorder int     `json:"fontborder"` // 描边类型 0:重墨 1:描边 2:45°投影
	DrawType   string  `json:"drawType"`   // 渲染类型 未启用
}
type Emote struct {
	ID        int64  `json:"id"`         // 表情id
	PackageID int64  `json:"package_id"` // 表情包id
	Text      string `json:"text"`       // 表情转义符	颜文字时为该字串
	URL       string `json:"url"`        // 表情图片url	颜文字时为该字串
	Mtime     int64  `json:"mtime"`      // 创建时间	时间戳
	// 表情类型
	//
	// 1：普通
	//
	// 2：会员专属
	//
	// 3：购买所得
	//
	// 4：颜文字
	Type  int        `json:"type"`
	Attr  int        `json:"attr"`  // 作用尚不明确
	Meta  *EmoteMeta `json:"meta"`  // 属性信息
	Flags *EmoteFlag `json:"flags"` // 禁用标志
}
type EmoteMeta struct {
	Size    int      `json:"size"`    // 表情尺寸信息 1：小 2：大
	Alias   string   `json:"alias"`   // 简写名 无则无此项
	Suggest []string `json:"suggest"` // 文字对应的表情推荐
}
type EmoteFlag struct {
	Unlocked bool `json:"unlocked"` // true：启用 需要登录 否则恒为false
}
type EmotePack struct {
	ID    int64  `json:"id"`    // 表情包id
	Text  string `json:"text"`  // 表情包名称
	URL   string `json:"url"`   // 表情包标志图片url
	Mtime int64  `json:"mtime"` // 创建时间 时间戳
	// 表情包类型
	//
	// 1：普通
	//
	// 2：会员专属
	//
	// 3：购买所得
	//
	// 4：颜文字
	Type  int            `json:"type"`
	Attr  int            `json:"attr"`  // 作用尚不明确
	Meta  *EmotePackMeta `json:"meta"`  // 属性信息
	Emote []*Emote       `json:"emote"` // 表情列表
	Flags *EmotePackFlag `json:"flags"` // 是否添加标志
}
type EmotePackMeta struct {
	Size    int    `json:"size"`     // 表情尺寸信息	1：小 2：大
	ItemID  int64  `json:"item_id"`  // 购买物品id
	ItemURL string `json:"item_url"` // 购买物品页面url 无则无此项
}
type EmotePackFlag struct {
	// 是否已添加
	//
	// true：已添加
	//
	// false：未添加
	//
	// 需要登录（SESSDATA）
	// 否则恒为false
	Added bool `json:"added"`
}
type AudioInfo struct {
	ID         int64          `json:"id"`         // 音频auid
	UID        int64          `json:"uid"`        // UP主mid
	Uname      string         `json:"uname"`      // UP主昵称
	Author     string         `json:"author"`     // 作者名
	Title      string         `json:"title"`      // 歌曲标题
	Cover      string         `json:"cover"`      // 封面图片url
	Intro      string         `json:"intro"`      // 歌曲简介
	Lyric      string         `json:"lyric"`      // lrc歌词url
	CrType     int            `json:"crtype"`     // 1 作用尚不明确
	Duration   int64          `json:"duration"`   // 歌曲时间长度 单位为秒
	PassTime   int64          `json:"passtime"`   // 歌曲发布时间 时间戳
	CurTime    int64          `json:"curtime"`    // 当前请求时间	时间戳
	AID        int64          `json:"aid"`        // 关联稿件avid 无为0
	BVID       string         `json:"bvid"`       // 关联稿件bvid 无为空
	CID        int64          `json:"cid"`        // 关联视频cid 无为0
	MsID       int            `json:"msid"`       // 0 作用尚不明确
	Attr       int            `json:"attr"`       // 0 作用尚不明确
	Limit      int            `json:"limit"`      // 0 作用尚不明确
	ActivityID int            `json:"activityId"` // 0 作用尚不明确
	LimitDesc  string         `json:"limitdesc"`  // 0 作用尚不明确
	Ctime      int64          `json:"ctime"`      // 0 作用尚不明确
	Statistic  *AudioInfoStat `json:"statistic"`  // 状态数
	VipInfo    *AudioInfoVip  `json:"vipInfo"`    // UP主会员状态
	CollectIDs []int64        `json:"collectIds"` // 歌曲所在的收藏夹mlids 需要登录(SESSDATA)
	CoinNum    int            `json:"coin_num"`   // 投币数
}
type AudioInfoStat struct {
	SID     int64 `json:"sid"`     // 音频auid
	Play    int64 `json:"play"`    // 播放次数
	Collect int   `json:"collect"` // 收藏数
	Comment int   `json:"comment"` // 评论数
	Share   int   `json:"share"`   // 分享数
}
type AudioInfoVip struct {
	// 会员类型
	//
	// 0：无
	//
	// 1：月会员
	//
	// 2：年会员
	Type       int   `json:"type"`
	Status     int   `json:"status"`       // 会员状态 0：无 1：有
	DueDate    int64 `json:"due_date"`     // 会员到期时间 时间戳 毫秒
	VipPayType int   `json:"vip_pay_type"` // 会员开通状态 0：无 1：有
}
type AudioTag struct {
	Type    string `json:"type"`    // song 作用尚不明确
	Subtype int    `json:"subtype"` // ？？？ 作用尚不明确
	Key     int64  `json:"key"`     // TAG id？？ 作用尚不明确
	Info    string `json:"info"`    // TAG名
}
type AudioMember struct {
	List []*AudioMemberList `json:"list"` // 成员列表
	// 成员类型代码
	//
	// 1：歌手
	//
	// 2：作词
	//
	// 3：作曲
	//
	// 4：编曲
	//
	// 5：后期/混音
	//
	// 7：封面制作
	//
	// 8：音源
	//
	// 9：调音
	//
	// 10：演奏
	//
	// 11：乐器
	//
	// 127：UP主
	Type int `json:"type"`
}
type AudioMemberList struct {
	MID      int64  `json:"mid"`       // 0 作用尚不明确
	Name     string `json:"name"`      // 成员名
	MemberID int64  `json:"member_id"` // 成员id？？反正不是mid 作用尚不明确
}
type AudioMyFavLists struct {
	CurPage   int             `json:"curPage"`   // 当前页码
	PageCount int             `json:"pageCount"` // 总计页数
	TotalSize int             `json:"totalSize"` // 总计收藏夹数
	PageSize  int             `json:"pageSize"`  // 当前页面项数
	Data      []*AudioFavList `json:"data"`      // 歌单列表
}
type AudioFavList struct {
	ID        int64             `json:"id"`        // 音频收藏夹mlid,歌单url里显示的是这个
	UID       int64             `json:"uid"`       // 创建用户mid
	Uname     string            `json:"uname"`     // 创建用户昵称
	Title     string            `json:"title"`     // 歌单标题
	Type      int               `json:"type"`      // 收藏夹属性 0：普通收藏夹 1：默认收藏夹
	Published int               `json:"published"` // 是否公开 0：不公开 1：公开
	Cover     string            `json:"cover"`     // 歌单封面图片url
	Ctime     int64             `json:"ctime"`     // 歌单创建时间 时间戳
	Song      int               `json:"song"`      // 歌单中的音乐数量
	Desc      string            `json:"desc"`      // 歌单备注信息
	Sids      []int64           `json:"sids"`      // 歌单中的音乐的auid
	MenuID    int64             `json:"menuId"`    // 音频收藏夹对应的歌单amid
	Statistic *AudioFavListStat `json:"statistic"` // 歌单状态数信息
}
type AudioFavListStat struct {
	SID     int64 `json:"sid"`     // 音频收藏夹对应的歌单amid
	Play    int64 `json:"play"`    // 播放次数
	Collect int   `json:"collect"` // 收藏数
	Comment int   `json:"comment"` // 评论数
	Share   int   `json:"share"`   // 分享数
}
type AudioPlayURL struct {
	SID int64 `json:"sid"` // 音频auid
	// 音质标识
	//
	// -1：试听片段（192K）
	//
	// 0：128K
	//
	// 1：192K
	//
	// 2：320K
	//
	// 3：FLAC
	Type      int               `json:"type"`
	Info      string            `json:"info"`      // 空	作用尚不明确
	Timeout   int64             `json:"timeout"`   // 有效时长 单位为秒 一般为3h
	Size      int64             `json:"size"`      // 文件大小	单位为字节 当type为-1时size为0
	CDNs      []string          `json:"cdns"`      // 音频流url
	Qualities []*AudioPlayURLQn `json:"qualities"` // 音质列表
	Title     string            `json:"title"`     // 音频标题
	Cover     string            `json:"cover"`     // 音频封面url
}
type AudioPlayURLQn struct {
	Type        int    `json:"type"`        // 音质代码
	Desc        string `json:"desc"`        // 音质名称
	Size        int64  `json:"size"`        // 该音质的文件大小 单位为字节
	Bps         string `json:"bps"`         // 比特率标签
	Tag         string `json:"tag"`         // 音质标签
	Require     int    `json:"require"`     // 是否需要会员权限 0：不需要 1：需要
	RequireDesc string `json:"requiredesc"` // 会员权限标签
}
type ChargeBpResult struct {
	MID     int64  `json:"mid"`      // 本用户mid
	UpMID   int64  `json:"up_mid"`   // 目标用户mid
	OrderNo string `json:"order_no"` // 订单号 用于添加充电留言
	BpNum   string `json:"bp_num"`   // 充电B币数？(不知道按B币算还是换算成贝壳) 不知为何返回类型为string
	Exp     int    `json:"exp"`      // 获得经验数
	// 返回结果
	//
	// 4：成功
	//
	// -2：低于充电下限
	//
	// -4：B币不足
	Status int    `json:"status"`
	Msg    string `json:"msg"` // 错误信息 默认为空
}
type ChargeSpaceList struct {
	DisplayNum int           `json:"display_num"` // 0 作用尚不明确
	Count      int           `json:"count"`       // 本月充电人数
	TotalCount int           `json:"total_count"` // 总计充电人数
	List       []*ChargeItem `json:"list"`        // 本月充电用户列表
}
type ChargeItem struct {
	MID        int64          `json:"mid"`         // 充电对象mid
	PayMID     int64          `json:"pay_mid"`     // 充电用户mid
	Rank       int            `json:"rank"`        // 充电用户排名 取决于充电的多少
	Uname      string         `json:"uname"`       // 充电用户昵称
	Avatar     string         `json:"avatar"`      // 充电用户头像url
	Message    string         `json:"message"`     // 充电留言 无为空
	MsgDeleted int            `json:"msg_deleted"` // 0 作用尚不明确
	VipInfo    *ChargeItemVip `json:"vip_info"`    // 充电用户会员信息
	TrendType  int            `json:"trend_type"`  // 0 作用尚不明确
}
type ChargeItemVip struct {
	// 大会员类型
	//
	// 0：无
	//
	// 1：月会员
	//
	// 2：年会员
	VipType    int `json:"vipType"`
	VipDueMsec int `json:"vipDueMsec"` // 0 作用尚不明确
	VipStatus  int `json:"vipStatus"`  // 大会员状态 0：无 1：有
}
type ChargeVideoList struct {
	ShowInfo   *ChargeVideoShow `json:"show_info"`   // 展示选项
	AvCount    int              `json:"av_count"`    // 目标视频充电人数
	Count      int              `json:"count"`       // 本月充电人数
	TotalCount int              `json:"total_count"` // 总计充电人数
	SpecialDay int              `json:"special_day"` // 0 作用尚不明确
	DisplayNum int              `json:"display_num"` // 0 作用尚不明确
	List       []*ChargeItem    `json:"list"`        // 本月充电用户列表
}
type ChargeVideoShow struct {
	Show  bool `json:"show"`  // 是否展示视频充电鸣谢名单 false：不展示 true：展示
	State int  `json:"state"` // 0
}
type ChargeCreateQrCode struct {
	QrCodeURL string `json:"qr_code_url"` // 支付二维码生成内容 存在转义
	QrToken   string `json:"qr_token"`    // 扫码秘钥
	Exp       int    `json:"exp"`         // 获得经验数
}
type ChargeQrCodeStatus struct {
	QrToken string `json:"qr_token"` // 扫码秘钥
	OrderNo string `json:"order_no"` // 留言token 未成功则无此项 用于添加充电留言
	MID     int64  `json:"mid"`      // 当前用户mid
	// 状态值 若秘钥错误则无此项
	//
	// 1：已支付
	//
	// 2：未扫描
	//
	// 3：未确认
	Status int `json:"status"`
}
type LiveRoomInfoByMID struct {
	RoomStatus    int    `json:"roomStatus"`     // 直播间状态 0：无房间 1：有房间
	RoundStatus   int    `json:"roundStatus"`    // 轮播状态 0：未轮播 1：轮播
	LiveStatus    int    `json:"liveStatus"`     // 直播状态 0：未开播 1：直播中
	URL           string `json:"url"`            // 直播间网页url
	Title         string `json:"title"`          // 直播间标题
	Cover         string `json:"cover"`          // 直播间封面url
	Online        int    `json:"online"`         // 直播间人气 值为上次直播时刷新
	RoomID        int    `json:"roomid"`         // 直播间id(真实ID)
	BroadcastType int    `json:"broadcast_type"` // 0
	OnlineHidden  int    `json:"online_hidden"`  // 已废弃
}
type FollowingsDetail struct {
	ReVersion int               `json:"re_version"` // 0
	Total     int               `json:"total"`      // 关注总数
	List      []*FollowingsItem `json:"list"`       // 关注详细信息
}

type FollowingsItem struct {
	MID          int64    `json:"mid"`       // mid
	Attribute    int      `json:"attribute"` // 属性值
	Mtime        int64    `json:"mtime"`     // 关注时间
	Tag          []string `json:"tag,omitempty"`
	Special      int      `json:"special"`
	ContractInfo *struct {
		IsContractor bool `json:"is_contractor"`
		TS           int  `json:"ts"`
		IsContract   bool `json:"is_contract"`
		UserAttr     int  `json:"user_attr"`
	} `json:"contract_info"`
	Uname          string `json:"uname"` // 昵称
	Face           string `json:"face"`  // 头像
	Sign           string `json:"sign"`  // 个人签名
	OfficialVerify *struct {
		Type int    `json:"type"`
		Desc string `json:"desc"`
	} `json:"official_verify"` // 认证信息
	Vip *struct {
		VipType       int    `json:"vipType"`
		VipDueDate    int64  `json:"vipDueDate"`
		DueRemark     string `json:"dueRemark"`
		AccessStatus  int    `json:"accessStatus"`
		VipStatus     int    `json:"vipStatus"`
		VipStatusWarn string `json:"vipStatusWarn"`
		ThemeType     int    `json:"themeType"`
		Label         *struct {
			Path        string `json:"path"`
			Text        string `json:"text"`
			LabelTheme  string `json:"label_theme"`
			TextColor   string `json:"text_color"`
			BgStyle     int    `json:"bg_style"`
			BgColor     string `json:"bg_color"`
			BorderColor string `json:"border_color"`
		} `json:"label"`
		AvatarSubscript    int    `json:"avatar_subscript"`
		NicknameColor      string `json:"nickname_color"`
		AvatarSubscriptURL string `json:"avatar_subscript_url"`
	} `json:"vip"` // VIP信息
}
type dynaCtrl struct {
	Location int    `json:"location"`
	Type     int    `json:"type"`
	Length   int    `json:"length"`
	Data     string `json:"data"`
}
type FileUpload struct {
	Field string
	Name  string
	File  io.Reader
}
type DynaUploadPic struct {
	ImageURL    string `json:"image_url"`
	ImageWidth  int    `json:"image_width"`
	ImageHeight int    `json:"image_height"`
}
type dynaPic struct {
	ImgSrc      string `json:"img_src"`
	ImageWidth  int    `json:"img_width"`
	ImageHeight int    `json:"img_height"`
}
type dynaDraft struct {
	Biz         int    `json:"biz"`
	Category    int    `json:"category"`
	Type        int    `json:"type"`
	Pictures    string `json:"pictures"`
	Description string `json:"description"`
	From        string `json:"from"`
	Content     string `json:"content"`
	AtUIDs      string `json:"at_uids"`
	AtControl   string `json:"at_control"`
}
type DynaGetDraft struct {
	Drafts []*DynaDraft `json:"drafts"`
}
type DynaDraft struct {
	DraftID       int64           `json:"draft_id"`       // 定时发布ID
	UID           int64           `json:"uid"`            // mid
	Type          int             `json:"type"`           // 未知
	PublishTime   int64           `json:"publish_time"`   // 指定发布时间
	Request       json.RawMessage `json:"request"`        // 动态信息，不同动态类型内容不同，请根据需要自行提取
	UpdateTime    int64           `json:"update_time"`    // 动态更新时间
	PublishStatus int             `json:"publish_status"` // 发布状态 0:未发布 3:错误?
	ErrorCode     int             `json:"error_code"`     // 动态错误码 0:无错误 500003:系统错误
	ErrorMsg      string          `json:"error_msg"`      // 动态错误描述
	UserProfile   *struct {
		Info *struct {
			UID   int64  `json:"uid"`
			Uname string `json:"uname"`
			Face  string `json:"face"`
		} `json:"info"`
		Card *struct {
			OfficialVerify *struct {
				Type int    `json:"type"`
				Desc string `json:"desc"`
			} `json:"official_verify"`
		} `json:"card"`
		Vip *struct {
			VipType    int `json:"vipType"`
			VipDueDate int `json:"vipDueDate"`
			VipStatus  int `json:"vipStatus"`
			ThemeType  int `json:"themeType"`
			Label      *struct {
				Path       string `json:"path"`
				Text       string `json:"text"`
				LabelTheme string `json:"label_theme"`
			} `json:"label"`
			AvatarSubscript int    `json:"avatar_subscript"`
			NicknameColor   string `json:"nickname_color"`
		} `json:"vip"`
		Pendant *struct {
			PID               int64  `json:"pid"`
			Name              string `json:"name"`
			Image             string `json:"image"`
			Expire            int    `json:"expire"`
			ImageEnhance      string `json:"image_enhance"`
			ImageEnhanceFrame string `json:"image_enhance_frame"`
		} `json:"pendant"`
		Rank      string `json:"rank"`
		Sign      string `json:"sign"`
		LevelInfo *struct {
			CurrentLevel int `json:"current_level"`
		} `json:"level_info"`
	} `json:"user_profile"` // 动态作者各种信息
}

type LiveRoomInfoByID struct {
	RoomID          int64 `json:"room_id"`      // 真实直播间ID
	ShortID         int   `json:"short_id"`     // 短号
	UID             int64 `json:"uid"`          // 主播mid
	NeedP2P         int   `json:"need_p2p"`     // 需要P2P
	IsHidden        bool  `json:"is_hidden"`    // 直播间是否隐藏
	IsLocked        bool  `json:"is_locked"`    // 直播间是否被封锁
	IsPortrait      bool  `json:"is_portrait"`  // 是否为竖屏直播间
	LiveStatus      int   `json:"live_status"`  // 0:未开播 1:开播
	HiddenTill      int64 `json:"hidden_till"`  // 隐藏截止时间戳?
	LockTill        int64 `json:"lock_till"`    // 封锁截止时间戳?
	Encrypted       bool  `json:"encrypted"`    // 直播间是否加密
	PwdVerified     bool  `json:"pwd_verified"` // 直播间是否需要密码验证
	LiveTime        int64 `json:"live_time"`    // 开播时间,-1为未开播
	RoomShield      int   `json:"room_shield"`
	IsSp            int   `json:"is_sp"`
	SpecialType     int   `json:"special_type"`
	AllSpecialTypes []int `json:"all_special_types"`
}
type LiveWsConf struct {
	RefreshRowFactor float64 `json:"refresh_row_factor"`
	RefreshRate      int     `json:"refresh_rate"`
	MaxDelay         int     `json:"max_delay"`
	Port             int     `json:"port"`
	Host             string  `json:"host"`
	HostServerList   []*struct {
		Host    string `json:"host"`
		Port    int    `json:"port"`
		WssPort int    `json:"wss_port"`
		WsPort  int    `json:"ws_port"`
	} `json:"host_server_list"`
	ServerList []*struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server_list"`
	Token string `json:"token"`
}
type LiveAreaInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	List []*struct {
		ID              string `json:"id"`
		ParentID        string `json:"parent_id"`
		OldAreaID       string `json:"old_area_id"`
		Name            string `json:"name"`
		ActID           string `json:"act_id"`
		PkStatus        string `json:"pk_status"`
		HotStatus       int    `json:"hot_status"`
		LockStatus      string `json:"lock_status"`
		Pic             string `json:"pic"`
		ComplexAreaName string `json:"complex_area_name"`
		ParentName      string `json:"parent_name"`
		AreaType        int    `json:"area_type"`
		CateID          string `json:"cate_id,omitempty"`
	} `json:"list"`
}
type LiveGuardList struct {
	Info *struct {
		Num              int `json:"num"`  // 大航海总数
		Page             int `json:"page"` // 总页数
		Now              int `json:"now"`  // 该次请求的页数
		AchievementLevel int `json:"achievement_level"`
	} `json:"info"`
	List []*struct {
		UID           int64  `json:"uid"`
		RUID          int64  `json:"ruid"` // 主播mid
		Rank          int    `json:"rank"` // 在该数组中的排名
		Username      string `json:"username"`
		Face          string `json:"face"`
		IsAlive       int    `json:"is_alive"`
		GuardLevel    int    `json:"guard_level"` // 1:总督 2:提督 3:舰长
		GuardSubLevel int    `json:"guard_sub_level"`
	} `json:"list"`
	Top3 []*struct {
		UID           int    `json:"uid"`
		RUID          int    `json:"ruid"` // 主播mid
		Rank          int    `json:"rank"` // 在该数组中的排名
		Username      string `json:"username"`
		Face          string `json:"face"`
		IsAlive       int    `json:"is_alive"`
		GuardLevel    int    `json:"guard_level"` // 1:总督 2:提督 3:舰长
		GuardSubLevel int    `json:"guard_sub_level"`
	} `json:"top3"`
}
type LiveMedalRank struct {
	Medal struct {
		Status int `json:"status"`
	} `json:"medal"`
	List []*struct {
		UID              int64  `json:"uid"`        // mid
		Uname            string `json:"uname"`      // 昵称
		Face             string `json:"face"`       // 头像url
		Rank             int    `json:"rank"`       // 排名
		MedalName        string `json:"medal_name"` // 勋章名字
		Level            int    `json:"level"`      // 勋章等级
		Color            int64  `json:"color"`      // 勋章颜色
		TargetID         int64  `json:"target_id"`  // 主播mid
		Special          string `json:"special"`
		IsSelf           int    `json:"isSelf"`
		GuardLevel       int    `json:"guard_level"` // 1:总督 2:提督 3:舰长
		MedalColorStart  int64  `json:"medal_color_start"`
		MedalColorEnd    int64  `json:"medal_color_end"`
		MedalColorBorder int64  `json:"medal_color_border"`
		IsLighted        int    `json:"is_lighted"`
	} `json:"list"`
}
type LivePlayURL struct {
	CurrentQn          int `json:"current_qn"`
	QualityDescription []*struct {
		Qn   int    `json:"qn"`
		Desc string `json:"desc"`
	} `json:"quality_description"` // 清晰度列表
	DURL []*struct {
		URL        string `json:"url"`
		Length     int    `json:"length"`
		Order      int    `json:"order"`
		StreamType int    `json:"stream_type"`
		PTag       int    `json:"ptag"`
		P2PType    int    `json:"p2p_type"`
	} `json:"durl"`
	IsDashAuto bool `json:"is_dash_auto"`
}
type LiveAllGiftInfo struct {
	List []*struct {
		ID                int    `json:"id"`
		Name              string `json:"name"`
		Price             int    `json:"price"`
		Type              int    `json:"type"`
		CoinType          string `json:"coin_type"`
		BagGift           int    `json:"bag_gift"`
		Effect            int    `json:"effect"`
		CornerMark        string `json:"corner_mark"`
		CornerBackground  string `json:"corner_background"`
		Broadcast         int    `json:"broadcast"`
		Draw              int    `json:"draw"`
		StayTime          int    `json:"stay_time"`
		AnimationFrameNum int    `json:"animation_frame_num"`
		Desc              string `json:"desc"`
		Rule              string `json:"rule"`
		Rights            string `json:"rights"`
		PrivilegeRequired int    `json:"privilege_required"`
		CountMap          []*struct {
			Num            int    `json:"num"`
			Text           string `json:"text"`
			WebSvga        string `json:"web_svga"`
			VerticalSvga   string `json:"vertical_svga"`
			HorizontalSvga string `json:"horizontal_svga"`
			SpecialColor   string `json:"special_color"`
			EffectID       int    `json:"effect_id"`
		} `json:"count_map"`
		ImgBasic             string `json:"img_basic"`
		ImgDynamic           string `json:"img_dynamic"`
		FrameAnimation       string `json:"frame_animation"`
		GIF                  string `json:"gif"`
		Webp                 string `json:"webp"`
		FullScWeb            string `json:"full_sc_web"`
		FullScHorizontal     string `json:"full_sc_horizontal"`
		FullScVertical       string `json:"full_sc_vertical"`
		FullScHorizontalSvga string `json:"full_sc_horizontal_svga"`
		FullScVerticalSvga   string `json:"full_sc_vertical_svga"`
		BulletHead           string `json:"bullet_head"`
		BulletTail           string `json:"bullet_tail"`
		LimitInterval        int    `json:"limit_interval"`
		BindRUID             int    `json:"bind_ruid"`
		BindRoomID           int    `json:"bind_roomid"`
		GiftType             int    `json:"gift_type"`
		ComboResourcesID     int    `json:"combo_resources_id"`
		MaxSendLimit         int    `json:"max_send_limit"`
		Weight               int    `json:"weight"`
		GoodsID              int    `json:"goods_id"`
		HasImagedGift        int    `json:"has_imaged_gift"`
		LeftCornerText       string `json:"left_corner_text"`
		LeftCornerBackground string `json:"left_corner_background"`
		GiftBanner           struct {
			AppPic         string `json:"app_pic"`
			WebPic         string `json:"web_pic"`
			LeftText       string `json:"left_text"`
			LeftColor      string `json:"left_color"`
			ButtonText     string `json:"button_text"`
			ButtonColor    string `json:"button_color"`
			ButtonPicColor string `json:"button_pic_color"`
			JumpURL        string `json:"jump_url"`
			JumpTo         int    `json:"jump_to"`
			WebPicURL      string `json:"web_pic_url"`
			WebJumpURL     string `json:"web_jump_url"`
		} `json:"gift_banner"`
		DiyCountMap int `json:"diy_count_map"`
		EffectID    int `json:"effect_id"`
	} `json:"list"`
	ComboResources []*struct {
		ComboResourcesId int    `json:"combo_resources_id"`
		ImgOne           string `json:"img_one"`
		ImgTwo           string `json:"img_two"`
		ImgThree         string `json:"img_three"`
		ImgFour          string `json:"img_four"`
		ColorOne         string `json:"color_one"`
		ColorTwo         string `json:"color_two"`
	} `json:"combo_resources"`
	GuardResources []*struct {
		Level int    `json:"level"`
		Img   string `json:"img"`
		Name  string `json:"name"`
	} `json:"guard_resources"`
}
type Comment struct {
	RPID int64 `json:"rpid"` // 评论rpid
	OID  int64 `json:"oid"`  // 评论区对象id
	// 评论区类型代码
	//
	// https://github.com/SocialSisterYi/bilibili-API-collect/tree/master/comment#%E8%AF%84%E8%AE%BA%E5%8C%BA%E7%B1%BB%E5%9E%8B%E4%BB%A3%E7%A0%81
	Type int   `json:"type"`
	MID  int64 `json:"mid"` // 发送者mid
	// 根评论rpid
	//
	// 若为一级评论则为0
	//
	// 大于一级评论则为根评论id
	Root int64 `json:"root"`
	// 回复父评论rpid
	//
	// 若为一级评论则为0
	//
	// 若为二级评论则为根评论id
	//
	// 大于二级评论为上一级评论id
	Parent int64 `json:"parent"`
	// 回复对方rpid
	//
	// 若为一级评论则为0
	//
	// 若为二级评论则为该评论id
	//
	// 大于二级评论为上一级评论id
	Dialog    int64  `json:"dialog"`
	Count     int    `json:"count"`           // 二级评论条数
	Rcount    int    `json:"rcount"`          // 回复评论条数
	Floor     int    `json:"floor,omitempty"` // 评论楼层号 若不支持楼层则无此项
	State     int    `json:"state"`           // 0 作用尚不明确
	FansGrade int    `json:"fansgrade"`       // 是否具有粉丝标签	0：无 1：有
	Attr      int    `json:"attr"`            // 属性位
	Ctime     int64  `json:"ctime"`           // 评论发送时间 时间戳
	RpidStr   string `json:"rpid_str"`        // 评论rpid	字串格式
	RootStr   string `json:"root_str"`        // 根评论rpid 字串格式
	ParentStr string `json:"parent_str"`      // 回复父评论rpid 字串格式
	Like      int    `json:"like"`            // 评论获赞数
	// 当前用户操作状态	需要登录(Cookie或APP)
	//
	// 否则恒为0
	//
	// 0：无
	//
	// 1：已点赞
	//
	// 2：已点踩
	Action  int           `json:"action"`
	Member  CommentMember `json:"member"` // 评论发送者信息
	Content struct {
		Message string `json:"message"` // 评论内容
		// 评论发送端
		//
		// 1：web端
		//
		// 2：安卓客户端
		//
		// 3：ios客户端
		//
		// 4：wp客户端
		Plat    int              `json:"plat"`
		Device  string           `json:"device"`  // 评论发送平台设备
		Members []*CommentMember `json:"members"` // at到的用户信息
		Emote   map[string]struct {
			ID        int64  `json:"id"`
			PackageID int64  `json:"package_id"`
			Status    int    `json:"status"`
			Type      int    `json:"type"`
			Attr      int    `json:"attr"`
			Text      string `json:"text"`
			URL       string `json:"url"`
			Meta      struct {
				Size  int    `json:"size"`
				Alias string `json:"alias"`
			} `json:"meta"`
			Mtime int64 `json:"mtime"`
		} `json:"emote"` // 需要渲染的表情转义	评论内容无表情则无此项
		JumpURL map[string]struct {
			Title          string `json:"title"` // 标题
			State          int    `json:"state"`
			PrefixIcon     string `json:"prefixIcon"` // 图标url
			AppURLSchema   string `json:"appUrlSchema"`
			AppName        string `json:"appName"`
			AppPackageName string `json:"appPackageName"`
			ClickReport    string `json:"clickReport"`
		} `json:"jump_url"` // 需要高亮的超链转义
		MaxLine int `json:"max_line"` // 6	收起最大行数
	} `json:"content"` // 评论信息
	Replies []*Comment `json:"replies"` // 评论回复条目预览 仅嵌套一层 否则为null
	Assist  int        `json:"assist"`
	Folder  struct {
		HasFolded bool   `json:"has_folded"` // 是否有被折叠的二级评论
		IsFolded  bool   `json:"is_folded"`  // 评论是否被折叠
		Rule      string `json:"rule"`       // 相关规则页面url
	} `json:"folder"` // 折叠信息
	UpAction struct {
		Like  bool `json:"like"`  // 是否UP主觉得很赞
		Reply bool `json:"reply"` // 是否被UP主回复
	} `json:"up_action"` // 评论UP主操作信息
	ShowFollow   bool `json:"show_follow"`
	Invisible    bool `json:"invisible"`
	ReplyControl struct {
	} `json:"reply_control"` // 未知
}
type CommentMember struct {
	MID         string `json:"mid"`    // 发送者mid
	Uname       string `json:"uname"`  // 发送者昵称
	Sex         string `json:"sex"`    // 发送者性别	男 女 保密
	Sign        string `json:"sign"`   // 发送者签名
	Avatar      string `json:"avatar"` // 发送者头像
	Rank        string `json:"rank"`   // 10000
	DisplayRank string `json:"DisplayRank"`
	LevelInfo   struct {
		CurrentLevel int `json:"current_level"` // 用户等级
		CurrentMin   int `json:"current_min"`   // 0
		CurrentExp   int `json:"current_exp"`   // 0
		NextExp      int `json:"next_exp"`      // 0
	} `json:"level_info"` // 发送者等级
	Pendant struct {
		PID               int64  `json:"pid"`           // 头像框id
		Name              string `json:"name"`          // 头像框名称
		Image             string `json:"image"`         // 头像框图片url
		Expire            int    `json:"expire"`        // 0
		ImageEnhance      string `json:"image_enhance"` // 头像框图片url
		ImageEnhanceFrame string `json:"image_enhance_frame"`
	} `json:"pendant"` // 发送者头像框信息
	Nameplate struct {
		NID        int    `json:"nid"`         // 勋章id
		Name       string `json:"name"`        // 勋章名称
		Image      string `json:"image"`       // 挂件图片url
		ImageSmall string `json:"image_small"` // 勋章图片url 小
		Level      string `json:"level"`       // 勋章等级
		Condition  string `json:"condition"`   // 勋章条件
	} `json:"nameplate"` // 发送者勋章信息
	OfficialVerify struct {
		// 是否认证
		//
		// -1：无
		//
		// 0：个人认证
		//
		// 1：机构认证
		Type int    `json:"type"`
		Desc string `json:"desc"` // 认证信息
	} `json:"official_verify"` // 发送者认证信息
	Vip struct {
		// 大会员类型
		//
		// 0：无
		//
		// 1：月会员
		//
		// 2：年以上会员
		VipType       int    `json:"vipType"`
		VipDueDate    int64  `json:"vipDueDate"` // 大会员到期时间	毫秒 时间戳
		DueRemark     string `json:"dueRemark"`
		AccessStatus  int    `json:"accessStatus"`
		VipStatus     int    `json:"vipStatus"` // 大会员状态	0：无 1：有
		VipStatusWarn string `json:"vipStatusWarn"`
		ThemeType     int    `json:"themeType"` // 会员样式id
		Label         struct {
			Path string `json:"path"`
			Text string `json:"text"` // 会员类型文案
			// 会员类型
			//
			// vip：大会员
			//
			// annual_vip：年度大会员
			//
			// ten_annual_vip：十年大会员
			//
			// hundred_annual_vip：百年大会员
			LabelTheme  string `json:"label_theme"`
			TextColor   string `json:"text_color"`
			BgStyle     int    `json:"bg_style"`
			BgColor     string `json:"bg_color"`
			BorderColor string `json:"border_color"`
		} `json:"label"` // 会员铭牌样式
		AvatarSubscript int    `json:"avatar_subscript"`
		NicknameColor   string `json:"nickname_color"`
	} `json:"vip"` // 发送者会员信息
	FansDetail struct {
		UID          int64  `json:"uid"`           // 用户mid
		MedalID      int64  `json:"medal_id"`      // 粉丝标签id
		MedalName    string `json:"medal_name"`    // 粉丝标签名
		Score        int    `json:"score"`         // 0
		Level        int    `json:"level"`         // 当前标签等级
		Intimacy     int    `json:"intimacy"`      // 0
		MasterStatus int    `json:"master_status"` // 1
		IsReceive    int    `json:"is_receive"`    // 1
	} `json:"fans_detail"` // 发送者粉丝标签
	// 是否关注该用户	需要登录(Cookie或APP) 否则恒为0
	//
	// 0：未关注
	//
	// 1：已关注
	Following int `json:"following"`
	// 是否被该用户关注	需要登录(Cookie或APP) 否则恒为0
	//
	// 0：未关注
	//
	// 1：已关注
	IsFollowed  int `json:"is_followed"`
	UserSailing struct {
		Pendant struct {
			ID      int64  `json:"id"`
			Name    string `json:"name"`
			Image   string `json:"image"`
			JumpURL string `json:"jump_url"`
			Type    string `json:"type"`
		} `json:"pendant"` // 头像框信息
		Cardbg struct {
			ID      int64  `json:"id"`       // 评论条目装扮id
			Name    string `json:"name"`     // 评论条目装扮名称
			Image   string `json:"image"`    // 评论条目装扮图片url
			JumpURL string `json:"jump_url"` // 评论条目装扮商城页面url
			Fan     struct {
				IsFan   int    `json:"is_fan"`   // 是否为粉丝专属装扮	0：否 1：是
				Number  int64  `json:"number"`   // 粉丝专属编号
				Color   string `json:"color"`    // 数字颜色	颜色码
				Name    string `json:"name"`     // 装扮名称
				NumDesc string `json:"num_desc"` // 粉丝专属编号	字串格式
			} `json:"fan"` // 粉丝专属信息
			// 装扮类型
			//
			// suit：一般装扮
			//
			// vip_suit：vip装扮
			Type string `json:"type"`
		} `json:"cardbg"` // 评论卡片装扮
		CardbgWithFocus interface{} `json:"cardbg_with_focus"` // 作用尚不明确
	} `json:"user_sailing"` // 发送者评论条目装扮信息
	IsContractor bool   `json:"is_contractor"` // 是否为合作用户？
	ContractDesc string `json:"contract_desc"` // 合作者信息
}
type CommentSend struct {
	SuccessAction int    `json:"success_action"`  // 0
	SuccessToast  string `json:"success_toast"`   // 状态文字
	NeedCaptcha   bool   `json:"need_captcha"`    // 评论需要验证码
	NeedCaptchaV2 bool   `json:"need_captcha_v2"` // 评论需要验证码v2
	URL           string `json:"url"`
	URLV2         string `json:"url_v2"`
	RPID          int64  `json:"rpid"`     // 评论rpid
	RpidStr       string `json:"rpid_str"` // 评论rpid4
	// 回复对方rpid
	//
	// 若为一级评论则为0
	//
	// 若为二级评论则为该评论id
	//
	// 大于二级评论为上一级评论id
	Dialog    int64  `json:"dialog"`
	DialogStr string `json:"dialog_str"`
	// 根评论rpid
	//
	// 若为一级评论则为0
	//
	// 大于一级评论则为根评论id
	Root    int64  `json:"root"`
	RootStr string `json:"root_str"`
	// 回复父评论rpid
	//
	// 若为一级评论则为0
	//
	// 若为二级评论则为根评论id
	//
	// 大于二级评论为上一级评论id
	Parent    int64    `json:"parent"`
	ParentStr string   `json:"parent_str"`
	Reply     *Comment `json:"reply"`
}
type CommentMain struct {
	Cursor struct {
		AllCount    int    `json:"all_count"`    // 全部评论条数
		IsBegin     bool   `json:"is_begin"`     // 是否为第一页
		Prev        int    `json:"prev"`         // 上页页码
		Next        int    `json:"next"`         // 下页页码
		IsEnd       bool   `json:"is_end"`       // 是否为最后页
		Mode        int    `json:"mode"`         // 排序方式
		ShowType    int    `json:"show_type"`    // 1
		SupportMode []int  `json:"support_mode"` // 支持的排序方式
		Name        string `json:"name"`         // 评论区类型名
	} `json:"cursor"` // 游标信息
	Hots   []*Comment // 热评列表
	Notice struct {
		Content string `json:"content"` // 公告正文
		ID      int64  `json:"id"`      // 公告id
		Link    string `json:"link"`    // 公告页面链接url
		Title   string `json:"title"`   // 公告标题
	} `json:"notice"` // 评论区公告信息
	Replies []*Comment // 评论列表
	Top     struct {
		Admin *Comment `json:"admin"`
		Upper *Comment `json:"upper"`
		Vote  *Comment `json:"vote"`
	} `json:"top"` // 置顶评论
	Folder struct {
		HasFolded bool   `json:"has_folded"`
		IsFolded  bool   `json:"is_folded"`
		Rule      string `json:"rule"`
	} `json:"folder"` // 评论折叠信息
	Assist    int `json:"assist"`    // 0
	Blacklist int `json:"blacklist"` // 0
	Vote      int `json:"vote"`      // 0
	Lottery   int `json:"lottery"`   // 0
	Config    struct {
		ShowAdmin  int  `json:"showadmin"`
		ShowEntry  int  `json:"showentry"`
		ShowFloor  int  `json:"showfloor"`
		ShowTopic  int  `json:"showtopic"`
		ShowUpFlag bool `json:"show_up_flag"`
		ReadOnly   bool `json:"read_only"`
		ShowDelLog bool `json:"show_del_log"`
	} `json:"config"` // 评论区显示控制
	Upper struct {
		MID int64 `json:"mid"`
	} `json:"upper"` // UP主信息
	ShowBvid bool `json:"show_bvid"`
	Control  struct {
		InputDisable          bool   `json:"input_disable"`            // 禁止评论?
		RootInputText         string `json:"root_input_text"`          // 评论框文字
		ChildInputText        string `json:"child_input_text"`         // 评论框文字
		GiveUpInputText       string `json:"giveup_input_text"`        // 放弃评论后的评论框文字
		BgText                string `json:"bg_text"`                  // 空评论区文字
		WebSelection          bool   `json:"web_selection"`            // 评论是否筛选后可见 false：无需筛选 true：需要筛选
		AnswerGuideText       string `json:"answer_guide_text"`        // 答题页面链接文字
		AnswerGuideIconURL    string `json:"answer_guide_icon_url"`    // 答题页面图标url
		AnswerGuideIosURL     string `json:"answer_guide_ios_url"`     // 答题页面ios url
		AnswerGuideAndroidURL string `json:"answer_guide_android_url"` // 答题页面安卓url
		ShowType              int    `json:"show_type"`
		ShowText              string `json:"show_text"`
	} `json:"control"`
}

type CommentReply struct {
	Config struct {
		ShowAdmin  int  `json:"showadmin"`
		ShowEntry  int  `json:"showentry"`
		ShowFloor  int  `json:"showfloor"`
		Showtopic  int  `json:"showtopic"`
		ShowUpFlag bool `json:"show_up_flag"`
		ReadOnly   bool `json:"read_only"`
		ShowDelLog bool `json:"show_del_log"`
	} `json:"config"`
	Control struct {
		InputDisable          bool   `json:"input_disable"`
		RootInputText         string `json:"root_input_text"`
		ChildInputText        string `json:"child_input_text"`
		GiveUpInputText       string `json:"giveup_input_text"`
		BgText                string `json:"bg_text"`
		WebSelection          bool   `json:"web_selection"`
		AnswerGuideText       string `json:"answer_guide_text"`
		AnswerGuideIconURL    string `json:"answer_guide_icon_url"`
		AnswerGuideIosURL     string `json:"answer_guide_ios_url"`
		AnswerGuideAndroidURL string `json:"answer_guide_android_url"`
		ShowType              int    `json:"show_type"`
		ShowText              string `json:"show_text"`
	} `json:"control"`
	Page struct {
		Count int `json:"count"`
		Num   int `json:"num"`
		Size  int `json:"size"`
	} `json:"page"`
	Root     Comment    `json:"root"`
	Replies  []*Comment `json:"replies"`
	ShowBvid bool       `json:"show_bvid"`
	ShowText string     `json:"show_text"`
	ShowType int        `json:"show_type"`
	Upper    struct {
		MID int64 `json:"mid"`
	} `json:"upper"`
}
type UserInfo struct {
	MID      int64  `json:"mid"`      // mid
	Name     string `json:"name"`     // 昵称
	Sex      string `json:"sex"`      // 性别 男/女/保密
	Face     string `json:"face"`     // 头像链接
	Sign     string `json:"sign"`     // 签名
	Rank     int    `json:"rank"`     // 10000
	Level    int    `json:"level"`    // 当前等级	0-6级
	JoinTime int64  `json:"jointime"` // 0
	Moral    int    `json:"moral"`    // 0
	Silence  int    `json:"silence"`  // 封禁状态 0：正常 1：被封
	// 硬币数 需要登录(Cookie)
	//
	// 只能查看自己的
	//
	// 默认为0
	Coins     int  `json:"coins"`
	FansBadge bool `json:"fans_badge"` // 是否具有粉丝勋章 false：无 true：有
	Official  struct {
		// 认证类型
		//
		// 0：无
		//
		// 1 2 7：个人认证
		//
		// 3 4 5 6：机构认证
		Role  int    `json:"role"`
		Title string `json:"title"` // 认证信息
		Desc  string `json:"desc"`  // 认证备注
		Type  int    `json:"type"`  // 是否认证 -1：无 0：认证
	} `json:"official"`
	Vip struct {
		// 会员类型
		//
		// 0：无
		//
		// 1：月大会员
		//
		// 2：年度及以上大会员
		Type       int   `json:"type"`
		Status     int   `json:"status"`       // 会员状态 0：无 1：有
		DueDate    int64 `json:"due_date"`     // 会员过期时间 Unix时间戳(毫秒)
		VipPayType int   `json:"vip_pay_type"` // 支付类型?
		ThemeType  int   `json:"theme_type"`   // 0
		Label      struct {
			Path string `json:"path"` //
			Text string `json:"text"` // 会员类型文案
			// 会员标签
			//
			// vip：大会员
			//
			// annual_vip：年度大会员
			//
			// ten_annual_vip：十年大会员
			//
			// hundred_annual_vip：百年大会员
			LabelTheme  string `json:"label_theme"`
			TextColor   string `json:"text_color"`   // 文字颜色
			BgStyle     int    `json:"bg_style"`     // 背景类型
			BgColor     string `json:"bg_color"`     // 背景颜色
			BorderColor string `json:"border_color"` // 边角颜色
		} `json:"label"` //
		AvatarSubscript    int    `json:"avatar_subscript"`     // 是否显示会员图标	0：不显示 1：显示
		NicknameColor      string `json:"nickname_color"`       // 会员昵称颜色 颜色码
		Role               int    `json:"role"`                 //
		AvatarSubscriptUrl string `json:"avatar_subscript_url"` // 会员图标url
	} `json:"vip"` //
	Pendant struct {
		PID               int64  `json:"pid"`                 // 头像框id
		Name              string `json:"name"`                // 头像框名称
		Image             string `json:"image"`               // 头像框图片url
		Expire            int64  `json:"expire"`              // 0
		ImageEnhance      string `json:"image_enhance"`       //
		ImageEnhanceFrame string `json:"image_enhance_frame"` //
	} `json:"pendant"` //
	Nameplate struct {
		NID        int    `json:"nid"`         // 勋章id
		Name       string `json:"name"`        // 勋章名称
		Image      string `json:"image"`       // 挂件图片url 正常
		ImageSmall string `json:"image_small"` // 勋章图片url 小
		Level      string `json:"level"`       // 勋章等级
		Condition  string `json:"condition"`   // 勋章条件
	} `json:"nameplate"` //
	// 是否关注此用户	true：已关注 false：未关注
	//
	// 需要登录(Cookie)
	//
	// 未登录恒为false
	IsFollowed bool   `json:"is_followed"` //
	TopPhoto   string `json:"top_photo"`   // 主页头图链接
	SysNotice  struct {
		// 系统提示类型id
		//
		// 8：争议账号
		//
		// 20：纪念账号
		//
		// 22：合约诉讼
		ID         int    `json:"id"`
		Content    string `json:"content"`     // 提示文案
		URL        string `json:"url"`         // 提示信息页面url
		NoticeType int    `json:"notice_type"` //
		Icon       string `json:"icon"`        // 提示图标url
		TextColor  string `json:"text_color"`  // 提示文字颜色
		BgColor    string `json:"bg_color"`    // 提示背景颜色
	} `json:"sys_notice"` //
	LiveRoom struct {
		RoomStatus    int    `json:"roomStatus"`     // 直播间状态 0：无房间 1：有房间
		RoundStatus   int    `json:"roundStatus"`    // 轮播状态 0：未轮播 1：轮播
		LiveStatus    int    `json:"liveStatus"`     // 直播状态 0：未开播 1：直播中
		URL           string `json:"url"`            // 直播间网页url
		Title         string `json:"title"`          // 直播间标题
		Cover         string `json:"cover"`          // 直播间封面url
		Online        int    `json:"online"`         // 直播间人气 值为上次直播时刷新
		RoomID        int    `json:"roomid"`         // 直播间id(真实ID)
		BroadcastType int    `json:"broadcast_type"` // 0
		OnlineHidden  int    `json:"online_hidden"`  // 已废弃
	} `json:"live_room"` //
	Birthday string `json:"birthday"` // 生日 MM-DD 如设置隐私为空
	School   struct {
		Name string `json:"name"` // 就读大学名称
	} `json:"school"` //
	Profession struct {
		Name string `json:"name"` //
	} `json:"profession"` //
	Series struct {
		UserUpgradeStatus int  `json:"user_upgrade_status"` //
		ShowUpgradeWindow bool `json:"show_upgrade_window"` //
	} `json:"series"` //
}
