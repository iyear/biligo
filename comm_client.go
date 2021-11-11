package biligo

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/iyear/biligo/internal/util"
	"github.com/iyear/biligo/proto/dm"
	"github.com/tidwall/gjson"
	"net/http"
	"strconv"
)

type CommClient struct {
	*baseClient
}
type CommSetting struct {

	// 自定义http client
	//
	// 默认为 http.http.DefaultClient
	Client *http.Client

	// Debug模式 true将输出请求信息 false不输出
	//
	// 默认false
	DebugMode bool

	// 自定义UserAgent
	//
	// 默认Chrome随机Agent
	UserAgent string
}

// NewCommClient
//
// Setting的Auth属性可以随意填写或传入nil，Auth不起到作用，用于访问公共API
func NewCommClient(setting *CommSetting) *CommClient {
	return &CommClient{baseClient: newBaseClient(&baseSetting{
		Client:    setting.Client,
		DebugMode: setting.DebugMode,
		UserAgent: setting.UserAgent,
		Prefix:    "CommClient ",
	})}
}

// SetClient
//
// 设置Client,可以用来更换代理等操作
func (c *CommClient) SetClient(client *http.Client) {
	c.client = client
}

// SetUA
//
// 设置UA
func (c *CommClient) SetUA(ua string) {
	c.ua = ua
}

// Raw
//
// base末尾带/
func (c *CommClient) Raw(base, endpoint, method string, payload map[string]string) ([]byte, error) {
	// 不用侵入处理则传入nil
	raw, err := c.raw(base, endpoint, method, payload, nil, nil)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// RawParse
//
// base末尾带/
func (c *CommClient) RawParse(base, endpoint, method string, payload map[string]string) (*Response, error) {
	raw, err := c.Raw(base, endpoint, method, payload)
	if err != nil {
		return nil, err
	}
	return c.parse(raw)
}

// GetGeoInfo 调用哔哩哔哩API获取地理位置等信息
//
// 会受到自定义 http.Client 代理的影响
func (c *CommClient) GetGeoInfo() (*GeoInfo, error) {
	resp, err := c.RawParse(BiliApiURL,
		"x/web-interface/zone",
		"GET",
		nil,
	)
	if err != nil {
		return nil, err
	}
	var info *GeoInfo
	if err = json.Unmarshal(resp.Data, &info); err != nil {
		return nil, err
	}
	return info, nil
}

// FollowingsGetDetail 获取个人详细的关注列表
//
// pn 页码
//
// ps 每页项数，最大50
func (c *CommClient) FollowingsGetDetail(mid int64, pn int, ps int) (*FollowingsDetail, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/relation/followings",
		"GET",
		map[string]string{
			"vmid": strconv.FormatInt(mid, 10),
			"pn":   strconv.Itoa(pn),
			"ps":   strconv.Itoa(ps),
		},
	)
	if err != nil {
		return nil, err
	}

	var detail = &FollowingsDetail{}
	if err = json.Unmarshal(resp.Data, &detail); err != nil {
		return nil, err
	}
	return detail, nil
}

