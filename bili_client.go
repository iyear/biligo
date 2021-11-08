package biligo

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/iyear/biligo/internal/util"
	"github.com/iyear/biligo/proto/dm"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type BiliClient struct {
	Me   *Account
	auth *CookieAuth

	*baseClient
}

type CookieAuth struct {
	DedeUserID      string // DedeUserID
	DedeUserIDCkMd5 string // DedeUserID__ckMd5
	SESSDATA        string // SESSDATA
	BiliJCT         string // bili_jct
}

type BiliSetting struct {
	// Cookie
	Auth *CookieAuth
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

// NewBiliClient
//
// 带有账户Cookie的Client，用于访问私人操作API
func NewBiliClient(setting *BiliSetting) (*BiliClient, error) {
	if setting.Auth == nil {
		return nil, errors.New("auth cannot be nil")
	}

	bili := &BiliClient{
		auth: setting.Auth,
		baseClient: newBaseClient(&baseSetting{
			Client:    setting.Client,
			DebugMode: setting.DebugMode,
			UserAgent: setting.UserAgent,
			Prefix:    "BiliClient ",
		}),
	}

	account, err := bili.GetMe()
	if err != nil {
		return nil, err
	}

	bili.Me = account

	return bili, nil
}

// GetMe
//
// 获取个人基本信息
func (b *BiliClient) GetMe() (*Account, error) {
	resp, err := b.RawParse(BiliApiURL,
		"x/member/web/account",
		"GET",
		nil,
	)
	if err != nil {
		return nil, err
	}
	var account *Account
	if err = json.Unmarshal(resp.Data, &account); err != nil {
		return nil, err
	}
	return account, nil
}

// SetClient
//
// 设置Client,可以用来更换代理等操作
func (b *BiliClient) SetClient(client *http.Client) {
	b.client = client
}

// SetUA
//
// 设置UA
func (b *BiliClient) SetUA(ua string) {
	b.ua = ua
}

// Raw
//
// base末尾带/
func (b *BiliClient) Raw(base, endpoint, method string, payload map[string]string) ([]byte, error) {
	raw, err := b.raw(base, endpoint, method, payload,
		func(d *url.Values) {
			switch method {
			case "POST":
				d.Add("csrf", b.auth.BiliJCT)
			}
		},
		func(r *http.Request) {
			r.Header.Add("Cookie", fmt.Sprintf("DedeUserID=%s;SESSDATA=%s;DedeUserID__ckMd5=%s",
				b.auth.DedeUserID, b.auth.SESSDATA, b.auth.DedeUserIDCkMd5))
		})
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// RawParse
//
// base末尾带/
func (b *BiliClient) RawParse(base, endpoint, method string, payload map[string]string) (*Response, error) {
	raw, err := b.Raw(base, endpoint, method, payload)
	if err != nil {
		return nil, err
	}
	return b.parse(raw)
}

// Upload 上传文件
//
// base末尾带/
func (b *BiliClient) Upload(base, endpoint string, payload map[string]string, files []*FileUpload) ([]byte, error) {
	raw, err := b.upload(base, endpoint, payload, files, func(m *multipart.Writer) error {
		return m.WriteField("csrf", b.auth.BiliJCT)
	}, func(r *http.Request) {
		r.Header.Add("Cookie", fmt.Sprintf("DedeUserID=%s;SESSDATA=%s;DedeUserID__ckMd5=%s",
			b.auth.DedeUserID, b.auth.SESSDATA, b.auth.DedeUserIDCkMd5))
	})
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// UploadParse 上传文件
//
// base末尾带/
func (b *BiliClient) UploadParse(base, endpoint string, payload map[string]string, files []*FileUpload) (*Response, error) {
	raw, err := b.Upload(base, endpoint, payload, files)
	if err != nil {
		return nil, err
	}
	return b.parse(raw)
}

// GetCookieAuth
//
// 获取Cookie信息
func (b *BiliClient) GetCookieAuth() *CookieAuth {
	return b.auth
}

// GetNavInfo
//
// 获取我的导航栏信息(大部分的用户信息都在这里了)
func (b *BiliClient) GetNavInfo() (*NavInfo, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/web-interface/nav",
		"GET",
		nil,
	)
	if err != nil {
		return nil, err
	}
	var info *NavInfo
	if err = json.Unmarshal(resp.Data, &info); err != nil {
		return nil, err
	}
	return info, nil
}

// GetNavStat
//
// 获取我的状态数
func (b *BiliClient) GetNavStat() (*NavStat, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/web-interface/nav/stat",
		"GET",
		nil,
	)
	if err != nil {
		return nil, err
	}
	var stat *NavStat
	if err = json.Unmarshal(resp.Data, &stat); err != nil {
		return nil, err
	}
	return stat, nil
}

// GetExpRewardStat
//
// 获取每日经验任务和认证相关任务的完成情况
func (b *BiliClient) GetExpRewardStat() (*ExpRewardStat, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/member/web/exp/reward",
		"GET",
		nil,
	)
	if err != nil {
		return nil, err
	}
	var stat *ExpRewardStat
	if err = json.Unmarshal(resp.Data, &stat); err != nil {
		return nil, err
	}
	return stat, nil
}

// GetExpCoinReward
//
// 获取每日投币获得经验数，不存在延迟 上限50经验
func (b *BiliClient) GetExpCoinReward() (int, error) {
	resp, err := b.Raw(
		BiliMainURL,
		"plus/account/exp.php",
		"GET",
		nil,
	)
	if err != nil {
		return -1, err
	}
	var tResp struct {
		Response
		Number int `json:"number"`
	}
	if err = json.Unmarshal(resp, &tResp); err != nil {
		return -1, err
	}
	if tResp.Code != 0 {
		return -1, fmt.Errorf("(%d) %s", tResp.Code, tResp.Message)
	}
	return tResp.Number, nil
}

// GetVipStat
//
// 获取大会员信息
func (b *BiliClient) GetVipStat() (*VipStat, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/vip/web/user/info",
		"GET",
		nil,
	)
	if err != nil {
		return nil, err
	}
	var info *VipStat
	if err = json.Unmarshal(resp.Data, &info); err != nil {
		return nil, err
	}
	return info, nil
}

// GetAccountSafetyStat
//
// 获取账户安全情况
func (b *BiliClient) GetAccountSafetyStat() (*AccountSafetyStat, error) {
	resp, err := b.RawParse(
		BiliPassportURL,
		"web/site/user/info",
		"GET",
		nil,
	)
	if err != nil {
		return nil, err
	}
	var info *AccountSafetyStat
	if err = json.Unmarshal(resp.Data, &info); err != nil {
		return nil, err
	}
	return info, nil
}

// GetRealNameStat
//
// 获取账户实名认证状态
func (b *BiliClient) GetRealNameStat() (bool, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/member/realname/status",
		"GET",
		nil,
	)
	if err != nil {
		return false, err
	}
	var flag struct {
		Status int `json:"status,omitempty"`
	}
	if err = json.Unmarshal(resp.Data, &flag); err != nil {
		return false, err
	}
	return flag.Status == 1, nil
}

// GetRealNameInfo
//
// 获取账户实名认证详细信息
func (b *BiliClient) GetRealNameInfo() (*RealNameInfo, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/member/realname/apply/status",
		"GET",
		nil,
	)
	if err != nil {
		return nil, err
	}
	var info *RealNameInfo
	if err = json.Unmarshal(resp.Data, &info); err != nil {
		return nil, err
	}
	return info, nil
}

// FollowingsGetMy 获取我的关注列表
//
// mid数组
func (b *BiliClient) FollowingsGetMy() ([]int64, error) {
	resp, err := b.RawParse(
		BiliVcURL,
		"feed/v1/feed/get_attention_list",
		"GET",
		map[string]string{},
	)
	if err != nil {
		return []int64{}, err
	}

	var mids struct {
		List []int64 `json:"list"`
	}
	if err = json.Unmarshal(resp.Data, &mids); err != nil {
		return []int64{}, err
	}
	return mids.List, nil
}

