package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UserData struct {
	user          string
	avatar        string
	species       string
	lastTime      string
	totalRun      int64
	totalJump     int64
	totalDuration int64
}

func handler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	name := values.Get("user")
	userData := formatUserData(fetchGithubData(name))
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0")

	fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' width='495' height='195' viewBox='0 0 495 195' fill='none'>  \n")
	fmt.Fprintf(w, "  <style>.header { font: 600 18px 'Segoe UI', Ubuntu, Sans-Serif; fill: #2f80ed; animation: fadeInAnimation 0.8s ease-in-out forwards; } @supports(-moz-appearance: auto) { /* Selector detects Firefox */ .header { font-size: 15.5px; } } .stat { font: 600 14px 'Segoe UI', Ubuntu, 'Helvetica Neue', Sans-Serif; fill: #434d58; } @supports(-moz-appearance: auto) { /* Selector detects Firefox */ .stat { font-size:12px; } } .stagger { opacity: 0; animation: fadeInAnimation 0.3s ease-in-out forwards; } .avatar { animation: scaleInAnimation 0.3s ease-in-out forwards; } .bold { font-weight: 700 } @keyframes rankAnimation { from { stroke-dashoffset: 251.32741228718345; } to { stroke-dashoffset: 10.556231094032176; } } /* Animations */ @keyframes scaleInAnimation { from { transform: translate(-5px, 5px) scale(0); } to { transform: translate(-5px, 5px) scale(1); } } @keyframes fadeInAnimation { from { opacity: 0; } to { opacity: 1; } }</style>  \n")
	fmt.Fprintf(w, "  <rect data-testid='card-bg' x='0.5' y='0.5' rx='4.5' height='190' stroke='#e4e2e2' width='440' fill='#fffefe' stroke-opacity='1'/>  \n")
	fmt.Fprintf(w, "  <g data-testid='card-title' transform='translate(25, 35)'> \n")
	fmt.Fprintf(w, "    <g transform='translate(0, 0)'> \n")
	fmt.Fprintf(w, "      <text x='0' y='0' class='header' data-testid='header'>%s's Stats</text> \n", userData.user)
	fmt.Fprintf(w, "    </g> \n")
	fmt.Fprintf(w, "  </g>  \n")
	fmt.Fprintf(w, "  <g data-testid='main-card-body' transform='translate(0, 55)'> \n")
	fmt.Fprintf(w, "    <g transform='translate(340, 10)'> \n")
	fmt.Fprintf(w, "      <g class='avatar'> \n")
	// fmt.Fprintf(w, "      	<image width='80' height='80' href='%s'></image> \n", userData.avatar)
	fmt.Fprintf(w, "      </g> \n")
	fmt.Fprintf(w, "    </g>  \n")
	fmt.Fprintf(w, "    <svg x='0' y='0'> \n")
	fmt.Fprintf(w, "      <g transform='translate(0, 0)'> \n")
	fmt.Fprintf(w, "        <g class='stagger' style='animation-delay: 450ms' transform='translate(25, 0)'> \n")
	fmt.Fprintf(w, "          <text class='stat bold' y='12.5'>Last Time:</text>  \n")
	fmt.Fprintf(w, "          <text class='stat' x='130' y='12.5' data-testid='stars'>%s</text> \n", userData.lastTime)
	fmt.Fprintf(w, "        </g> \n")
	fmt.Fprintf(w, "      </g>\n")
	fmt.Fprintf(w, "      <g transform='translate(0, 25)'> \n")
	fmt.Fprintf(w, "        <g class='stagger' style='animation-delay: 600ms' transform='translate(25, 0)'> \n")
	fmt.Fprintf(w, "          <text class='stat bold' y='12.5'>Last Record:</text>  \n")
	fmt.Fprintf(w, "          <text class='stat' x='130' y='12.5' data-testid='commits'>%s</text> \n", userData.species)
	fmt.Fprintf(w, "        </g> \n")
	fmt.Fprintf(w, "      </g>\n")
	fmt.Fprintf(w, "      <g transform='translate(0, 50)'> \n")
	fmt.Fprintf(w, "        <g class='stagger' style='animation-delay: 750ms' transform='translate(25, 0)'> \n")
	fmt.Fprintf(w, "          <text class='stat bold' y='12.5'>Total Run:</text>  \n")
	fmt.Fprintf(w, "          <text class='stat' x='130' y='12.5' data-testid='prs'>%dm</text> \n", userData.totalRun)
	fmt.Fprintf(w, "        </g> \n")
	fmt.Fprintf(w, "      </g>\n")
	fmt.Fprintf(w, "      <g transform='translate(0, 75)'> \n")
	fmt.Fprintf(w, "        <g class='stagger' style='animation-delay: 900ms' transform='translate(25, 0)'> \n")
	fmt.Fprintf(w, "          <text class='stat bold' y='12.5'>Total Jump:</text>  \n")
	fmt.Fprintf(w, "          <text class='stat' x='130' y='12.5' data-testid='issues'>%d</text> \n", userData.totalJump)
	fmt.Fprintf(w, "        </g> \n")
	fmt.Fprintf(w, "      </g>\n")
	fmt.Fprintf(w, "      <g transform='translate(0, 100)'> \n")
	fmt.Fprintf(w, "        <g class='stagger' style='animation-delay: 1050ms' transform='translate(25, 0)'> \n")
	fmt.Fprintf(w, "          <text class='stat bold' y='12.5'>Total Duration:</text>  \n")
	fmt.Fprintf(w, "          <text class='stat' x='130' y='12.5' data-testid='contribs'>%d min</text> \n", userData.totalDuration)
	fmt.Fprintf(w, "        </g> \n")
	fmt.Fprintf(w, "      </g> \n")
	fmt.Fprintf(w, "    </svg> \n")
	fmt.Fprintf(w, "  </g> \n")
	fmt.Fprintf(w, "</svg>\n")
}

