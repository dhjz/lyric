// /qqmusic/qq_music_native_api.go
package qqmusic

import (
	"dlrc/service/base"
	"dlrc/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var verbatimXmlMappingDict = map[string]string{
	"content":     "orig",  // 原文
	"contentts":   "ts",    // 译文
	"contentroma": "roma",  // 罗马音
	"Lyric_1":     "lyric", // 解压后的内容
}

// QQMusicNativeApi 结构体，用于实现 QQ 音乐的 API 调用
type QQMusicNativeApi struct {
	*BaseNativeApi
}

//	q := NewQQMusicNativeApi(func() string {
//		return utils.GetParam(r, "cookie")
//	})
var QQApi = NewQQMusicNativeApi(nil)

func HandleSongs(w http.ResponseWriter, r *http.Request) {
	utils.Cors(w, r)
	// 初始化实例QQMusicNativeApi
	result, err := QQApi.Search(utils.GetParam(r, "name"), SearchTypeSong)
	if err != nil {
		utils.Fail(w)
		return
	}

	var targets []base.BaseSong
	if len(result.Req1.Data.Body.Song.List) > 0 {
		for _, song := range result.Req1.Data.Body.Song.List {
			singerName := ""
			if len(song.Singer) > 0 {
				singerName = song.Singer[0].Name
			}
			targets = append(targets, base.BaseSong{
				ID:       fmt.Sprintf("%d", song.ID),
				Mid:      song.Mid,
				Name:     song.Name,
				Singer:   singerName,
				Interval: song.Interval,
			})
		}
	}

	utils.Ok(w, targets)
}

func HandleLyric(w http.ResponseWriter, r *http.Request) {
	utils.Cors(w, r)
	result, err := QQApi.GetLyric(utils.GetParam(r, "id"))
	if err != nil {
		utils.Fail(w)
		return
	}
	result.Lyric = FormatLyric(result.Lyric)
	result.Trans = FormatLyric(result.Trans)
	result.Roma = FormatLyric(result.Roma)
	utils.Ok(w, result)
}

// NewQQMusicNativeApi 创建一个新的 QQMusicNativeApi 实例
func NewQQMusicNativeApi(cookieFunc func() string) *QQMusicNativeApi {
	return &QQMusicNativeApi{
		BaseNativeApi: NewBaseNativeApi(cookieFunc),
	}
}

// HttpRefer "重写" 基类的方法
func (q *QQMusicNativeApi) HttpRefer() string {
	return "https://c.y.qq.com/"
}

// Search 搜索歌曲、专辑或播放列表
func (q *QQMusicNativeApi) Search(keyword string, searchType SearchType) (*MusicFcgApiResult, error) {
	// 0单曲 2专辑 1歌手 3歌单 7歌词 12mv
	payload := map[string]interface{}{
		"req_1": map[string]interface{}{
			"method": "DoSearchForQQMusicDesktop",
			"module": "music.search.SearchCgiService",
			"param": map[string]interface{}{
				"num_per_page": 20,
				"page_num":     1,
				"query":        keyword,
				"search_type":  searchType,
			},
		},
	}

	respStr, err := q.sendJsonPost("https://u.y.qq.com/cgi-bin/musicu.fcg", payload)
	if err != nil {
		return nil, err
	}

	var result MusicFcgApiResult
	if err := json.Unmarshal([]byte(respStr), &result); err != nil {
		return nil, err
	}
	// result.Raw = respStr
	return &result, nil
}