// VideoGetStat
//
// 获取稿件状态数
func (c *CommClient) VideoGetStat(aid int64) (*VideoSingleStat, error) {
	resp, err := c.RawParse(BiliApiURL,
		"x/web-interface/archive/stat",
		"GET",
		map[string]string{
			"aid": strconv.FormatInt(aid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var stat *VideoSingleStat
	if err = json.Unmarshal(resp.Data, &stat); err != nil {
		return nil, err
	}
	return stat, nil
}

// VideoGetInfo 返回视频详细信息，数据较多，可以使用单独的接口获取部分数据
//
// 限制游客访问的视频会返回错误，请使用 BiliClient 发起请求
func (c *CommClient) VideoGetInfo(aid int64) (*VideoInfo, error) {
	resp, err := c.RawParse(BiliApiURL,
		"x/web-interface/view",
		"GET",
		map[string]string{
			"aid": strconv.FormatInt(aid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var info *VideoInfo
	if err = json.Unmarshal(resp.Data, &info); err != nil {
		return nil, err
	}
	return info, nil
}

// VideoGetDescription
//
// 获取稿件简介
func (c *CommClient) VideoGetDescription(aid int64) (string, error) {
	resp, err := c.RawParse(BiliApiURL,
		"x/web-interface/archive/desc",
		"GET",
		map[string]string{
			"aid": strconv.FormatInt(aid, 10),
		},
	)
	if err != nil {
		return "", err
	}
	var desc string
	if err = json.Unmarshal(resp.Data, &desc); err != nil {
		return "", err
	}
	return desc, nil
}

// VideoGetPageList
//
// 获取分P列表
func (c *CommClient) VideoGetPageList(aid int64) ([]*VideoPage, error) {
	resp, err := c.RawParse(BiliApiURL,
		"x/player/pagelist",
		"GET",
		map[string]string{
			"aid": strconv.FormatInt(aid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var list []*VideoPage
	if err = json.Unmarshal(resp.Data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// VideoGetOnlineNum
//
// 返回所有终端总计在线观看人数和WEB端在线观看人数 (用类似10万+的文字表示) cid用于分P标识
func (c *CommClient) VideoGetOnlineNum(aid int64, cid int64) (total string, web string, e error) {
	resp, err := c.RawParse(BiliApiURL,
		"x/player/online/total",
		"GET",
		map[string]string{
			"aid": strconv.FormatInt(aid, 10),
			"cid": strconv.FormatInt(cid, 10),
		},
	)
	if err != nil {
		return "", "", err
	}
	var num struct {
		Total string `json:"total,omitempty"`
		Count string `json:"count,omitempty"`
	}
	if err = json.Unmarshal(resp.Data, &num); err != nil {
		return "", "", err
	}
	return num.Total, num.Count, nil
}

// VideoTags
//
// 未登录无法获取 IsAtten,Liked,Hated 字段
func (c *CommClient) VideoTags(aid int64) ([]*VideoTag, error) {
	resp, err := c.RawParse(BiliApiURL,
		"x/tag/archive/tags",
		"GET",
		map[string]string{
			"aid": strconv.FormatInt(aid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var tags []*VideoTag
	if err = json.Unmarshal(resp.Data, &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

// VideoGetRecommend 获取视频的相关视频推荐
//
// 最多获取40条推荐视频
func (c *CommClient) VideoGetRecommend(aid int64) ([]*VideoRecommendInfo, error) {
	resp, err := c.RawParse(BiliApiURL,
		"x/web-interface/archive/related",
		"GET",
		map[string]string{
			"aid": strconv.FormatInt(aid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var videos []*VideoRecommendInfo
	if err = json.Unmarshal(resp.Data, &videos); err != nil {
		return nil, err
	}
	return videos, nil
}

// VideoGetPlayURL 获取视频取流地址
//
// 所有参数、返回信息和取流方法的说明请直接前往：https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/video/videostream_url.md
func (c *CommClient) VideoGetPlayURL(aid int64, cid int64, qn int, fnval int) (*VideoPlayURLResult, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/player/playurl",
		"GET",
		map[string]string{
			"avid":  strconv.FormatInt(aid, 10),
			"cid":   strconv.FormatInt(cid, 10),
			"qn":    strconv.Itoa(qn),
			"fnval": strconv.Itoa(fnval),
			"fnver": "0",
			"fourk": "1",
		},
	)
	if err != nil {
		return nil, err
	}
	var r *VideoPlayURLResult
	if err = json.Unmarshal(resp.Data, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// VideoShot 获取视频快照
//
// cid属性非必须 传入0表示1P
//
// index为JSON数组截取时间表 true:需要 false:不需要 传入false则Index属性为空
func (c *CommClient) VideoShot(aid int64, cid int64, index bool) (*VideoShot, error) {
	resp, err := c.RawParse(BiliApiURL,
		"x/player/videoshot",
		"GET",
		map[string]string{
			"aid":   strconv.FormatInt(aid, 10),
			"cid":   util.IF(cid == 0, "", strconv.FormatInt(cid, 10)).(string),
			"index": util.IF(index, "1", "0").(string),
		},
	)
	if err != nil {
		return nil, err
	}
	var shot *VideoShot
	if err = json.Unmarshal(resp.Data, &shot); err != nil {
		return nil, err
	}
	return shot, nil
}

// DanmakuGetLikes 获取弹幕点赞数，一次可以获取多条弹幕
//
// Link:https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/danmaku/action.md#%E6%9F%A5%E8%AF%A2%E5%BC%B9%E5%B9%95%E7%82%B9%E8%B5%9E%E6%95%B0
//
// 返回一个map，key为dmid，value为相关信息
// 未登录时UserLike属性恒为0
func (c *CommClient) DanmakuGetLikes(cid int64, dmids []uint64) (map[uint64]*DanmakuGetLikesResult, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/v2/dm/thumbup/stats",
		"GET",
		map[string]string{
			"oid": strconv.FormatInt(cid, 10),
			"ids": util.Uint64SliceToString(dmids, ","),
		},
	)
	if err != nil {
		return nil, err
	}
	var result = make(map[uint64]*DanmakuGetLikesResult)
	for _, dmid := range dmids {
		var r *DanmakuGetLikesResult
		if err = json.Unmarshal([]byte(gjson.Get(string(resp.Data), strconv.FormatUint(dmid, 10)).Raw), &r); err != nil {
			return nil, err
		}
		result[dmid] = r
	}
	return result, nil
}

// GetRelationStat
//
// 获取关系状态数，Whisper和Black恒为0
func (c *CommClient) GetRelationStat(mid int64) (*RelationStat, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/relation/stat",
		"GET",
		map[string]string{
			"vmid": strconv.FormatInt(mid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var stat *RelationStat
	if err = json.Unmarshal(resp.Data, &stat); err != nil {
		return nil, err
	}
	return stat, nil
}

// SpaceGetTopArchive
//
// 获取空间置顶稿件
func (c *CommClient) SpaceGetTopArchive(mid int64) (*SpaceVideo, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/space/top/arc",
		"GET",
		map[string]string{
			"vmid": strconv.FormatInt(mid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var top *SpaceVideo
	if err = json.Unmarshal(resp.Data, &top); err != nil {
		return nil, err
	}
	return top, nil
}

// SpaceGetMasterpieces
//
// 获取UP代表作
func (c *CommClient) SpaceGetMasterpieces(mid int64) ([]*SpaceVideo, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/space/masterpiece",
		"GET",
		map[string]string{
			"vmid": strconv.FormatInt(mid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var mp []*SpaceVideo
	if err = json.Unmarshal(resp.Data, &mp); err != nil {
		return nil, err
	}
	return mp, nil
}

// SpaceGetTags
//
// 获取空间用户个人TAG 上限5条，且内容由用户自定义 带有转义
func (c *CommClient) SpaceGetTags(mid int64) ([]string, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/space/acc/tags",
		"GET",
		map[string]string{
			"mid": strconv.FormatInt(mid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	// B站这写的是个啥玩意儿
	var tags []string
	for _, tag := range gjson.Get(string(resp.Data), "0.tags").Array() {
		tags = append(tags, tag.String())
	}
	return tags, nil
}

// SpaceGetNotice
//
// 获取空间公告内容
func (c *CommClient) SpaceGetNotice(mid int64) (string, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/space/notice",
		"GET",
		map[string]string{
			"mid": strconv.FormatInt(mid, 10),
		},
	)
	if err != nil {
		return "", err
	}
	// 新建一个变量再 unmarshal 可以把转义部分转回来
	// 直接返回 resp.Data 会带转义符
	var notice string
	if err = json.Unmarshal(resp.Data, &notice); err != nil {
		return "", err
	}
	return notice, nil
}

// SpaceGetLastPlayGame
//
// 获取用户空间近期玩的游戏
func (c *CommClient) SpaceGetLastPlayGame(mid int64) ([]*SpaceGame, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/space/lastplaygame",
		"GET",
		map[string]string{
			"mid": strconv.FormatInt(mid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var games []*SpaceGame
	if err = json.Unmarshal(resp.Data, &games); err != nil {
		return nil, err
	}
	return games, nil
}

// SpaceGetLastVideoCoin
//
// 获取用户最近投币的视频明细 如设置隐私查看自己的使用 BiliClient 访问
func (c *CommClient) SpaceGetLastVideoCoin(mid int64) ([]*SpaceVideoCoin, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/space/coin/video",
		"GET",
		map[string]string{
			"vmid": strconv.FormatInt(mid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var info []*SpaceVideoCoin
	if err = json.Unmarshal(resp.Data, &info); err != nil {
		return nil, err
	}
	return info, nil
}

// SpaceSearchVideo
//
// 获取用户投稿视频明细
//
// order 排序方式 默认为pubdate 可留空
//
// 最新发布:pubdate
// 最多播放:click
// 最多收藏:stow
//
// tid 筛选分区 0:不进行分区筛选
//
// keyword 关键词 可留空
//
// pn 页码
//
// ps 每页项数
func (c *CommClient) SpaceSearchVideo(mid int64, order string, tid int, keyword string, pn int, ps int) (*SpaceVideoSearchResult, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/space/arc/search",
		"GET",
		map[string]string{
			"mid":     strconv.FormatInt(mid, 10),
			"order":   order,
			"tid":     strconv.Itoa(tid),
			"keyword": keyword,
			"pn":      strconv.Itoa(pn),
			"ps":      strconv.Itoa(ps),
		},
	)
	if err != nil {
		return nil, err
	}
	var result *SpaceVideoSearchResult
	if err = json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ChanGet
//
// 获取用户频道列表
func (c *CommClient) ChanGet(mid int64) (*ChannelList, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/space/channel/list",
		"GET",
		map[string]string{
			"mid": strconv.FormatInt(mid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var list *ChannelList
	if err = json.Unmarshal(resp.Data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// ChanGetVideo
//
// 获取用户频道视频
//
// cid 频道ID
//
// pn 页码
//
// ps 每页项数
func (c *CommClient) ChanGetVideo(mid int64, cid int64, pn int, ps int) (*ChanVideo, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/space/channel/video",
		"GET",
		map[string]string{
			"mid": strconv.FormatInt(mid, 10),
			"cid": strconv.FormatInt(cid, 10),
			"pn":  strconv.Itoa(pn),
			"ps":  strconv.Itoa(ps),
		},
	)
	if err != nil {
		return nil, err
	}
	var videos *ChanVideo
	if err = json.Unmarshal(resp.Data, &videos); err != nil {
		return nil, err
	}
	return videos, nil
}

// FavGet
//
// 获取用户的公开收藏夹列表
func (c *CommClient) FavGet(mid int64) (*FavoritesList, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/v3/fav/folder/created/list-all",
		"GET",
		map[string]string{
			"up_mid": strconv.FormatInt(mid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var list *FavoritesList
	if err = json.Unmarshal(resp.Data, &list); err != nil {
		return nil, err
	}
	if list == nil {
		return &FavoritesList{}, nil
	}
	return list, nil
}

// FavGetDetail
//
// 获取收藏夹详细信息，部分信息需要登录，请使用 BiliClient 请求
func (c *CommClient) FavGetDetail(mlid int64) (*FavDetail, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/v3/fav/folder/info",
		"GET",
		map[string]string{
			"media_id": strconv.FormatInt(mlid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var detail *FavDetail
	if err = json.Unmarshal(resp.Data, &detail); err != nil {
		return nil, err
	}
	return detail, nil
}

// FavGetRes
//
// 获取收藏夹全部内容id 查询权限收藏夹时请使用 BiliClient 请求
func (c *CommClient) FavGetRes(mlid int64) ([]*FavRes, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/v3/fav/resource/ids",
		"GET",
		map[string]string{
			"media_id": strconv.FormatInt(mlid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var r []*FavRes
	if err = json.Unmarshal(resp.Data, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// FavGetResDetail 获取收藏夹内容详细内容，带过滤功能
//
// 查询权限收藏夹时请使用 BiliClient 请求
//
// tid 分区id，用于筛选，传入0代表所有分区
//
// keyword 关键词筛选 可留空
//
// order 留空默认按收藏时间
//
// 按收藏时间:mtime
// 按播放量: view
// 按投稿时间：pubtime
//
// tp 内容类型 不知道作用，传入0即可
//
// pn 页码
//
// ps 每页项数 ps不能太大，会报错
func (c *CommClient) FavGetResDetail(mlid int64, tid int, keyword string, order string, tp int, pn int, ps int) (*FavResDetail, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/v3/fav/resource/list",
		"GET",
		map[string]string{
			"media_id": strconv.FormatInt(mlid, 10),
			"tid":      strconv.Itoa(tid),
			"keyword":  keyword,
			"order":    order,
			"type":     strconv.Itoa(tp),
			"ps":       strconv.Itoa(ps),
			"pn":       strconv.Itoa(pn),
		},
	)
	if err != nil {
		return nil, err
	}
	var detail *FavResDetail
	if err = json.Unmarshal(resp.Data, &detail); err != nil {
		return nil, err
	}
	return detail, nil
}

// GetDailyNum
//
// 获取每日分区投稿数
func (c *CommClient) GetDailyNum() (map[int]int, error) {
	resp, err := c.RawParse(BiliApiURL, "x/web-interface/online", "GET", nil)
	if err != nil {
		return nil, err
	}
	var result = make(map[int]int)
	gjson.Get(string(resp.Data), "region_count").ForEach(func(key, value gjson.Result) bool {
		result[int(key.Int())] = int(value.Int())
		return true
	})
	return result, nil
}

// GetUnixNow
//
// 获取服务器的Unix时间戳
func (c *CommClient) GetUnixNow() (int64, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/report/click/now",
		"GET",
		nil,
	)
	if err != nil {
		return -1, err
	}
	var t struct {
		Now int64 `json:"now,omitempty"`
	}
	if err = json.Unmarshal(resp.Data, &t); err != nil {
		return -1, err
	}
	return t.Now, nil
}

// DanmakuGetByPb
//
// 获取实时弹幕(protobuf接口)
func (c *CommClient) DanmakuGetByPb(tp int, cid int64, seg int) (*DanmakuResp, error) {
	resp, err := c.Raw(
		BiliApiURL,
		"x/v2/dm/web/seg.so",
		"GET",
		map[string]string{
			"type":          strconv.Itoa(tp),
			"oid":           strconv.FormatInt(cid, 10),
			"segment_index": strconv.Itoa(seg),
		},
	)
	if err != nil {
		return nil, err
	}
	var reply dm.DmSegMobileReply
	var r = &DanmakuResp{}
	if err := proto.Unmarshal(resp, &reply); err != nil {
		return nil, err
	}
	for _, elem := range reply.GetElems() {
		r.Danmaku = append(r.Danmaku, &Danmaku{
			ID:       uint64(elem.Id),
			Progress: int64(elem.Progress),
			Mode:     int(elem.Mode),
			FontSize: int(elem.Fontsize),
			Color:    int(elem.Color),
			MidHash:  elem.MidHash,
			Content:  elem.Content,
			Ctime:    elem.Ctime,
			Weight:   int(elem.Weight),
			Action:   elem.Action,
			Pool:     int(elem.Pool),
			IDStr:    elem.IdStr,
			Attr:     int(elem.Attr),
		})
	}
	return r, nil

}

// DanmakuGetShot
//
// 获取弹幕快照(最新的几条弹幕)
func (c *CommClient) DanmakuGetShot(aid int64) ([]string, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/v2/dm/ajax",
		"GET",
		map[string]string{
			"aid": strconv.FormatInt(aid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var strings []string
	if err = json.Unmarshal(resp.Data, &strings); err != nil {
		return nil, err
	}
	return strings, nil
}

// EmoteGetFreePack 获取免费表情包列表
//
// business 使用场景	reply：评论区 dynamic：动态
//
// 全为免费表情包，如需获取个人专属表情包请使用 BiliClient 请求
func (c *CommClient) EmoteGetFreePack(business string) ([]*EmotePack, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/emote/user/panel/web",
		"GET",
		map[string]string{
			"business": business,
		},
	)
	if err != nil {
		return nil, err
	}
	var pack struct {
		Packages []*EmotePack `json:"packages,omitempty"`
	}
	if err = json.Unmarshal(resp.Data, &pack); err != nil {
		return nil, err
	}
	return pack.Packages, nil
}

// EmoteGetPackDetail 获取指定表情包明细
//
// business 使用场景	reply：评论区 dynamic：动态
//
// ids 多个表情包id的数组
func (c *CommClient) EmoteGetPackDetail(business string, ids []int64) ([]*EmotePack, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/emote/package",
		"GET",
		map[string]string{
			"business": business,
			"ids":      util.Int64SliceToString(ids, ","),
		},
	)
	if err != nil {
		return nil, err
	}
	var packs struct {
		Packages []*EmotePack `json:"packages,omitempty"`
	}
	if err = json.Unmarshal(resp.Data, &packs); err != nil {
		return nil, err
	}
	return packs.Packages, nil
}

// AudioGetInfo
//
// 获取音频信息 部分属性需要登录，请使用 BiliClient 请求
func (c *CommClient) AudioGetInfo(auid int64) (*AudioInfo, error) {
	resp, err := c.RawParse(
		BiliMainURL,
		"audio/music-service-c/web/song/info",
		"GET",
		map[string]string{
			"sid": strconv.FormatInt(auid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var info *AudioInfo
	if err = json.Unmarshal(resp.Data, &info); err != nil {
		return nil, err
	}
	return info, nil
}

// AudioGetTags 获取音频TAGs
//
// 根据页面显示观察，应该是歌曲分类
func (c *CommClient) AudioGetTags(auid int64) ([]*AudioTag, error) {
	resp, err := c.RawParse(
		BiliMainURL,
		"audio/music-service-c/web/tag/song",
		"GET",
		map[string]string{
			"sid": strconv.FormatInt(auid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var tags []*AudioTag
	if err = json.Unmarshal(resp.Data, &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

// AudioGetMembers
//
// 获取音频创作者信息
func (c *CommClient) AudioGetMembers(auid int64) ([]*AudioMember, error) {
	resp, err := c.RawParse(
		BiliMainURL,
		"audio/music-service-c/web/member/song",
		"GET",
		map[string]string{
			"sid": strconv.FormatInt(auid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var members []*AudioMember
	if err = json.Unmarshal(resp.Data, &members); err != nil {
		return nil, err
	}
	return members, nil
}

// AudioGetLyric 获取音频歌词
//
// 同 AudioGetInfo 中的lrc歌词
func (c *CommClient) AudioGetLyric(auid int64) (string, error) {
	resp, err := c.RawParse(
		BiliMainURL,
		"audio/music-service-c/web/song/lyric",
		"GET",
		map[string]string{
			"sid": strconv.FormatInt(auid, 10),
		},
	)
	if err != nil {
		return "", err
	}
	var lrc string
	if err = json.Unmarshal(resp.Data, &lrc); err != nil {
		return "", err
	}
	return lrc, nil
}

// AudioGetStat 获取歌曲状态数
//
// 没有投币数 获取投币数请使用 AudioGetInfo
func (c *CommClient) AudioGetStat(auid int64) (*AudioInfoStat, error) {
	resp, err := c.RawParse(
		BiliMainURL,
		"audio/music-service-c/web/stat/song",
		"GET",
		map[string]string{
			"sid": strconv.FormatInt(auid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var stat = &AudioInfoStat{}
	if err = json.Unmarshal(resp.Data, &stat); err != nil {
		return nil, err
	}
	return stat, nil
}

// AudioGetPlayURL 获取音频流URL
//
// 最多获取到
//
// qn 音质
//
// 0 流畅 128K
//
// 1 标准 192K
//
// 2 高品质 320K
//
// 3 无损 FLAC（大会员）
//
// 最高获取到 320K 音质,更高音质请使用 BiliClient 请求
//
// 取流：https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/audio/musicstream_url.md#%E9%9F%B3%E9%A2%91%E6%B5%81%E7%9A%84%E8%8E%B7%E5%8F%96
func (c *CommClient) AudioGetPlayURL(auid int64, qn int) (*AudioPlayURL, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"audio/music-service-c/url",
		"GET",
		map[string]string{
			"songid":    strconv.FormatInt(auid, 10),
			"quality":   strconv.Itoa(qn),
			"privilege": "2",
			"mid":       "2",
			"platform":  "web",
		},
	)
	if err != nil {
		return nil, err
	}
	var play *AudioPlayURL
	if err = json.Unmarshal(resp.Data, &play); err != nil {
		return nil, err
	}
	return play, nil
}

// ChargeSpaceGetList
//
// 获取用户空间充电名单
func (c *CommClient) ChargeSpaceGetList(mid int64) (*ChargeSpaceList, error) {
	resp, err := c.RawParse(
		BiliElecURL,
		"api/query.rank.do",
		"GET",
		map[string]string{
			"mid": strconv.FormatInt(mid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var list *ChargeSpaceList
	if err = json.Unmarshal(resp.Data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// ChargeVideoGetList
//
// 获取用户视频充电名单
func (c *CommClient) ChargeVideoGetList(mid int64, aid int64) (*ChargeVideoList, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/web-interface/elec/show",
		"GET",
		map[string]string{
			"mid": strconv.FormatInt(mid, 10),
			"aid": strconv.FormatInt(aid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var list *ChargeVideoList
	if err = json.Unmarshal(resp.Data, &list); err != nil {
		return nil, err
	}
	return list, nil

}

// LiveGetRoomInfoByMID
//
// 从mid获取直播间信息
func (c *CommClient) LiveGetRoomInfoByMID(mid int64) (*LiveRoomInfoByMID, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/space/acc/info",
		"GET",
		map[string]string{
			"mid": strconv.FormatInt(mid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var info struct {
		LiveRoom *LiveRoomInfoByMID `json:"live_room"`
	}
	if err = json.Unmarshal(resp.Data, &info); err != nil {
		return nil, err
	}
	return info.LiveRoom, nil
}

// LiveGetRoomInfoByID 从roomID获取直播间信息
//
// roomID 可为短号也可以是真实房号
func (c *CommClient) LiveGetRoomInfoByID(roomID int64) (*LiveRoomInfoByID, error) {
	resp, err := c.RawParse(
		BiliLiveURL,
		"xlive/web-room/v1/index/getRoomPlayInfo",
		"GET",
		map[string]string{
			"room_id": strconv.FormatInt(roomID, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var r = &LiveRoomInfoByID{}
	if err = json.Unmarshal(resp.Data, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// LiveGetWsConf 获取直播websocket服务器信息
//
// roomID: 真实直播间ID
func (c *CommClient) LiveGetWsConf(roomID int64) (*LiveWsConf, error) {
	resp, err := c.RawParse(
		BiliLiveURL,
		"room/v1/Danmu/getConf",
		"GET",
		map[string]string{
			"room_id": strconv.FormatInt(roomID, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var r = &LiveWsConf{}
	if err = json.Unmarshal(resp.Data, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// LiveGetAreaInfo
//
// 获取直播分区信息
func (c *CommClient) LiveGetAreaInfo() ([]*LiveAreaInfo, error) {
	resp, err := c.RawParse(
		BiliLiveURL,
		"room/v1/Area/getList",
		"GET",
		map[string]string{},
	)
	if err != nil {
		return nil, err
	}
	var r []*LiveAreaInfo
	if err = json.Unmarshal(resp.Data, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// LiveGetGuardList 获取直播间大航海列表
//
// roomID: 真实直播间ID
//
// mid: 主播mid
//
// pn: 页码
//
// ps: 每页项数
func (c *CommClient) LiveGetGuardList(roomID int64, mid int64, pn int, ps int) (*LiveGuardList, error) {
	resp, err := c.RawParse(
		BiliLiveURL,
		"xlive/app-room/v1/guardTab/topList",
		"GET",
		map[string]string{
			"roomid":    strconv.FormatInt(roomID, 10),
			"ruid":      strconv.FormatInt(mid, 10),
			"page":      strconv.Itoa(pn),
			"page_size": strconv.Itoa(ps),
		},
	)
	if err != nil {
		return nil, err
	}
	var r = &LiveGuardList{}
	if err = json.Unmarshal(resp.Data, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// LiveGetMedalRank 获取直播间粉丝勋章榜
//
// roomID: 真实直播间ID
//
// mid: 主播mid
func (c *CommClient) LiveGetMedalRank(roomID int64, mid int64) (*LiveMedalRank, error) {
	resp, err := c.RawParse(
		BiliLiveURL,
		"rankdb/v1/RoomRank/webMedalRank",
		"GET",
		map[string]string{
			"roomid": strconv.FormatInt(roomID, 10),
			"ruid":   strconv.FormatInt(mid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var r = &LiveMedalRank{}
	if err = json.Unmarshal(resp.Data, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// LiveGetPlayURL 获取直播流信息
//
// qn: 原画:10000 蓝光:400 超清:250 高清:150 流畅:80
func (c *CommClient) LiveGetPlayURL(roomID int64, qn int) (*LivePlayURL, error) {
	resp, err := c.RawParse(
		BiliLiveURL,
		"xlive/web-room/v1/playUrl/playUrl",
		"GET",
		map[string]string{
			"cid":           strconv.FormatInt(roomID, 10),
			"qn":            strconv.Itoa(qn),
			"platform":      "web",
			"https_url_req": "1",
			"ptype":         "16",
		},
	)
	if err != nil {
		return nil, err
	}
	var r = &LivePlayURL{}
	if err = json.Unmarshal(resp.Data, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// LiveGetAllGiftInfo 获取所有礼物信息
//
// areaID: 子分区ID 从 LiveGetAreaInfo 获取
//
// areaParentID: 父分区ID 从 LiveGetAreaInfo 获取
//
// 三个字段可以不用填，但填了有助于减小返回内容的大小，置空(传入0)返回约 2.7w 行，填了三个对应值返回约 1.4w 行
func (c *CommClient) LiveGetAllGiftInfo(roomID int64, areaID int, areaParentID int) (*LiveAllGiftInfo, error) {
	resp, err := c.RawParse(
		BiliLiveURL,
		"xlive/web-room/v1/giftPanel/giftConfig",
		"GET",
		map[string]string{
			"room_id":        strconv.FormatInt(roomID, 10),
			"platform":       "pc",
			"source":         "live",
			"area_id":        strconv.Itoa(areaID),
			"area_parent_id": strconv.Itoa(areaParentID),
		},
	)
	if err != nil {
		return nil, err
	}
	var r = &LiveAllGiftInfo{}
	if err = json.Unmarshal(resp.Data, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// CommentGetCount 获取评论总数
//
// oid: 对应类型的ID
//
// tp: 类型。https://github.com/SocialSisterYi/bilibili-API-collect/tree/master/comment#%E8%AF%84%E8%AE%BA%E5%8C%BA%E7%B1%BB%E5%9E%8B%E4%BB%A3%E7%A0%81
func (c *CommClient) CommentGetCount(oid int64, tp int) (int, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/v2/reply/count",
		"GET",
		map[string]string{
			"oid":  strconv.FormatInt(oid, 10),
			"type": strconv.Itoa(tp),
		},
	)
	if err != nil {
		return -1, err
	}
	var r struct {
		Count int `json:"count"`
	}
	if err = json.Unmarshal(resp.Data, &r); err != nil {
		return -1, err
	}
	return r.Count, nil
}

// CommentGetMain 获取评论区内容
//
// oid: 对应类型的ID
//
// tp: 类型。https://github.com/SocialSisterYi/bilibili-API-collect/tree/master/comment#%E8%AF%84%E8%AE%BA%E5%8C%BA%E7%B1%BB%E5%9E%8B%E4%BB%A3%E7%A0%81
//
// mode: 排序方式
//
// 0 3：仅按热度
//
// 1：按热度+按时间
//
// 2：仅按时间
//
// next: 评论页选择 按热度时：热度顺序页码（0为第一页） 按时间时：时间倒序楼层号
//
// ps: 每页项数
//
// 具体用法请看测试样例
func (c *CommClient) CommentGetMain(oid int64, tp int, mode int, next int, ps int) (*CommentMain, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/v2/reply/main",
		"GET",
		map[string]string{
			"oid":  strconv.FormatInt(oid, 10),
			"type": strconv.Itoa(tp),
			"mode": strconv.Itoa(mode),
			"next": strconv.Itoa(next),
			"ps":   strconv.Itoa(ps),
		},
	)
	if err != nil {
		return nil, err
	}
	var r = &CommentMain{}
	if err = json.Unmarshal(resp.Data, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// CommentGetReply 获取指定评论和二级回复
//
//
// oid: 对应类型的ID
//
// tp: 类型。https://github.com/SocialSisterYi/bilibili-API-collect/tree/master/comment#%E8%AF%84%E8%AE%BA%E5%8C%BA%E7%B1%BB%E5%9E%8B%E4%BB%A3%E7%A0%81
//
// root: 目标一级评论rpid
//
// pn: 二级评论页码 从1开始
//
// ps: 二级评论每页项数 定义域：1-49
func (c *CommClient) CommentGetReply(oid int64, tp int, root int64, pn int, ps int) (*CommentReply, error) {
	resp, err := c.RawParse(
		BiliApiURL,
		"x/v2/reply/reply",
		"GET",
		map[string]string{
			"oid":  strconv.FormatInt(oid, 10),
			"type": strconv.Itoa(tp),
			"root": strconv.FormatInt(root, 10),
			"pn":   strconv.Itoa(pn),
			"ps":   strconv.Itoa(ps),
		},
	)
	if err != nil {
		return nil, err
	}
	var r = &CommentReply{}
	if err = json.Unmarshal(resp.Data, &r); err != nil {
		return nil, err
	}
	return r, nil
}
