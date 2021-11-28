package main

import (
	"fmt"
	"os"
)

var VERSION = "0.0.1"

func showUsage() {
	fmt.Println("USAGE:  go run esab46.go '文字列'")
	os.Exit(1)
}

func main() {

	if len(os.Args) != 2 {
		showUsage()
	}
	originText := os.Args[1]

	fmt.Println(base64encode(originText))
}

func base64encode(originText string) string {
	originTextByteList := []byte(originText)

	//元の文字列のバイト配列でどこまで読み込みしたかのindex
	var originTextIndex int

	//3バイト(24bit)ごとに処理するため、3で割ってあまりを捨てた後の文字の長さを取得
	var threeByteBlockNum int = (len(originTextByteList) / 3) * 3

	var resultByte []byte //slice
	/**
	 * 3バイト単位に処理し、処理できなかった残りの1バイトもしくは2バイトがある場合は、このループの後に処理する
	 */
	for originTextIndex = 0; originTextIndex < threeByteBlockNum; originTextIndex += 3 {
		// 最初の3バイトを配列から取得し、最初の1バイト目を16bitシフトさせて24bitにする。
		// 例: 1バイト目が11111111だったら16bitシフトで 11111111 00000000 00000000 となる
		// 例: 2バイト目は11111111だったら 8bitシフトで          11111111 00000000 となる
		// 例: 3バイト目は11111111だったら 0bitシフトで                   11111111 となる
		// これらをorで演算すると、                    11111111 11111111 11111111  となる
		// これで3バイトのデータを連結したビットの情報が得られる
		first := uint(originTextByteList[originTextIndex]) << 16
		second := uint(originTextByteList[originTextIndex+1]) << 8
		third := uint(originTextByteList[originTextIndex+2])
		bit24 := first | second | third

		//3バイト連結データbit24変数(uint型)から4つの6bitデータに分割する
		//最初の情報は18bitシフトさせると6bitデータになる
		//次の情報は12bitシフトさせると12bitデータになる
		//   bit24が              11111111 00111111 11111111 の場合
		//   2番目の情報は          111111[11 0011]1111 11111111 カッコの値がほしい
		//   12bitシフトさせると、   00000000 00001111 11[110011]
		//   6bit取得のためAND演算  00000000 00000000 00111111  <- 0b00111111
		//   結果                 00000000 00000000 00110011
		first6bit := bit24 >> 18 & 0b00111111
		second6bit := bit24 >> 12 & 0b00111111
		third6bit := bit24 >> 6 & 0b00111111
		forth6bit := bit24 >> 0 & 0b00111111

		resultByte = append(resultByte, getChar(first6bit))
		resultByte = append(resultByte, getChar(second6bit))
		resultByte = append(resultByte, getChar(third6bit))
		resultByte = append(resultByte, getChar(forth6bit))
	}

	// ここからは残りの余ったバイト数が0,1,2のパターンとなる
	// 0バイト余った状態は何もしない
	remainNumChar := len(originTextByteList) - originTextIndex
	if remainNumChar == 0 {
		return string(resultByte)
	}

	// 1バイト余った状態
	// 8bitに0000の4bitを足して12bitにして、6bit, 6bitの64encodeの2文字を取得
	// エンコード後は4文字単位のため、残りの2文字は==を連結
	if remainNumChar == 1 {
		first := uint(originTextByteList[originTextIndex]) << 4
		first6bit := first >> 6 & 0b00111111
		second6bit := first >> 0 & 0b00111111
		resultByte = append(resultByte, getChar(first6bit))
		resultByte = append(resultByte, getChar(second6bit))
		resultByte = append(resultByte, '=')
		resultByte = append(resultByte, '=')
	}

	// 2バイト余った状態
	// 16bitに00の2bitを足して18bitにして、6bit, 6bit, 6bitの64encodeの3文字を取得
	// エンコード後は4文字単位のため、残りの1文字は=を連結
	if remainNumChar == 2 {
		first := uint(originTextByteList[originTextIndex]) << 10
		second := uint(originTextByteList[originTextIndex+1]) << 2
		bit18 := first | second
		first6bit := bit18 >> 12 & 0b00111111
		second6bit := bit18 >> 6 & 0b00111111
		third6bit := bit18 >> 0 & 0b00111111
		resultByte = append(resultByte, getChar(first6bit))
		resultByte = append(resultByte, getChar(second6bit))
		resultByte = append(resultByte, getChar(third6bit))
		resultByte = append(resultByte, '=')
	}
	return string(resultByte)
}

/**
 * 0-63までの数字に対応するBase64の文字列を取得する
 * 対応文字列はこちらを参照
 * https://ja.wikipedia.org/wiki/Base64
 */
func getChar(position uint) byte {
	encodeCharList := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	return encodeCharList[position]
}