// FollowingsGetMyDetail 获取我的详细的关注列表，获取他人的请使用 CommClient
//
// pn 页码
//
// ps 每页项数，最大50
//
// order 1:最常访问 2:最近关注
func (b *BiliClient) FollowingsGetMyDetail(pn int, ps int, order int) (*FollowingsDetail, error) {
	var o = map[int]string{
		1: "attention",
		2: "",
	}
	resp, err := b.RawParse(
		BiliApiURL,
		"x/relation/followings",
		"GET",
		map[string]string{
			"vmid":       strconv.FormatInt(b.Me.MID, 10),
			"pn":         strconv.Itoa(pn),
			"ps":         strconv.Itoa(ps),
			"order_type": o[order],
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

// GetCoinLogs
//
// 获取最近一周的硬币变化情况
func (b *BiliClient) GetCoinLogs() ([]*CoinLog, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/member/web/coin/log",
		"GET",
		nil,
	)
	if err != nil {
		return nil, err
	}
	var logs struct {
		List []*CoinLog `json:"list,omitempty"`
	}
	if err = json.Unmarshal(resp.Data, &logs); err != nil {
		return nil, err
	}
	return logs.List, nil
}

// GetRelationStat
//
// 获取关系状态数
func (b *BiliClient) GetRelationStat(mid int64) (*RelationStat, error) {
	resp, err := b.RawParse(
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

// GetUpStat
//
// 获取UP主状态数，该接口需要任意用户登录，否则不会返回任何数据
func (b *BiliClient) GetUpStat(mid int64) (*UpStat, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/space/upstat",
		"GET",
		map[string]string{
			"mid": strconv.FormatInt(mid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var stat *UpStat
	if err = json.Unmarshal(resp.Data, &stat); err != nil {
		return nil, err
	}
	return stat, nil
}

// GetMsgUnread 获取未读消息数
//
// Link:https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/message/msg.md#%E9%80%9A%E7%9F%A5%E6%B6%88%E6%81%AF
func (b *BiliClient) GetMsgUnread() (*MsgUnRead, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/msgfeed/unread",
		"GET",
		nil,
	)
	if err != nil {
		return nil, err
	}
	var unread *MsgUnRead
	if err = json.Unmarshal(resp.Data, &unread); err != nil {
		return nil, err
	}
	return unread, nil
}

// SpaceSetTopArchive 设置置顶视频
//
// reason 备注 最大40字符
func (b *BiliClient) SpaceSetTopArchive(aid int64, reason string) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/space/top/arc/set",
		"POST",
		map[string]string{
			"aid":    strconv.FormatInt(aid, 10),
			"reason": reason,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// SpaceCancelTopArchive
//
// 取消置顶视频
func (b *BiliClient) SpaceCancelTopArchive() error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/space/top/arc/cancel",
		"POST",
		nil,
	)
	return err
}

// SpaceAddMasterpieces 添加代表作 上限为3个稿件
//
// reason 备注 最大40字符
func (b *BiliClient) SpaceAddMasterpieces(aid int64, reason string) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/space/masterpiece/add",
		"POST",
		map[string]string{
			"aid":    strconv.FormatInt(aid, 10),
			"reason": reason,
		},
	)
	return err
}

// SpaceCancelMasterpiece
//
// 取消代表作视频
func (b *BiliClient) SpaceCancelMasterpiece(aid int64) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/space/masterpiece/cancel",
		"POST",
		map[string]string{
			"aid": strconv.FormatInt(aid, 10),
		},
	)
	return err
}

// SpaceSetTags 设置用户个人TAG TAG里不要包含逗号会被当作分隔符
//
// 经过测试好像可以建54个TAG 各TAG长度小于10字符
//
// 删除TAG留空或省去即可 重复TAG会重复创建
//
// 感觉这个功能已经被废弃了
func (b *BiliClient) SpaceSetTags(tags []string) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/space/acc/tags/set",
		"POST",
		map[string]string{
			"tags": util.StringSliceToString(tags, ","),
		},
	)
	return err
}

// SpaceSetNotice 修改公告内容
//
// 删除公告留空即可 少于150字符
func (b *BiliClient) SpaceSetNotice(notice string) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/space/notice/set",
		"POST",
		map[string]string{
			"notice": notice,
		},
	)
	return err
}

// SpaceGetMyLastPlayGame
//
// 获取我的空间近期玩的游戏
func (b *BiliClient) SpaceGetMyLastPlayGame() ([]*SpaceGame, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/space/lastplaygame",
		"GET",
		map[string]string{
			"mid": strconv.FormatInt(b.Me.MID, 10),
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

// SpaceGetMyLastVideoCoin
//
// 获取我的最近投币的视频明细
func (b *BiliClient) SpaceGetMyLastVideoCoin() ([]*SpaceVideoCoin, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/space/coin/video",
		"GET",
		map[string]string{
			"vmid": strconv.FormatInt(b.Me.MID, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var info []*SpaceVideoCoin
	if err = json.Unmarshal(resp.Data, &info); err != nil {
		return nil, err
	}
	if info == nil {
		return []*SpaceVideoCoin{}, nil
	}
	return info, nil
}

// ChanGetMy
//
// 获取我的频道列表
func (b *BiliClient) ChanGetMy() (*ChannelList, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/space/channel/list",
		"GET",
		map[string]string{
			"mid": strconv.FormatInt(b.Me.MID, 10),
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

// ChanAdd
//
// 创建频道 创建成功后会返回新建频道的id
//
// name 频道名
//
// intro 频道介绍
func (b *BiliClient) ChanAdd(name string, intro string) (int64, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/space/channel/add",
		"POST",
		map[string]string{
			"name":  name,
			"intro": intro,
		},
	)
	if err != nil {
		return -1, err
	}
	var result struct {
		CID int64 `json:"cid,omitempty"`
	}
	if err = json.Unmarshal(resp.Data, &result); err != nil {
		return -1, err
	}
	return result.CID, nil
}

// ChanEdit 编辑频道
//
// cid 频道id
//
// name 频道名
//
// intro 频道介绍
func (b *BiliClient) ChanEdit(cid int64, name string, intro string) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/space/channel/edit",
		"POST",
		map[string]string{
			"cid":   strconv.FormatInt(cid, 10),
			"name":  name,
			"intro": intro,
		},
	)
	return err
}

// ChanDel 删除频道
//
// cid 频道id
func (b *BiliClient) ChanDel(cid int64) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/space/channel/del",
		"POST",
		map[string]string{
			"cid": strconv.FormatInt(cid, 10),
		},
	)
	return err
}

// ChanAddVideo 频道添加视频，返回添加错误的视频aid
//
// 仅能添加自己是UP主的视频
//
// 如添加多个视频，仅会添加正确的
func (b *BiliClient) ChanAddVideo(cid int64, aids []int64) ([]int64, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/space/channel/video/add",
		"POST",
		map[string]string{
			"cid":  strconv.FormatInt(cid, 10),
			"aids": util.Int64SliceToString(aids, ","),
		},
	)

	if err != nil {
		return nil, err
	}
	var result []int64
	if err = json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}
	// 完成后需要使用接口「查询用户频道中的视频」刷新
	if _, err := b.ChanGetMyVideo(cid, 1, 1); err != nil {
		return nil, err
	}
	return result, nil
}

// ChanDelVideo
//
// 删除频道视频
func (b *BiliClient) ChanDelVideo(cid int64, aid int64) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/space/channel/video/del",
		"POST",
		map[string]string{
			"cid": strconv.FormatInt(cid, 10),
			"aid": strconv.FormatInt(aid, 10),
		},
	)
	if err != nil {
		return err
	}
	// 完成后需要使用接口「查询用户频道中的视频」刷新
	if _, err = b.ChanGetMyVideo(cid, 1, 1); err != nil {
		return err
	}
	return nil
}

// ChanSetVideoSort 调整频道视频顺序
//
// to 视频排序倒数位置 1为列表底部，视频总数为首端，与显示顺序恰好相反
func (b *BiliClient) ChanSetVideoSort(cid int64, aid int64, to int) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/space/channel/video/sort",
		"POST",
		map[string]string{
			"cid": strconv.FormatInt(cid, 10),
			"aid": strconv.FormatInt(aid, 10),
			"to":  strconv.Itoa(to),
		},
	)
	return err
}

