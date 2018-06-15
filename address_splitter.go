package main

import "fmt"
import "regexp"
import "os"
import "bufio"
import "strings"

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

func trimStringRegexp(input string, regexpString string) string {
	rep := regexp.MustCompile(regexpString)
	return rep.ReplaceAllString(input, "")
}

func getPrefecture(input string) string {
	rep := regexp.MustCompile(`[^\x00-\x7F]{2,3}県|..府|東京都|北海道`)
	if rep.MatchString(input) {
		return rep.FindAllStringSubmatch(input , -1)[0][0]
	}
	return ""
}

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
func norm_addr1(input string) string {
    addr1_temp := input
    // ハイフン以外のハイフンっぽい記号を置き換える
		rep := regexp.MustCompile(`-|‐|ー|−`)
		hoge := rep.ReplaceAllString(addr1_temp, "-")
    // 「丁目」などをハイフンに置き換える
		rep2 := regexp.MustCompile(`丁目|丁|番地|番|号|の`)
		hoge = rep2.ReplaceAllString(hoge, "-")
		rep3 := regexp.MustCompile(`-{2,}`)
		hoge = rep3.ReplaceAllString(hoge, "-")
		rep4 := regexp.MustCompile(`(^-)|(-$)`)
		hoge = rep4.ReplaceAllString(hoge, "")
    // # 漢数字をアラビア数字に置き換える
    // pattern = /[一二三四五六七八九十百千万]+/
    // while addr1_temp =~ pattern
    //     match_string = addr1_temp.match(pattern)[0]
    //     arabia_number_string = "#{kan_to_arabia(match_string)}"
    //     addr1_temp.sub!(match_string, arabia_number_string)
    // end
    return hoge
}

// # 漢数字をアラビア数字に変換する
// # 実は「十一万」以上の文字列で変換ミスが発生するが、
// # 番地変換でそこまで大きな数を考慮することはないと思われる
// func kan_to_arabia(str)
//     // 変換するためのハッシュ
// 		m := map[string]int{
// 			"一": 1, "二": 2, "三": 3, "四": 4, "五": 5,
// 			"六": 6, "七": 7, "八": 8, "九": 9, "○": 0,
// 			"十": 10, "百": 100, "千": 1000, "万": 10000
// 		}
//     # 漢数字を数字に置き換える
//     num_array = str.chars.to_a.map{|c| hash[c]}
//     # 10未満の数字を横方向に繋げる
//     # 例：[1,9,4,5]→[1945]
//     num_array2 = []
//     temp = 0
//     num_array.each{|num|
//         if num < 10
//             temp *= 10
//             temp += num
//         else
//             if temp != 0
//                 num_array2.push(temp)
//             else
//                 num_array2.push(1)
//             end
//             num_array2.push(num)
//             temp = 0
//         end
//     }
//     num_array2.push(temp)
//     # 10・100・1000・10000の直前にある数字とで積和する
//     # 例：[2,100,5,10,3]→253
//     val = 0
//     0.upto(num_array2.size / 2 - 1).each{|i|
//         val += num_array2[i * 2] * num_array2[i * 2 + 1]
//     }
//     val += num_array2.last
//     return val
// end


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
	normAddress := norm_addr1(address)

	town := trimStringRegexp(trimedPrefectureCityStr, address)

	etc := trimStringRegexp(trimedPrefectureStr, city)

	fmt.Println(trimedStr)
	fmt.Println(prefecture)
	fmt.Println(city)
	fmt.Println(town)
	fmt.Println(address)
	fmt.Println(normAddress)
	fmt.Println(etc)
}
