package service

import (
	"dlrc/utils"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery" // 引入 goquery 库
)

// SongInfo 结构体用于存储歌曲的 ID 和 Name
type SongInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var Cookie string

func HandleCookie(w http.ResponseWriter, r *http.Request) {
	utils.Cors(w, r)
	Cookie = utils.GetParam(r, "cookie", "c")
	fmt.Println("set Cookie:", Cookie)
	utils.Ok(w, Cookie)
}

func HandleSongs(w http.ResponseWriter, r *http.Request) {
	utils.Cors(w, r)
	name := utils.GetParam(r, "name", "n")
	if name == "" {
		utils.FailMsg(w, "请传入歌曲名称参数: name")
		return
	}

	songs, err := GetTop10SongsFromURL("http://www.2t58.com/so/"+name+".html", utils.GetParamInt(r, "limit", 20))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.Ok(w, songs)
}

func HandleLyric(w http.ResponseWriter, r *http.Request) {
	utils.Cors(w, r)

	id := utils.GetParam(r, "id")
	if id == "" {
		utils.FailMsg(w, "请传入歌曲ID参数: id")
		return
	}

	req, err := http.NewRequest("GET", "http://www.2t58.com/js/lrc.php?cid="+id, nil) // 第三个参数是请求体，GET 请求为 nil
	if err != nil {
		http.Error(w, fmt.Sprintf("后端请求失败: %v", err), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Host", "www.2t58.com")
	req.Header.Set("Referer", "http://www.2t58.com/song/"+id+".html")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.5845.97 Safari/537.36 Core/1.116.508.400 QQBrowser/19.1.6429.400")
	if Cookie != "" {
		req.Header.Set("Cookie", Cookie)
	}
	client := &http.Client{}
	resp, err := client.Do(req)

	// resp, err := http.Get("http://www.2t58.com/js/lrc.php?cid=" + id)

	if err != nil {
		fmt.Println("后端请求失败:", err)
		http.Error(w, fmt.Sprintf("后端请求失败: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("后端请求失败:", err)
		http.Error(w, fmt.Sprintf("后端请求失败: %v", err), http.StatusInternalServerError)
		return
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	} else {
		w.Header().Set("Content-Type", "application/json")
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		fmt.Printf("写入前端响应时出错: %v\n", err)
	}
}

// GetTop10SongsFromURL 函数发送 HTTP GET 请求到指定的 URL，
// 解析 HTML 并提取前 10 个歌曲的 ID 和 Name。
func GetTop10SongsFromURL(url string, max int) ([]SongInfo, error) {
	req, err := http.NewRequest("GET", url, nil) // 第三个参数是请求体，GET 请求为 nil
	if err != nil {
		return nil, fmt.Errorf("发送 HTTP 请求失败: %w", err)
	}
	req.Header.Set("Host", "www.2t58.com")
	req.Header.Set("Referer", url)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.5845.97 Safari/537.36 Core/1.116.508.400 QQBrowser/19.1.6429.400")
	if Cookie != "" {
		req.Header.Set("Cookie", Cookie)
	}
	client := &http.Client{}
	resp, err := client.Do(req)

	// resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("发送 HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("后端请求失败:", err)
		return nil, fmt.Errorf("HTTP 请求返回非 200 状态码: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}
	htmlContent := string(bodyBytes)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("解析 HTML 失败: %w", err)
	}

	// 4. 提取 `.play_list ul li .name a` 元素
	// CSS 选择器指定了你要查找的元素
	songElements := doc.Find(".play_list ul li .name a")

	var songs []SongInfo
	count := 0

	songElements.Each(func(i int, s *goquery.Selection) {
		if count >= max {
			return
		}

		href, exists := s.Attr("href")
		if !exists {
			return // 如果没有 href 属性，跳过此元素
		}

		startIndex := strings.Index(href, "song/")
		endIndex := strings.Index(href, ".htm")
		if startIndex == -1 || endIndex == -1 || startIndex+6 >= endIndex {
			return // 如果格式不匹配，跳过此元素
		}

		songs = append(songs, SongInfo{
			ID:   href[startIndex+5 : endIndex],
			Name: s.Text(),
		})
		count++
	})

	return songs, nil
}

// func main() {
// 	url := "http://www.222.com/so/22.html" // 替换成你实际要请求的URL
// 	songs, err := GetTop10SongsFromURL(url)
// 	if err != nil {
// 		fmt.Println("错误:", err)
// 		return
// 	}

// 	if len(songs) == 0 {
// 		fmt.Println("未找到任何歌曲信息。")
// 		return
// 	}

// 	fmt.Println("成功获取前10首歌曲信息:")
// 	for i, song := range songs {
// 		fmt.Printf("%d. ID: %s, Name: %s\n", i+1, song.ID, song.Name)
// 	}
// }