// ChanHasInvalidVideo
//
// 检查频道是否有失效视频，若有以错误形式返回(错误码:53005)
//
// 若Err为nil则没有无效视频
func (b *BiliClient) ChanHasInvalidVideo(cid int64) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/space/channel/video/check",
		"GET",
		map[string]string{
			"cid": strconv.FormatInt(cid, 10),
		},
	)
	return err
}

// ChanGetMyVideo
//
// 获取我的频道视频
//
// cid 频道ID
//
// pn 页码
//
// ps 每页项数
func (b *BiliClient) ChanGetMyVideo(cid int64, pn int, ps int) (*ChanVideo, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/space/channel/video",
		"GET",
		map[string]string{
			"mid": strconv.FormatInt(b.Me.MID, 10),
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

// FavGetMy
//
// 获取我的收藏夹列表
func (b *BiliClient) FavGetMy() (*FavoritesList, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/v3/fav/folder/created/list-all",
		"GET",
		map[string]string{
			"up_mid": strconv.FormatInt(b.Me.MID, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var list = &FavoritesList{}
	if err = json.Unmarshal(resp.Data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// FavGetDetail
//
// 获取收藏夹详细信息
func (b *BiliClient) FavGetDetail(mlid int64) (*FavDetail, error) {
	resp, err := b.RawParse(
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
	var detail = &FavDetail{}
	if err = json.Unmarshal(resp.Data, &detail); err != nil {
		return nil, err
	}
	return detail, nil
}

// FavAdd 新建收藏夹
//
// title 收藏夹标题
//
// intro 收藏夹简介
//
// privacy 是否私密 true:私密 false:公开
//
// cover 封面图url，会审核，不需要请留空
func (b *BiliClient) FavAdd(title string, intro string, privacy bool, cover string) (*FavDetail, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/v3/fav/folder/add",
		"POST",
		map[string]string{
			"title":   title,
			"intro":   intro,
			"privacy": util.IF(privacy, "1", "0").(string),
			"cover":   cover,
		},
	)
	if err != nil {
		return nil, err
	}
	var detail = &FavDetail{}
	if err = json.Unmarshal(resp.Data, &detail); err != nil {
		return nil, err
	}
	return detail, nil
}

// FavEdit
//
// 编辑收藏夹 参数注释与 FavAdd 相同
func (b *BiliClient) FavEdit(mlid int64, title string, intro string, privacy bool, cover string) (*FavDetail, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/v3/fav/folder/edit",
		"POST",
		map[string]string{
			"media_id": strconv.FormatInt(mlid, 10),
			"title":    title,
			"intro":    intro,
			"privacy":  util.IF(privacy, "1", "0").(string),
			"cover":    cover,
		},
	)
	if err != nil {
		return nil, err
	}
	var detail = &FavDetail{}
	if err = json.Unmarshal(resp.Data, &detail); err != nil {
		return nil, err
	}
	return detail, nil
}

// FavDel
//
// 删除收藏夹，传入需要删除的mlid数组
func (b *BiliClient) FavDel(mlids []int64) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/v3/fav/folder/del",
		"POST",
		map[string]string{
			"media_ids": util.Int64SliceToString(mlids, ","),
		},
	)
	return err
}

// FavGetRes
//
// 获取收藏夹全部内容id
func (b *BiliClient) FavGetRes(mlid int64) ([]*FavRes, error) {
	resp, err := b.RawParse(
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
	var r = make([]*FavRes, 0)
	if err = json.Unmarshal(resp.Data, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// FavGetResDetail 获取收藏夹内容详细内容，带过滤功能
//
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
func (b *BiliClient) FavGetResDetail(mlid int64, tid int, keyword string, order string, tp int, pn int, ps int) (*FavResDetail, error) {
	resp, err := b.RawParse(
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

// FavCopyRes 收藏夹批量复制内容
//
// 例子：[]string{"21822819:2", "21918689:2", "22288065:2"}
//
// from 源收藏夹mlid
//
// to 目标收藏夹mlid
//
// mid 当前用户mid
//
// resources 目标内容id列表
//
// 字符串数组成员格式：{内容id}:{内容类型}
// 例子：
//
// 类型：
// 2：视频稿件
// 12：音频
// 21：视频合集
//
// 内容id：
//
// 视频稿件：视频稿件avid
//
// 音频：音频auid
//
// 视频合集：视频合集id
func (b *BiliClient) FavCopyRes(from int64, to int64, mid int64, resources []string) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/v3/fav/resource/copy",
		"POST",
		map[string]string{
			"src_media_id": strconv.FormatInt(from, 10),
			"tar_media_id": strconv.FormatInt(to, 10),
			"mid":          strconv.FormatInt(mid, 10),
			"resources":    util.StringSliceToString(resources, ","),
			"platform":     "web",
		},
	)
	return err
}

// FavMoveRes 收藏夹批量移动内容
//
// 参数说明同 FavCopyRes
func (b *BiliClient) FavMoveRes(from int64, to int64, mid int64, resources []string) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/v3/fav/resource/move",
		"POST",
		map[string]string{
			"src_media_id": strconv.FormatInt(from, 10),
			"tar_media_id": strconv.FormatInt(to, 10),
			"mid":          strconv.FormatInt(mid, 10),
			"resources":    util.StringSliceToString(resources, ","),
			"platform":     "web",
		},
	)
	return err
}

// FavDelRes 收藏夹批量删除内容
//
// resources 同 FavCopyRes
func (b *BiliClient) FavDelRes(mlid int64, resources []string) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/v3/fav/resource/batch-del",
		"POST",
		map[string]string{
			"media_id":  strconv.FormatInt(mlid, 10),
			"resources": util.StringSliceToString(resources, ","),
		},
	)
	return err
}

// FavCleanRes
//
// 清除收藏夹失效内容
func (b *BiliClient) FavCleanRes(mlid int64) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/v3/fav/resource/clean",
		"POST",
		map[string]string{
			"media_id": strconv.FormatInt(mlid, 10),
		},
	)
	return err
}

// SignUpdate
//
// 更新个性签名
func (b *BiliClient) SignUpdate(sign string) error {
	_, err := b.RawParse(BiliApiURL,
		"x/member/web/sign/update",
		"POST",
		map[string]string{"user_sign": sign},
	)
	return err
}

// VideoAddLike
//
// 点赞稿件
func (b *BiliClient) VideoAddLike(aid int64, like bool) error {
	_, err := b.RawParse(BiliApiURL,
		"x/web-interface/archive/like",
		"POST",
		map[string]string{
			"aid":  strconv.FormatInt(aid, 10),
			"like": util.IF(like, "1", "2").(string),
		},
	)
	return err
}

// VideoIsLiked
//
// 获取稿件是否被点赞
func (b *BiliClient) VideoIsLiked(aid int64) (bool, error) {
	resp, err := b.RawParse(BiliApiURL,
		"x/web-interface/archive/has/like",
		"GET",
		map[string]string{
			"aid": strconv.FormatInt(aid, 10),
		},
	)
	if err != nil {
		return false, err
	}

	var liked int
	if err = json.Unmarshal(resp.Data, &liked); err != nil {
		return false, err
	}

	return liked == 1, nil
}

// VideoAddCoins 视频投币
//
// aid 视频avid
//
// num 投币数量,上限为2
//
// like 是否附加点赞
func (b *BiliClient) VideoAddCoins(aid int64, num int, like bool) error {
	_, err := b.RawParse(BiliApiURL,
		"x/web-interface/coin/add",
		"POST",
		map[string]string{
			"aid":         strconv.FormatInt(aid, 10),
			"multiply":    strconv.Itoa(num),
			"select_like": util.IF(like, "1", "0").(string),
		},
	)
	return err
}

