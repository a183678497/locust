package utils

import (
	"encoding/binary"
	"github.com/axgle/mahonia"
	"strconv"
	"unsafe"
)


/*
字符串转bytes
*/
func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

/*
bytes转字符串
*/
func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Uint16ToBytes(i uint16) []byte {
	var buf = make([]byte, 2)
	binary.BigEndian.PutUint16(buf, i)
	return buf
}

func BytesToUint16(buf []byte) uint16 {
	return binary.BigEndian.Uint16(buf)
}

func Uint32ToBytes(i uint32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, i)
	return buf
}

func BytesToUint32(buf []byte) uint32 {
	return binary.BigEndian.Uint32(buf)
}

func Uint64ToBytes(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

func BytesToUint64(buf []byte) uint64 {
	return binary.BigEndian.Uint64(buf)
}

func BcdToStr(buf []byte) string{
	var res string
	for _,v := range buf {
		a := (v&0xf0)>>4
		b := v&0x0f
		res=res+strconv.Itoa(int(a))+strconv.Itoa(int(b))
	}
	return res
}


/*
字符编码切换
demo
	s := "驾驶员疲劳报警"
	b := str2bytes(s)
	c := hex.EncodeToString(b)
	test, _ := hex.DecodeString("BBC6CFD4BBAA20202020")
	s3 := bytes2str(test)
	output := ConvertToString(s3, "gb18030", "utf-8")
	fmt.Println(b, c, output)
*/
func ConvertToString(src string, srcCode string, tagCode string) string {

	srcCoder := mahonia.NewDecoder(srcCode)

	srcResult := srcCoder.ConvertString(src)

	tagCoder := mahonia.NewDecoder(tagCode)

	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)

	result := string(cdata)

	return result

}