func fetchGithubData(name string) []interface{} {
	url := fmt.Sprintf("https://api.github.com/search/issues?q=author:%s+repo:pengcu/sexy-programmer-online-sex", name)
	c := &http.Client{
		Timeout: 5 * time.Second,
	}
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
	response, err := c.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	var result map[string]interface{}
	json.NewDecoder(response.Body).Decode(&result)
	return result["items"].([]interface{})
}

func formatUserData(list []interface{}) (userData UserData) {
	var data []string
	if len(list) > 0 {
		lastItem := list[0].(map[string]interface{})
		userData.user = lastItem["user"].(map[string]interface{})["login"].(string)
		userData.avatar = lastItem["user"].(map[string]interface{})["avatar_url"].(string)
		userData.lastTime = lastItem["created_at"].(string)
		for _, val := range list {
			data = append(data, val.(map[string]interface{})["title"].(string))
		}
		userData.species = data[0]
		for _, val := range data {
			array := strings.Split(val, "，")
			switch array[0] {
			case "跳绳":
				jump, err := strconv.ParseInt(array[2], 10, 64)
				if err != nil {
					panic(err)
				}

				time, err := strconv.ParseInt(strings.Split(array[1], "分钟")[0], 10, 64)
				if err != nil {
					panic(err)
				}
				userData.totalJump += jump
				userData.totalDuration += time
			case "跑步":
				run, err := strconv.ParseInt(array[2], 10, 64)
				if err != nil {
					panic(err)
				}
				time, err := strconv.ParseInt(strings.Split(array[1], "分钟")[0], 10, 64)
				if err != nil {
					panic(err)
				}
				userData.totalRun += run
				userData.totalDuration += time
			default:
				fmt.Println("error type issue", array)
			}
		}
	}
	return userData
}

func main() {
	http.HandleFunc("/api", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