// VideoIsAddedCoins
//
// 返回投币数
func (b *BiliClient) VideoIsAddedCoins(aid int64) (int, error) {
	resp, err := b.RawParse(BiliApiURL,
		"x/web-interface/archive/coins",
		"GET",
		map[string]string{
			"aid": strconv.FormatInt(aid, 10),
		},
	)
	if err != nil {
		return -1, err
	}

	var coins struct {
		Multiply int `json:"multiply,omitempty"` // 投币枚数,未投币为0
	}
	if err = json.Unmarshal(resp.Data, &coins); err != nil {
		return -1, err
	}

	return coins.Multiply, nil
}

// VideoSetFavour 收藏视频，返回 [是否为未关注用户收藏] 的布尔值
//
// addMediaLists 需要加入的收藏夹id 非必须 传入空切片或nil留空
//
// delMediaLists 需要取消的收藏夹id 非必须 传入空切片或nil留空
func (b *BiliClient) VideoSetFavour(aid int64, addLists []int64, delLists []int64) (bool, error) {
	resp, err := b.RawParse(BiliApiURL,
		"medialist/gateway/coll/resource/deal",
		"POST",
		map[string]string{
			"rid":           strconv.FormatInt(aid, 10),
			"type":          "2",
			"add_media_ids": util.Int64SliceToString(addLists, ","),
			"del_media_ids": util.Int64SliceToString(delLists, ","),
		},
	)
	if err != nil {
		return false, err
	}

	var prompt struct {
		Prompt bool `json:"prompt,omitempty"` // 是否为未关注用户收藏
	}
	if err = json.Unmarshal(resp.Data, &prompt); err != nil {
		return false, err
	}

	return prompt.Prompt, nil
}

// VideoIsFavoured
//
// 返回 是否被收藏
func (b *BiliClient) VideoIsFavoured(aid int64) (bool, error) {
	resp, err := b.RawParse(BiliApiURL,
		"x/v2/fav/video/favoured",
		"GET",
		map[string]string{
			"aid": strconv.FormatInt(aid, 10),
		},
	)
	if err != nil {
		return false, err
	}

	var favour struct {
		Count    int  `json:"count,omitempty"`    // 作用尚不明确
		Favoured bool `json:"favoured,omitempty"` // true：已收藏  false：未收藏
	}
	if err = json.Unmarshal(resp.Data, &favour); err != nil {
		return false, err
	}

	return favour.Favoured, nil
}

// VideoTriple
//
// 返回是否点赞成功、投币成功、收藏成功和投币枚数
func (b *BiliClient) VideoTriple(aid int64) (like, coin, favour bool, multiply int, e error) {
	resp, err := b.RawParse(BiliApiURL,
		"x/web-interface/archive/like/triple",
		"POST",
		map[string]string{
			"aid": strconv.FormatInt(aid, 10),
		},
	)
	if err != nil {
		return false, false, false, -1, err
	}

	var triple struct {
		Like     bool `json:"like,omitempty"`     // 是否点赞成功
		Coin     bool `json:"coin,omitempty"`     // 是否投币成功
		Fav      bool `json:"fav,omitempty"`      // 是否收藏成功
		Multiply int  `json:"multiply,omitempty"` // 投币枚数
	}
	if err = json.Unmarshal(resp.Data, &triple); err != nil {
		return false, false, false, -1, err
	}

	return triple.Like, triple.Coin, triple.Fav, triple.Multiply, nil
}

// VideoShare
//
// 完成分享并返回该视频当前分享数
func (b *BiliClient) VideoShare(aid int64) (int, error) {
	resp, err := b.RawParse(BiliApiURL,
		"x/web-interface/share/add",
		"POST",
		map[string]string{
			"aid": strconv.FormatInt(aid, 10),
		},
	)
	if err != nil {
		return -1, err
	}

	var shareNum int
	if err = json.Unmarshal(resp.Data, &shareNum); err != nil {
		return -1, err
	}

	return shareNum, nil
}

// VideoGetInfo
//
// 返回视频详细信息，数据较多，可以使用单独的接口获取部分数据
func (b *BiliClient) VideoGetInfo(aid int64) (*VideoInfo, error) {
	resp, err := b.RawParse(BiliApiURL,
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

// VideoReportProgress 视频观看进度上报
//
// cid用于分P标识,progress为观看进度(单位为秒)
//
// 不是心跳包，应该就是个历史记录和下次播放自动跳转的功能，一般是关闭当前视频页时请求
func (b *BiliClient) VideoReportProgress(aid int64, cid int64, progress int64) error {
	_, err := b.RawParse(BiliApiURL,
		"x/v2/history/report",
		"POST",
		map[string]string{
			"aid":      strconv.FormatInt(aid, 10),
			"cid":      strconv.FormatInt(cid, 10),
			"progress": strconv.FormatInt(progress, 10),
		},
	)
	return err
}

// VideoGetPlayURL 获取视频取流地址
//
// 所有参数、返回信息和取流方法的说明请直接前往：https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/video/videostream_url.md
func (b *BiliClient) VideoGetPlayURL(aid int64, cid int64, qn int, fnval int) (*VideoPlayURLResult, error) {
	resp, err := b.RawParse(
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

// VideoHeartBeat 视频心跳包上报
//
// 默认间隔15秒一次，不要过慢或过快上报，控制好时间间隔
//
// playedTime 已播放时间 单位为秒 默认为0
//
// 包含了 VideoReportProgress 的功能(应该)
func (b *BiliClient) VideoHeartBeat(aid int64, cid int64, playedTime int64) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/click-interface/web/heartbeat",
		"POST",
		map[string]string{
			"aid":         strconv.FormatInt(aid, 10),
			"cid":         strconv.FormatInt(cid, 10),
			"mid":         strconv.FormatInt(b.Me.MID, 10),
			"start_ts":    strconv.FormatInt(util.GetCST8Time(time.Now()).Unix(), 10),
			"played_time": strconv.FormatInt(playedTime, 10),
		},
	)
	return err
}

// VideoGetTags
//
// 获取稿件Tags
func (b *BiliClient) VideoGetTags(aid int64) ([]*VideoTag, error) {
	resp, err := b.RawParse(BiliApiURL,
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

// VideoLikeTag 点赞视频的TAG
//
// 重复请求为取消
func (b *BiliClient) VideoLikeTag(aid int64, tagID int64) error {
	_, err := b.RawParse(BiliApiURL,
		"x/tag/archive/like2",
		"POST",
		map[string]string{
			"aid":    strconv.FormatInt(aid, 10),
			"tag_id": strconv.FormatInt(tagID, 10),
		},
	)
	return err
}

// VideoHateTag 点踩视频的TAG
//
// 重复请求为取消
func (b *BiliClient) VideoHateTag(aid int64, tagID int64) error {
	_, err := b.RawParse(BiliApiURL,
		"x/tag/archive/hate2",
		"POST",
		map[string]string{
			"aid":    strconv.FormatInt(aid, 10),
			"tag_id": strconv.FormatInt(tagID, 10),
		},
	)
	return err
}

// CommentSend 发送评论
//
// oid: 对应类型的ID
//
// tp: 类型。https://github.com/SocialSisterYi/bilibili-API-collect/tree/master/comment#%E8%AF%84%E8%AE%BA%E5%8C%BA%E7%B1%BB%E5%9E%8B%E4%BB%A3%E7%A0%81
//
// content: 评论内容，最大1000字符 表情使用表情转义符
//
// platform: 平台标识 1：web端 2：安卓客户端 3：ios客户端 4：wp客户端
//
// root: 二级评论以上使用 没有填0
//
// parent: 二级评论同根评论id 大于二级评论为要回复的评论id
func (b *BiliClient) CommentSend(oid int64, tp int, content string, platform int, root int64, parent int64) (*CommentSend, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/v2/reply/add",
		"POST",
		map[string]string{
			"oid":      strconv.FormatInt(oid, 10),
			"type":     strconv.Itoa(tp),
			"root":     strconv.FormatInt(root, 10),
			"parent":   strconv.FormatInt(parent, 10),
			"ordering": "heat", // 暂时不知道作用
			"message":  content,
			"plat":     strconv.Itoa(platform),
		},
	)
	if err != nil {
		return nil, err
	}
	var r = &CommentSend{}
	if err = json.Unmarshal(resp.Data, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// CommentLike 点赞评论，点赞成功后会同时消去该评论的点踩
//
// oid,tp: 同 BiliClient.CommentSend
//
// rpid: 评论ID
//
// like: true为点赞，false为取消点赞
func (b *BiliClient) CommentLike(oid int64, tp int, rpid int64, like bool) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/v2/reply/action",
		"POST",
		map[string]string{
			"oid":    strconv.FormatInt(oid, 10),
			"type":   strconv.Itoa(tp),
			"rpid":   strconv.FormatInt(rpid, 10),
			"action": util.IF(like, "1", "0").(string),
		},
	)
	return err
}

// CommentHate 点踩评论，点踩成功后会同时消去该评论的点赞
//
// oid,tp: 同 BiliClient.CommentSend
//
// rpid: 评论ID
//
// like: true为点踩，false为取消点踩
func (b *BiliClient) CommentHate(oid int64, tp int, rpid int64, hate bool) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/v2/reply/hate",
		"POST",
		map[string]string{
			"oid":    strconv.FormatInt(oid, 10),
			"type":   strconv.Itoa(tp),
			"rpid":   strconv.FormatInt(rpid, 10),
			"action": util.IF(hate, "1", "0").(string),
		},
	)
	return err
}

// CommentDel 删除评论 只能删除自己的评论，或自己管理的评论区下的评论
//
// oid,tp: 同 BiliClient.CommentSend
//
// rpid: 评论ID
func (b *BiliClient) CommentDel(oid int64, tp int, rpid int64) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/v2/reply/del",
		"POST",
		map[string]string{
			"oid":  strconv.FormatInt(oid, 10),
			"type": strconv.Itoa(tp),
			"rpid": strconv.FormatInt(rpid, 10),
		},
	)
	return err
}

