// /qqmusic/models.go
package qqmusic

// SearchTypeEnum 模拟 C# 的枚举
type SearchType int

// 0单曲 2专辑 1歌手 3歌单 7歌词 12mv
const (
	SearchTypeSong     SearchType = 0
	SearchTypeAlbum    SearchType = 2
	SearchTypePlaylist SearchType = 3
)

// MusicFcgApiResult 对应 QQMusicBean.MusicFcgApiResult
type MusicFcgApiResult struct {
	Code int `json:"code"`
	// Raw  string `json:"raw"`
	Req struct {
		Code int `json:"code"`
		Data struct {
			Sip []string `json:"sip"`
		} `json:"data"`
	} `json:"req"`
	Req0 struct {
		Code int `json:"code"`
		Data struct {
			MidURLInfo []struct {
				PURL string `json:"purl"`
			} `json:"midurlinfo"`
		} `json:"data"`
	} `json:"req_0"`
	Req1 struct {
		Code int `json:"code"`
		Data struct {
			Body struct {
				Song struct {
					List []SongInfo `json:"list"`
				} `json:"song"`
			} `json:"body"`
		} `json:"data"`
	} `json:"req_1"`
}

// AlbumResult 对应 QQMusicBean.AlbumResult
type AlbumResult struct {
	Code int   `json:"code"`
	Data Album `json:"data"`
}

type Album struct {
	List []SongInfo `json:"list"`
}

// PlaylistResult 对应 QQMusicBean.PlaylistResult
type PlaylistResult struct {
	Code   int `json:"code"`
	Cdlist []struct {
		Songlist []SongInfo `json:"songlist"`
	} `json:"cdlist"`
}

// SongResult 对应 QQMusicBean.SongResult
type SongResult struct {
	Code int        `json:"code"`
	Data []SongInfo `json:"data"`
}

// SongInfo 包含歌曲的基本信息
type SongInfo struct {
	ID         int    `json:"id"`
	Mid        string `json:"mid"`
	Name       string `json:"name"`
	Interval   int    `json:"interval"`
	SingerName string `json:"singerName"`
	Singer     []struct {
		Name string `json:"name"`
	} `json:"singer"`
}

// LyricResult 对应 QQMusicBean.LyricResult
type LyricResult struct {
	Code  int    `json:"code"`
	Lyric string `json:"lyric"`
	Trans string `json:"trans"`
	Roma  string `json:"roma"`
}

// ResultVo 简单的泛型结果包装
type ResultVo[T any] struct {
	Data T
}
