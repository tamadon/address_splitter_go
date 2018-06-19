package main

import (
	"fmt"
	"regexp"
	"os"
	"bufio"
	"strings"
	"strconv"
	"github.com/kurehajime/cjk2num"
	"golang.org/x/text/unicode/norm"
)

// 条件にマッチした文字列をトリムする
func trimStringRegexp(input string, regexpString string) string {
	rep := regexp.MustCompile(regexpString)
	return rep.ReplaceAllString(input, "")
}


// 余計な文字列をトリムする
func trimExtraString(input string) string {
	// 電話番号と思わしき文字列を削除
	result := trimStringRegexp(input, `[\d\(\)-]{9,}`)
	result = trimStringRegexp(result, `TEL:|FAX:|TEL|FAX`)
	// 郵便番号と思わしき文字列を削除
	result = trimStringRegexp(result, `\d\d\d-\d\d\d\d`)
	result = trimStringRegexp(result, `〒|郵便番号|郵便`)
	// 括弧に囲われた部分を削除
	result = trimStringRegexp(result, `【.*?】`)
	result = trimStringRegexp(result, `≪.*?≫`)
	result = trimStringRegexp(result, `《.*?》`)
	result = trimStringRegexp(result, `◎.*?◎`)
	result = trimStringRegexp(result, `〔.*?〕`)
	result = trimStringRegexp(result, `\[.*?\]`)
	result = trimStringRegexp(result, `<.*?>`)
	result = trimStringRegexp(result, `\(.*?\)`)
	result = trimStringRegexp(result, `「.*?」`)
	// 特定フレーズの後にある文字を削除
	result = trimStringRegexp(result, `(◎|※|☆|★|◇|◆|□|■|●|○|~|〜).*`)
	result = trimStringRegexp(result, ` `)

	return result
}

// 都道府県名を取得する
func getPrefecture(input string) string {
	rep := regexp.MustCompile(`[^\x00-\x7F]{2,3}県|..府|東京都|北海道`)
	if rep.MatchString(input) {
		return rep.FindAllStringSubmatch(input , -1)[0][0]
	}
	return ""
}

// 市区名を取得する
func getCity(input string) string {
	regexPattern := []string{}

	regexPattern = append(regexPattern, `([^\x00-\x7F]{1,6}市[^\x00-\x7F]{1,4}区)`)
	regexPattern = append(regexPattern, `([^\x00-\x7F]{1,3}郡[^\x00-\x7F]{1,5}町)`)
	regexPattern = append(regexPattern, `(四日|廿日|野々)市市`)
	regexPattern = append(regexPattern, `([^\x00-\x7F市]{1,6}市)`)
	regexPattern = append(regexPattern, `([^\x00-\x7F]{1,4}区)`)

	for _, pattern := range regexPattern {
		rep := regexp.MustCompile(pattern)
		if rep.MatchString(input) {
			return rep.FindAllStringSubmatch(input , -1)[0][0]
		}
	}
	return ""
}

// 番地を取得する
func getAddress(input string) string {
	// 数字
	num := `[一二三四五六七八九十百千万]|[0-9]|[０-９]`
	// 繋ぎ文字1：数字と数字の間(末尾以外)
	// s_str1 := `(丁目|丁|番地|番|号|-|‐|ー|−|の|東|西|南|北)`
	s_str1 := `(丁目|丁|番地|番|号|-|‐|ー|−|の)`
	// 繋ぎ文字2：数字と数字の間(末尾)
	s_str2 := `(丁目|丁|番地|番|号)`
	// 全ての数字
	all_num := `(\\d+|` + num + `+)`

	pattern := all_num + `*(` + all_num + `|` + s_str1 + `{1,2})*(` + all_num + `|` + s_str2 + `{1,2})`
	rep := regexp.MustCompile(pattern)
	if rep.MatchString(input) {
		return rep.FindAllStringSubmatch(input , -1)[0][0]
	}
	return ""
}

// 番地を正規化する
func norm_address(input string) string {
    // ハイフン以外のハイフンっぽい記号を置き換える
		rep := regexp.MustCompile(`-|‐|ー|−`)
		result := rep.ReplaceAllString(input, "-")
    // 「丁目」などをハイフンに置き換える
		rep2 := regexp.MustCompile(`丁目|丁|番地|番|号|の`)
		result = rep2.ReplaceAllString(result, "-")
		rep3 := regexp.MustCompile(`-{2,}`)
		result = rep3.ReplaceAllString(result, "-")
		rep4 := regexp.MustCompile(`(^-)|(-$)`)
		result = rep4.ReplaceAllString(result, "")
    // 全角数字、漢数字を半角アラビア数字に置き換える
		halfNum := `[0-9]`
		halfNumRep := regexp.MustCompile(halfNum)
		fullNum := `[０-９]`
		fullNumRep := regexp.MustCompile(fullNum)

		var resultSlice []string
		arr := strings.Split(result, "-")

		for _, num := range arr {
			if halfNumRep.MatchString(num) { // 半角数字
				resultSlice = append(resultSlice, num)
			} else if fullNumRep.MatchString(num) { // 全角アラビア数字
				resultSlice = append(resultSlice, string(norm.NFKC.Bytes([]byte(num))))
			} else { // それ以外＝漢数字
				convertedNum, err := cjk2num.Convert(num)
			  if err != nil {
			    fmt.Println(err.Error())
			  }
				resultSlice = append(resultSlice, strconv.FormatInt(convertedNum, 10))
			}
		}
    return strings.Join(resultSlice, "-")
}

// 町村と建物名を取得する
func getTownAndBuilding(input string, splitter string) (string, string) {
	arr := strings.Split(input, splitter)
	town := arr[0]

	building := ""
	if len(arr) == 2 {
		building = arr[1]
	}
	return town, building
}

// 標準出力から文字列を1行取得する
func StrStdin() (stringInput string) {
    scanner := bufio.NewScanner(os.Stdin)

    scanner.Scan()
    stringInput = scanner.Text()
    stringInput = strings.TrimSpace(stringInput)
    return
}

func main() {
	inputStr := StrStdin()
	trimedStr := trimExtraString(inputStr)
	prefecture := getPrefecture(trimedStr)

	trimedPrefectureStr := trimStringRegexp(trimedStr, prefecture)
	city := getCity(trimedPrefectureStr)

	trimedPrefectureCityStr := trimStringRegexp(trimedPrefectureStr, city)
	address := getAddress(trimedPrefectureCityStr)
	normAddress := norm_address(address)

	town, building := getTownAndBuilding(trimedPrefectureCityStr, address)

	fmt.Println("\n")
	fmt.Println(prefecture)
	fmt.Println(city)
	fmt.Println(town)
	// fmt.Println(address)
	fmt.Println(normAddress)
	fmt.Println(building)
}