// CommentSetTop 置顶评论 只能置顶自己管理的评论区中的一级评论
//
// oid,tp: 同 BiliClient.CommentSend
//
// rpid: 评论ID
//
// top: true为置顶，false为取消置顶
func (b *BiliClient) CommentSetTop(oid int64, tp int, rpid int64, top bool) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/v2/reply/top",
		"POST",
		map[string]string{
			"oid":    strconv.FormatInt(oid, 10),
			"type":   strconv.Itoa(tp),
			"rpid":   strconv.FormatInt(rpid, 10),
			"action": util.IF(top, "1", "0").(string),
		},
	)
	return err
}

// CommentReport 举报评论
//
// oid,tp: 同 BiliClient.CommentSend
//
// rpid: 评论ID
//
// reason: 参考 https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/comment/action.md#%E4%B8%BE%E6%8A%A5%E8%AF%84%E8%AE%BA
//
// content: 其他举报备注 reason=0时有效 不需要时留空
func (b *BiliClient) CommentReport(oid int64, tp int, rpid int64, reason int, content string) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/v2/reply/report",
		"POST",
		map[string]string{
			"oid":     strconv.FormatInt(oid, 10),
			"type":    strconv.Itoa(tp),
			"rpid":    strconv.FormatInt(rpid, 10),
			"reason":  strconv.Itoa(reason),
			"content": content,
		},
	)
	return err
}

// DanmakuGetHistoryIndex
//
// 获取历史弹幕日期，返回的日期代表有历史弹幕，用于请求历史弹幕
func (b *BiliClient) DanmakuGetHistoryIndex(cid int64, year int, month int) ([]string, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/v2/dm/history/index",
		"GET",
		map[string]string{
			"type":  "1",
			"oid":   strconv.FormatInt(cid, 10),
			"month": fmt.Sprintf("%04d-%02d", year, month),
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

// DanmakuGetHistory
//
// 获取历史弹幕
//
// date 历史日期 YYYY-MM-DD
func (b *BiliClient) DanmakuGetHistory(cid int64, date string) (*DanmakuResp, error) {
	resp, err := b.Raw(
		BiliApiURL,
		"x/v2/dm/web/history/seg.so",
		"GET",
		map[string]string{
			"type": "1",
			"oid":  strconv.FormatInt(cid, 10),
			"date": date,
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

// DanmakuPost 发送普通弹幕
//
// Link:https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/danmaku/action.md#%E5%8F%91%E9%80%81%E8%A7%86%E9%A2%91%E5%BC%B9%E5%B9%95
//
// tp：类型 1:视频弹幕
//
// aid 稿件avid
//
// cid 用于区分分P
//
// msg 弹幕内容 长度小于100字符
//
// progress 弹幕出现在视频内的时间 单位为毫秒
//
// color 弹幕颜色 十进制RGB888值 [默认为16777215（#FFFFFF）白色]
//
// fontsize 弹幕字号 默认为25
// 极小:12
// 超小:16
// 小:18
// 标准:25
// 大:36
// 超大:45
// 极大:64
//
// pool 弹幕池
// 0:普通池
// 1:字幕池
// 2:特殊池（代码/BAS弹幕）
// 默认为0
//
// mode 弹幕类型
// 1:普通弹幕
// 4:底部弹幕
// 5:顶部弹幕
// 7:高级弹幕
// 9:BAS弹幕（pool必须为2）
func (b *BiliClient) DanmakuPost(tp int, aid int64, cid int64, msg string, progress int64, color int, fontsize int, pool int, mode int) (*DanmakuPostResult, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/v2/dm/post",
		"POST",
		map[string]string{
			"type":     strconv.Itoa(tp),
			"oid":      strconv.FormatInt(cid, 10),
			"msg":      msg,
			"aid":      strconv.FormatInt(aid, 10),
			"progress": strconv.FormatInt(progress, 10),
			"color":    strconv.Itoa(color),
			"fontsize": strconv.Itoa(fontsize),
			"pool":     strconv.Itoa(pool),
			"mode":     strconv.Itoa(mode),
			"rnd":      strconv.FormatInt(util.GetCST8Time(time.Now()).UnixNano(), 10),
		},
	)
	if err != nil {
		return nil, err
	}

	var result *DanmakuPostResult
	if err = json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// DanmakuRecall 仅能撤回自己两分钟内的弹幕，且每天机会有限额
//
// Link:https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/danmaku/action.md#%E6%92%A4%E5%9B%9E%E5%BC%B9%E5%B9%95
//
// 成功后显示剩余次数的文本信息 如 "撤回成功，你还有2次撤回机会"
func (b *BiliClient) DanmakuRecall(cid int64, dmid uint64) (string, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/dm/recall",
		"POST",
		map[string]string{
			"cid":  strconv.FormatInt(cid, 10),
			"dmid": strconv.FormatUint(dmid, 10),
		},
	)
	if err != nil {
		return "", err
	}

	return resp.Message, nil
}

// DanmakuGetLikes 获取弹幕点赞数，一次可以获取多条弹幕
//
// Link:https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/danmaku/action.md#%E6%9F%A5%E8%AF%A2%E5%BC%B9%E5%B9%95%E7%82%B9%E8%B5%9E%E6%95%B0
//
// 返回一个map，key为dmid，value为相关信息
func (b *BiliClient) DanmakuGetLikes(cid int64, dmids []uint64) (map[uint64]*DanmakuGetLikesResult, error) {
	resp, err := b.RawParse(
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
		var d *DanmakuGetLikesResult
		if err = json.Unmarshal([]byte(gjson.Get(string(resp.Data), strconv.FormatUint(dmid, 10)).Raw), &d); err != nil {
			return nil, err
		}
		result[dmid] = d
	}
	return result, nil
}

// DanmakuLike 点赞弹幕
//
// Link:https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/danmaku/action.md#%E7%82%B9%E8%B5%9E%E5%BC%B9%E5%B9%95
//
// op 1:点赞 2:取消点赞
func (b *BiliClient) DanmakuLike(cid int64, dmid uint64, op int) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/v2/dm/thumbup/add",
		"POST",
		map[string]string{
			"oid":  strconv.FormatInt(cid, 10),
			"dmid": strconv.FormatUint(dmid, 10),
			"op":   strconv.Itoa(op),
		},
	)
	return err
}

// DanmakuReport 举报弹幕
//
// Link:https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/danmaku/action.md#%E4%B8%BE%E6%8A%A5%E5%BC%B9%E5%B9%95
//
// reason
// 1:违法违禁
// 2:色情低俗
// 3:赌博诈骗
// 4:人身攻击
// 5:侵犯隐私
// 6:垃圾广告
// 7:引战
// 8:剧透
// 9:恶意刷屏
// 10:视频无关
// 11:其他
// 12:青少年不良
//
// content 其他举报备注(可空) reason=11时有效
func (b *BiliClient) DanmakuReport(cid int64, dmid uint64, reason int, content string) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/dm/report/add",
		"POST",
		map[string]string{
			"cid":     strconv.FormatInt(cid, 10),
			"dmid":    strconv.FormatUint(dmid, 10),
			"reason":  strconv.Itoa(reason),
			"content": content,
		},
	)
	return err
}

