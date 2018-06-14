package main

import "fmt"
import "regexp"
import "os"
import "bufio"
import "strings"

func trimExtraString(address string) string {
	// 電話番号と思わしき文字列を削除
	result := trimStringRegexp(address, `[\d\(\)-]{9,}`)
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

func trimStringRegexp(inputString string, regexpString string) string {
	rep := regexp.MustCompile(regexpString)
	return rep.ReplaceAllString(inputString, "")
}

func getPrefecture(inputString string) string {
	rep := regexp.MustCompile(`[^\x00-\x7F]{2,3}県|..府|東京都|北海道`)
	if rep.MatchString(inputString) {
		return rep.FindAllStringSubmatch(inputString , -1)[0][0]
	}
	return ""
}

func getCity(inputString string) string {
	regexPattern := []string{}

	regexPattern = append(regexPattern, `([^\x00-\x7F]{1,6}市[^\x00-\x7F]{1,4}区)`)
	regexPattern = append(regexPattern, `([^\x00-\x7F]{1,3}郡[^\x00-\x7F]{1,5}町)`)
	regexPattern = append(regexPattern, `(四日|廿日|野々)市市`)
	regexPattern = append(regexPattern, `([^\x00-\x7F市]{1,6}市)`)
	regexPattern = append(regexPattern, `([^\x00-\x7F]{1,4}区)`)

	for _, pattern := range regexPattern {
		rep := regexp.MustCompile(pattern)
		if rep.MatchString(inputString) {
			return rep.FindAllStringSubmatch(inputString , -1)[0][0]
		}
	}

	return ""
}

func getAddress(inputString string) string {
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
	if rep.MatchString(inputString) {
		return rep.FindAllStringSubmatch(inputString , -1)[0][0]
	}

	return ""
}

// 文字列を1行入力
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

	town := trimStringRegexp(trimedPrefectureCityStr, address)

	etc := trimStringRegexp(trimedPrefectureStr, city)

	fmt.Println(trimedStr)
	fmt.Println(prefecture)
	fmt.Println(city)
	fmt.Println(town)
	fmt.Println(address)
	fmt.Println(etc)
}
