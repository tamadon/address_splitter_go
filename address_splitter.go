package main

import "fmt"
import "regexp"

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

func main() {
	address := "〒105-0011 東京都港区芝公園4丁目2-8東京タワー大展望台2F TEL:03-3433-5111"
	trimedAddress := trimExtraString(address)
	prefecture := getPrefecture(trimedAddress)
	city := getCity(trimedAddress)

	result := prefecture + city

	fmt.Println(result)
}