// DanmakuEditState 保护&删除弹幕，只能操作自己的稿件或有骑士权限的稿件
//
// Link:https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/danmaku/action.md#%E4%BF%9D%E6%8A%A4%E5%88%A0%E9%99%A4%E5%BC%B9%E5%B9%95
//
// tp 弹幕类选择 1:视频弹幕
//
// dmids 弹幕dmid数组
//
// state 操作代码
// 1:删除弹幕
// 2:弹幕保护
// 3:取消保护
func (b *BiliClient) DanmakuEditState(tp int, cid int64, dmids []uint64, state int) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/v2/dm/edit/state",
		"POST",
		map[string]string{
			"type":  strconv.Itoa(tp),
			"oid":   strconv.FormatInt(cid, 10),
			"dmids": util.Uint64SliceToString(dmids, ","),
			"state": strconv.Itoa(state),
		},
	)
	return err
}

// DanmakuEditPool 修改字幕池，只能操作自己的稿件或有骑士权限的稿件
//
// Link:https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/danmaku/action.md#%E4%BF%AE%E6%94%B9%E5%AD%97%E5%B9%95%E6%B1%A0
//
// tp 弹幕类选择 1:视频弹幕
//
// dmids 弹幕dmid数组
//
// pool 操作代码
// 0:移出字幕池
// 1:移入字幕池
func (b *BiliClient) DanmakuEditPool(tp int, cid int64, dmids []uint64, pool int) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/v2/dm/edit/pool",
		"POST",
		map[string]string{
			"type":  strconv.Itoa(tp),
			"oid":   strconv.FormatInt(cid, 10),
			"dmids": util.Uint64SliceToString(dmids, ","),
			"pool":  strconv.Itoa(pool),
		},
	)
	return err
}

// DanmakuCommandPost 发送互动弹幕，只能在自己的视频发
//
// Link:https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/danmaku/action.md#%E5%8F%91%E9%80%81%E4%BA%92%E5%8A%A8%E5%BC%B9%E5%B9%95
//
// tp 互动弹幕类型
// 1:UP主头像弹幕
// 2:关联视频弹幕
// 5:视频内嵌引导关注按钮
//
// aid 稿件avid
//
// cid 视频cid
//
// progress 弹幕出现在视频内的时间(发送UP主头像弹幕时无用，传入即可)
//
// platform 平台标识
// 1:web端
// 2:安卓端
// 8:视频管理页面
//
// data JSON序列 具体请看 https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/danmaku/action.md#%E5%8F%91%E9%80%81%E4%BA%92%E5%8A%A8%E5%BC%B9%E5%B9%95
//
// dmid 修改互动弹幕的弹幕id 不需要传入0即可 注意:修改弹幕platform必须为8
func (b *BiliClient) DanmakuCommandPost(tp int, aid int64, cid int64, progress int64, platform int, data string, dmid uint64) (*DanmakuCommandPostResult, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/v2/dm/command/post",
		"POST",
		map[string]string{
			"type":     strconv.Itoa(tp),
			"cid":      strconv.FormatInt(cid, 10),
			"aid":      strconv.FormatInt(aid, 10),
			"progress": strconv.FormatInt(progress, 10),
			"plat":     strconv.Itoa(platform),
			"data":     data,
			"dmid":     util.IF(dmid == 0, "", strconv.FormatUint(dmid, 10)).(string),
		},
	)
	if err != nil {
		return nil, err
	}

	var result *DanmakuCommandPostResult
	if err = json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// DanmakuSetConfig
//
// 弹幕个人配置修改
func (b *BiliClient) DanmakuSetConfig(conf *DanmakuConfig) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/v2/dm/web/config",
		"POST",
		map[string]string{
			"dm_switch":    strconv.FormatBool(conf.DmSwitch),
			"blockscroll":  strconv.FormatBool(conf.BlockScroll),
			"blocktop":     strconv.FormatBool(conf.BlockTop),
			"blockbottom":  strconv.FormatBool(conf.BlockBottom),
			"blockcolor":   strconv.FormatBool(conf.BlockColor),
			"blockspecial": strconv.FormatBool(conf.BlockSpecial),
			"ai_switch":    strconv.FormatBool(conf.AISwitch),
			"ai_level":     strconv.Itoa(conf.AILevel),
			"preventshade": strconv.FormatBool(conf.PreventShade),
			"dmask":        strconv.FormatBool(conf.DmMask),
			"opacity":      fmt.Sprintf("%.1f", conf.Opacity),
			"dmarea":       strconv.Itoa(conf.DmArea),
			"speedplus":    fmt.Sprintf("%.1f", conf.SpeedPlus),
			"fontsize":     fmt.Sprintf("%.1f", conf.FontSize),
			"screensync":   strconv.FormatBool(conf.ScreenSync),
			"speedsync":    strconv.FormatBool(conf.SpeedSync),
			"fontfamily":   conf.FontFamily,
			"bold":         strconv.FormatBool(conf.Bold),
			"fontborder":   strconv.Itoa(conf.FontBorder),
			"drawType":     conf.DrawType,
			"ts":           strconv.FormatInt(util.GetCST8Time(time.Now()).Unix(), 10),
		},
	)
	return err
}

// EmotePackGetMy 获取我的表情包列表
//
// business 使用场景	reply：评论区 dynamic：动态
func (b *BiliClient) EmotePackGetMy(business string) ([]*EmotePack, error) {
	resp, err := b.RawParse(
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

// EmotePackGetAll 获取全部表情包
//
// business 使用场景	reply：评论区 dynamic：动态
//
// B站接口导致必须登录才能获取
func (b *BiliClient) EmotePackGetAll(business string) ([]*EmotePack, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/emote/setting/panel",
		"GET",
		map[string]string{
			"business": business,
		},
	)
	if err != nil {
		return nil, err
	}
	var packs struct {
		AllPackages []*EmotePack `json:"all_packages,omitempty"`
	}
	if err = json.Unmarshal(resp.Data, &packs); err != nil {
		return nil, err
	}
	return packs.AllPackages, nil
}

// EmotePackAdd 添加表情包
//
// 只能添加有会员权限或已购买的表情包
//
// id 表情包id
//
// business 使用场景	reply：评论区 dynamic：动态
func (b *BiliClient) EmotePackAdd(id int64, business string) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/emote/package/add",
		"POST",
		map[string]string{
			"package_id": strconv.FormatInt(id, 10),
			"business":   business,
		},
	)
	return err
}

// EmotePackRemove 移除表情包
//
// id 表情包id
//
// business 使用场景	reply：评论区 dynamic：动态
func (b *BiliClient) EmotePackRemove(id int64, business string) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/emote/package/remove",
		"POST",
		map[string]string{
			"package_id": strconv.FormatInt(id, 10),
			"business":   business,
		},
	)
	return err
}

