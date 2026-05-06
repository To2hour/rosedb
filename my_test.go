package rosedb

import (
	"encoding/binary"
	"fmt"
	"os"
	"testing"
)

func TestName(t *testing.T) {
	dir, err := os.ReadDir("./")
	if err != nil {
		return
	}
	for _, value := range dir {
		fmt.Println(value)
	}
}
func TestBinaryUvarint(t *testing.T) {
	// 构造数据：
	// 'a' (97) -> [0x61] (最高位是0，1字节读完)
	// 300      -> [0xAC, 0x02] (0xAC最高位是1，继续读；0x02最高位是0，结束)
	str := "hello world,123455"
	var buf []byte = []byte(str)
	idx := 0
	for {
		// 尝试从当前位置 idx 开始解码
		val, n := binary.Uvarint(buf[idx:])

		if n > 0 {
			fmt.Printf("【成功解析】\n")
			fmt.Printf("  解析出的数值: %d\n", val)
			fmt.Printf("  该数值占用字节数: %d\n", n)
			idx += n // 移动光标
			fmt.Printf("  当前光标位置 (idx): %d\n", idx)
			fmt.Printf("  剩余待处理字节: %v\n", buf[idx:])
		} else if n == 0 {
			fmt.Printf("【解析中断】数据不足或已读完，剩余字节: %v\n", buf[idx:])
			break
		} else {
			fmt.Printf("【解析错误】数据溢出或损坏\n")
			break
		}
		fmt.Println("--------------------")
	}
}
func TestBinaryUvarint2(t *testing.T) {
	// 手工构造一个字节数组
	// 0x80, 0x01 是数字 128 的 Varint 编码（占 2 字节）
	// 0xff 是随后的无关数据
	buf := []byte{0x80, 0x01, 0xff, 0xee, 0xdd}

	val, n := binary.Uvarint(buf)

	fmt.Printf("读到的值: %d\n", val) // 输出 128
	fmt.Printf("消耗的字节数: %d\n", n) // 输出 2

	// 验证：buf 剩下的部分完全没被改动，也没有被 Uvarint 关心
	fmt.Printf("剩下的数据: %v\n", buf[n:]) // 输出 [255 238 221]
}
func TestBigValue(t *testing.T) {
	// 1. 准备一个数值为 1 亿的数
	var originalVal uint64 = 100000000
	buf := make([]byte, binary.MaxVarintLen64) // MaxVarintLen64 是 10

	// 2. 编码
	n := binary.PutUvarint(buf, originalVal)
	fmt.Printf("数值 %d 编码后的字节: %v, 长度: %d\n", originalVal, buf[:n], n)

	// 3. 解码
	decodedVal, readNum := binary.Uvarint(buf[:n])
	fmt.Printf("解码出的值: %d, 读取字节: %d\n", decodedVal, readNum)
}