// GetAlbum 获取专辑信息
func (q *QQMusicNativeApi) GetAlbum(albumID string) (*AlbumResult, error) {
	key := "albummid"
	if isNumeric(albumID) {
		key = "albumid"
	}
	data := map[string]string{key: albumID}

	respStr, err := q.sendPost("https://c.y.qq.com/v8/fcg-bin/fcg_v8_album_info_cp.fcg", data)
	if err != nil {
		return nil, err
	}

	var result AlbumResult
	if err := json.Unmarshal([]byte(respStr), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetPlaylist 获取播放列表信息
func (q *QQMusicNativeApi) GetPlaylist(playlistID string) (*PlaylistResult, error) {
	data := map[string]string{
		"disstid":    playlistID,
		"format":     "json",
		"outCharset": "utf8",
		"type":       "1",
		"json":       "1",
		"utf8":       "1",
		"onlysong":   "0",
		"new_format": "1",
	}

	respStr, err := q.sendPost("https://c.y.qq.com/qzone/fcg-bin/fcg_ucc_getcdinfo_byids_cp.fcg", data)
	if err != nil {
		return nil, err
	}

	var result PlaylistResult
	if err := json.Unmarshal([]byte(respStr), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSong 获取单曲信息
func (q *QQMusicNativeApi) GetSong(id string) (*SongResult, error) {
	const callBack = "getOneSongInfoCallback"
	key := "songmid"
	if isNumeric(id) {
		key = "songid"
	}

	data := map[string]string{
		key:             id,
		"tpl":           "yqq_song_detail",
		"format":        "jsonp",
		"callback":      callBack,
		"g_tk":          "5381",
		"jsonpCallback": callBack,
		"loginUin":      "0",
		"hostUin":       "0",
		"outCharset":    "utf8",
		"notice":        "0",
		"platform":      "yqq",
		"needNewCode":   "0",
	}

	respStr, err := q.sendPost("https://c.y.qq.com/v8/fcg-bin/fcg_play_single_song.fcg", data)
	if err != nil {
		return nil, err
	}

	jsonStr := resolveRespJson(callBack, respStr)
	var result SongResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetLyric 获取歌词
func (q *QQMusicNativeApi) GetLyric(songID string) (*LyricResult, error) {
	data := map[string]string{
		"version":     "15",
		"miniversion": "82",
		"lrctype":     "4",
		"musicid":     songID,
	}

	respStr, err := q.sendPost("https://c.y.qq.com/qqmusic/fcgi-bin/lyric_download.fcg", data)
	if err != nil {
		return nil, err
	}
	// log.Printf("GetLyric  songID %s: %v", songID, respStr)

	respStr = strings.ReplaceAll(respStr, "<!--", "")
	respStr = strings.ReplaceAll(respStr, "-->", "")
	respStr = regexp.MustCompile(`<miniversion.*?/>`).ReplaceAllString(respStr, "")

	foundElements, err := recursionFindElement(respStr, verbatimXmlMappingDict)
	if err != nil {
		return nil, fmt.Errorf("failed to parse lyric xml: %w", err)
	}
	// log.Printf("GetLyric  songID %s: %v", songID, foundElements)

	result := &LyricResult{Code: 0}

	for key, text := range foundElements {
		if strings.TrimSpace(text) == "" {
			continue
		}

		decompressText, err := DecryptLyrics(text)
		if err != nil {
			log.Printf("GetLyric DecryptLyrics failed for songID %s: %v", songID, err)
			continue
		}

		lyricContent, err := parseLyricContentXML(decompressText)
		if err != nil {
			log.Printf("Failed to parse inner lyric XML for songID %s: %v", songID, err)
			lyricContent = decompressText // Fallback to raw decompressed text
		}

		switch key {
		case "orig":
			result.Lyric = lyricContent
		case "ts":
			result.Trans = lyricContent
		case "roma":
			result.Roma = lyricContent
		}
	}

	return result, nil
}

// GetSongLink 获取歌曲播放链接
func (q *QQMusicNativeApi) GetSongLink(songMid string) (*ResultVo[string], error) {
	guid := GetGuid()

	payload := map[string]interface{}{
		"req": map[string]interface{}{
			"method": "GetCdnDispatch",
			"module": "CDN.SrfCdnDispatchServer",
			"param": map[string]interface{}{
				"guid":     guid,
				"calltype": "0",
				"userip":   "",
			},
		},
		"req_0": map[string]interface{}{
			"method": "CgiGetVkey",
			"module": "vkey.GetVkeyServer",
			"param": map[string]interface{}{
				"guid":      "8348972662",
				"songmid":   []string{songMid},
				"songtype":  []int{1},
				"uin":       "0",
				"loginflag": 1,
				"platform":  "20",
			},
		},
		"comm": map[string]interface{}{
			"uin":    0,
			"format": "json",
			"ct":     24,
			"cv":     0,
		},
	}

	respStr, err := q.sendJsonPost("https://u.y.qq.com/cgi-bin/musicu.fcg", payload)
	if err != nil {
		return nil, err
	}

	var res MusicFcgApiResult
	if err := json.Unmarshal([]byte(respStr), &res); err != nil {
		return nil, err
	}

	link := ""
	if res.Code == 0 && res.Req.Code == 0 && res.Req0.Code == 0 &&
		len(res.Req.Data.Sip) > 0 && len(res.Req0.Data.MidURLInfo) > 0 {
		link = res.Req.Data.Sip[0] + res.Req0.Data.MidURLInfo[0].PURL
	}

	return &ResultVo[string]{Data: link}, nil
}