// AudioGetInfo
//
// 获取音频信息
func (b *BiliClient) AudioGetInfo(auid int64) (*AudioInfo, error) {
	resp, err := b.RawParse(
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

// AudioGetMyFavLists 获取自己创建的歌单
//
// pn 页码
//
// ps 每页项数
func (b *BiliClient) AudioGetMyFavLists(pn int, ps int) (*AudioMyFavLists, error) {
	resp, err := b.RawParse(
		BiliMainURL,
		"audio/music-service-c/web/collections/list",
		"GET",
		map[string]string{
			"pn": strconv.Itoa(pn),
			"ps": strconv.Itoa(ps),
		},
	)
	if err != nil {
		return nil, err
	}
	var coll *AudioMyFavLists
	if err = json.Unmarshal(resp.Data, &coll); err != nil {
		return nil, err
	}
	return coll, nil
}

// AudioGetPlayURL 获取音频流URL
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
// 取流：https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/audio/musicstream_url.md#%E9%9F%B3%E9%A2%91%E6%B5%81%E7%9A%84%E8%8E%B7%E5%8F%96
func (b *BiliClient) AudioGetPlayURL(auid int64, qn int) (*AudioPlayURL, error) {
	resp, err := b.RawParse(
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

// AudioIsFavored
//
// 查询音频是否被收藏
func (b *BiliClient) AudioIsFavored(auid int64) (bool, error) {
	resp, err := b.RawParse(
		BiliMainURL,
		"audio/music-service-c/web/collections/songs-coll",
		"GET",
		map[string]string{
			"sid": strconv.FormatInt(auid, 10),
		},
	)
	if err != nil {
		return false, err
	}
	var is bool
	if err = json.Unmarshal(resp.Data, &is); err != nil {
		return false, err
	}
	return is, nil
}

// AudioIsCoined 获取音频是否被投币
//
// 返回投币数
func (b *BiliClient) AudioIsCoined(auid int64) (int, error) {
	resp, err := b.RawParse(
		BiliMainURL,
		"audio/music-service-c/web/coin/audio",
		"GET",
		map[string]string{
			"sid": strconv.FormatInt(auid, 10),
		},
	)
	if err != nil {
		return -1, err
	}
	var coin int
	if err = json.Unmarshal(resp.Data, &coin); err != nil {
		return -1, err
	}
	return coin, nil
}

// ChargeTradeCreateBp 充电
//
// num B币数量 必须在2-9999之间
//
// mid 充电对象用户mid
//
// otype 充电来源	up：空间充电 archive：视频充电
//
// oid 充电来源代码 空间充电：充电对象用户mid 视频充电：稿件aid
func (b *BiliClient) ChargeTradeCreateBp(num int, mid int64, otype string, oid int64) (*ChargeBpResult, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/ugcpay/web/v2/trade/elec/pay/quick",
		"POST",
		map[string]string{
			"bp_num":              strconv.Itoa(num),
			"is_bp_remains_prior": "true",
			"up_mid":              strconv.FormatInt(mid, 10),
			"otype":               otype,
			"oid":                 strconv.FormatInt(oid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var result *ChargeBpResult
	if err = json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ChargeSetMessage 发送充电留言
//
// order 订单号，从充电成功的响应中获取
func (b *BiliClient) ChargeSetMessage(order string, message string) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/ugcpay/trade/elec/message",
		"POST",
		map[string]string{
			"order_id": order,
			"message":  message,
		},
	)
	return err
}

// ChargeTradeCreateQrCode 第三方如支付宝、微信充电
//
// num B币数量 必须在2-9999之间，1-19区间视为充值B币（未测试）
//
// prior 是否优先扣除B币余额
//
// mid 充电对象用户mid
//
// otype 充电来源	up：空间充电 archive：视频充电
//
// oid 充电来源代码 空间充电：充电对象用户mid 视频充电：稿件aid
//
// 整个支付流程请看：https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/electric/WeChat&Alipay.md
func (b *BiliClient) ChargeTradeCreateQrCode(num int, prior bool, mid int64, otype string, oid int64) (*ChargeCreateQrCode, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/ugcpay/web/v2/trade/elec/pay/qr_code/create",
		"POST",
		map[string]string{
			"bp_num":              strconv.Itoa(num),
			"up_mid":              strconv.FormatInt(mid, 10),
			"is_bp_remains_prior": strconv.FormatBool(prior),
			"otype":               otype,
			"oid":                 strconv.FormatInt(oid, 10),
		},
	)
	if err != nil {
		return nil, err
	}
	var qr *ChargeCreateQrCode
	if err = json.Unmarshal(resp.Data, &qr); err != nil {
		return nil, err
	}
	return qr, nil
}

// ChargeTradeCheckQrCode 获取第三方充电支付状态
//
// token ChargeTradeCreateQrCode 中返回的token
func (b *BiliClient) ChargeTradeCheckQrCode(token string) (*ChargeQrCodeStatus, error) {
	resp, err := b.RawParse(
		BiliApiURL,
		"x/ugcpay/trade/elec/pay/order/status",
		"GET",
		map[string]string{
			"qr_token": token,
		},
	)
	if err != nil {
		return nil, err
	}
	var status *ChargeQrCodeStatus
	if err = json.Unmarshal(resp.Data, &status); err != nil {
		return nil, err
	}
	return status, nil
}

// FollowUser
//
// 关注用户或取消关注
//
// follow true:关注 false:取消关注
func (b *BiliClient) FollowUser(mid int64, follow bool) error {
	_, err := b.RawParse(
		BiliApiURL,
		"x/relation/modify",
		"POST",
		map[string]string{
			"fid": strconv.FormatInt(mid, 10),
			"act": util.IF(follow, "1", "2").(string),
			// 以下为未知用途参数
			"re_src":         "11",
			"spmid":          "333.999.0.0",
			"extend_content": fmt.Sprintf(`{"entity":"user","entity_id":%s}`, strconv.FormatInt(mid, 10)),
		},
	)
	return err
}

// DynaCreatePlain 创建普通动态,返回创建的动态ID
//
// 具体请看测试样例
//
// 支持表情、at、话题
//
// 表情请使用 EmotePackGetMy 获取
//
// at请使用 [@刘庸干净又卫生 ] 注意末尾空格
//
// at map 里提供本次动态at的用户名与mid的映射
//
// 话题直接用 #xxx# 包裹即可
func (b *BiliClient) DynaCreatePlain(content string, at map[string]int64) (int64, error) {
	var ids []int64
	for _, id := range at {
		ids = append(ids, id)
	}

	ctrl, err := json.Marshal(parseDynaAt(1, content, at))
	if err != nil {
		return -1, err
	}
	resp, err := b.RawParse(
		BiliVcURL,
		"dynamic_svr/v1/dynamic_svr/create",
		"POST",
		map[string]string{
			"dynamic_id": "0",
			"type":       "4",
			"rid":        "0",
			"content":    content,
			"at_uids":    util.Int64SliceToString(ids, ","),
			"ctrl":       string(ctrl),
		},
	)
	if err != nil {
		return -1, err
	}

	var D struct {
		DynamicID int64 `json:"dynamic_id"`
	}
	if err = json.Unmarshal(resp.Data, &D); err != nil {
		return -1, err
	}
	return D.DynamicID, nil
}

// DynaLike 点赞动态
//
// like true:点赞 false: 不点赞
func (b *BiliClient) DynaLike(dyid int64, like bool) error {
	_, err := b.RawParse(
		BiliVcURL,
		"dynamic_like/v1/dynamic_like/thumb",
		"POST",
		map[string]string{
			"dynamic_id": strconv.FormatInt(dyid, 10),
			"up":         util.IF(like, "1", "2").(string),
		},
	)
	return err
}

// DynaUploadPics 上传动态图片
//
// 接口一次只能传一张，该函数为循环上传，如有速度需求请使用并发实现
//
// 返回的结构体用于创建
func (b *BiliClient) DynaUploadPics(pics []io.Reader) ([]*DynaUploadPic, error) {
	var results []*DynaUploadPic
	for _, p := range pics {
		// 该接口一次只能传一张，循环发送
		resp, err := b.UploadParse(
			BiliApiURL,
			"x/dynamic/feed/draw/upload_bfs",
			map[string]string{
				"biz":      "dyn",
				"category": "daily",
			},
			[]*FileUpload{{
				Field: "file_up",
				Name:  "1.jpg", // B站通过文件头判断content-type，该字段无用
				File:  p,
			}},
		)
		if err != nil {
			return nil, err
		}
		var r = &DynaUploadPic{}
		if err = json.Unmarshal(resp.Data, &r); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

// DynaCreateDraw 创建图片动态
//
// content,at 同 DynaCreatePlain
//
// pics 从 DynaUploadPics 获取
func (b *BiliClient) DynaCreateDraw(content string, at map[string]int64, pic []*DynaUploadPic) (int64, error) {
	var ids []int64
	for _, id := range at {
		ids = append(ids, id)
	}

	ctrl, err := json.Marshal(parseDynaAt(1, content, at))
	if err != nil {
		return -1, err
	}

	pj, err := genDynaPic(pic)
	if err != nil {
		return -1, err
	}
	resp, err := b.RawParse(
		BiliVcURL,
		"dynamic_svr/v1/dynamic_svr/create_draw",
		"POST",
		map[string]string{
			"biz":        "3",
			"category":   "3",
			"type":       "0",
			"pictures":   pj,
			"content":    content,
			"at_uids":    util.Int64SliceToString(ids, ","),
			"at_control": string(ctrl),
		},
	)
	if err != nil {
		return -1, err
	}

	var D struct {
		DynamicID int64 `json:"dynamic_id"`
	}
	if err = json.Unmarshal(resp.Data, &D); err != nil {
		return -1, err
	}
	return D.DynamicID, nil
}

// DynaRepost 转发动态
//
// dyid 为转发的动态ID
func (b *BiliClient) DynaRepost(dyid int64, content string, at map[string]int64) error {
	var ids []int64
	for _, id := range at {
		ids = append(ids, id)
	}

	ctrl, err := json.Marshal(parseDynaAt(1, content, at))
	if err != nil {
		return err
	}
	_, err = b.RawParse(
		BiliVcURL,
		"dynamic_repost/v1/dynamic_repost/repost",
		"POST",
		map[string]string{
			"dynamic_id": strconv.FormatInt(dyid, 10),
			"content":    content,
			"at_uids":    util.Int64SliceToString(ids, ","),
			"ctrl":       string(ctrl),
		},
	)
	return err
}

// DynaDel
//
// 删除动态
func (b *BiliClient) DynaDel(dyid int64) error {
	_, err := b.RawParse(
		BiliVcURL,
		"dynamic_svr/v1/dynamic_svr/rm_dynamic",
		"POST",
		map[string]string{
			"dynamic_id": strconv.FormatInt(dyid, 10),
		},
	)
	return err
}

// DynaCreateDraft 创建定时发布动态
//
// 返回draft id
//
// content,at 同 DynaCreatePlain
//
// pics 从 DynaUploadPics 获取
//
// publish 为指定发布的时间戳,换算后为东八区
func (b *BiliClient) DynaCreateDraft(content string, at map[string]int64, pic []*DynaUploadPic, publish int64) (int64, error) {
	var ids []int64
	for _, id := range at {
		ids = append(ids, id)
	}

	ctrl, err := json.Marshal(parseDynaAt(1, content, at))
	if err != nil {
		return -1, err
	}

	pj, err := genDynaPic(pic)
	if err != nil {
		return -1, err
	}

	request, err := json.Marshal(&dynaDraft{
		Biz:         3,
		Category:    3,
		Type:        0,
		Pictures:    pj,
		Description: content,
		Content:     content,
		From:        "create.dynamic.web",
		AtUIDs:      util.Int64SliceToString(ids, ","),
		AtControl:   string(ctrl),
	})
	if err != nil {
		return -1, err
	}
	resp, err := b.RawParse(
		BiliVcURL,
		"dynamic_draft/v1/dynamic_draft/add_draft",
		"POST",
		map[string]string{
			"type":         "4",
			"publish_time": strconv.FormatInt(publish, 10),
			"request":      string(request),
		},
	)
	if err != nil {
		return -1, err
	}

	var D struct {
		DraftID int64 `json:"draft_id"`
	}
	if err = json.Unmarshal(resp.Data, &D); err != nil {
		return -1, err
	}
	return D.DraftID, nil
}

// DynaModifyDraft 修改定时发布动态
//
// dfid 定时发布ID
//
// 其他参数同 DynaCreateDraft
func (b *BiliClient) DynaModifyDraft(dfid int64, content string, at map[string]int64, pic []*DynaUploadPic, publish int64) error {
	var ids []int64
	for _, id := range at {
		ids = append(ids, id)
	}

	ctrl, err := json.Marshal(parseDynaAt(1, content, at))
	if err != nil {
		return err
	}

	pj, err := genDynaPic(pic)
	if err != nil {
		return err
	}

	request, err := json.Marshal(&dynaDraft{
		Biz:         3,
		Category:    3,
		Type:        0,
		Pictures:    pj,
		Description: content,
		Content:     content,
		From:        "create.dynamic.web",
		AtUIDs:      util.Int64SliceToString(ids, ","),
		AtControl:   string(ctrl),
	})
	if err != nil {
		return err
	}

	_, err = b.RawParse(
		BiliVcURL,
		"dynamic_draft/v1/dynamic_draft/modify_draft",
		"POST",
		map[string]string{
			"draft_id":     strconv.FormatInt(dfid, 10),
			"type":         "2",
			"publish_time": strconv.FormatInt(publish, 10),
			"request":      string(request),
		},
	)
	return err
}

// DynaDelDraft 删除定时发布动态
//
// dfid 定时发布ID
func (b *BiliClient) DynaDelDraft(dfid int64) error {
	_, err := b.RawParse(
		BiliVcURL,
		"dynamic_draft/v1/dynamic_draft/rm_draft",
		"POST",
		map[string]string{
			"draft_id": strconv.FormatInt(dfid, 10),
		},
	)
	return err
}

// DynaPublishDraft 立即发布定时动态
//
// 返回发布的动态ID
//
// dfid 定时发布ID
func (b *BiliClient) DynaPublishDraft(dfid int64) (int64, error) {
	resp, err := b.RawParse(
		BiliVcURL,
		"dynamic_draft/v1/dynamic_draft/publish_now",
		"POST",
		map[string]string{
			"draft_id": strconv.FormatInt(dfid, 10),
		},
	)
	if err != nil {
		return -1, err
	}
	var r struct {
		DynamicID int64 `json:"dynamic_id"`
		CreateEc  int   `json:"create_ec"`
	}
	if err = json.Unmarshal(resp.Data, &r); err != nil {
		return -1, err
	}
	// 一些特殊情况导致发布失败，还有一层错误需要判断
	if r.CreateEc != 0 {
		return -1, fmt.Errorf("(%d) publish error", r.CreateEc)
	}
	return r.DynamicID, nil
}

// DynaGetDrafts
//
// 获取所有定时发布动态
func (b *BiliClient) DynaGetDrafts() (*DynaGetDraft, error) {
	resp, err := b.RawParse(
		BiliVcURL,
		"dynamic_draft/v1/dynamic_draft/get_drafts",
		"GET",
		map[string]string{},
	)
	if err != nil {
		return nil, err
	}
	var r = &DynaGetDraft{}
	if err = json.Unmarshal(resp.Data, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// LiveSendDanmaku 发送弹幕
//
// roomID: 真实直播间ID
//
// color: 颜色十进制，有权限控制.默认白色:16777215
//
// fontsize: 默认25
//
// mode: 1:飞行 5:顶部 4:底部
//
// msg: 弹幕内容
//
// bubble: 气泡弹幕?默认0
func (b *BiliClient) LiveSendDanmaku(roomID int64, color int64, fontsize int, mode int, msg string, bubble int) error {
	_, err := b.RawParse(
		BiliLiveURL,
		"msg/send",
		"POST",
		map[string]string{
			"roomid":   strconv.FormatInt(roomID, 10),
			"color":    strconv.FormatInt(color, 10),
			"fontsize": strconv.Itoa(fontsize),
			"mode":     strconv.Itoa(mode),
			"msg":      msg,
			"bubble":   strconv.Itoa(bubble),
			"rnd":      strconv.FormatInt(time.Now().Unix(), 10),
		},
	)
	return err
}
