package test

import (
	"fmt"
	"html"
	"os"
	"testing"
)

func TestXss(t *testing.T) {
	rawString := `<script>alert("xss")</script>`
	safeString := html.EscapeString(rawString) //escape:编码转义

	file, err := os.Create("./xssexample.html") //浏览器查看.html文件，rawstring有弹窗，safeString正常显示
	if err != nil {
		fmt.Println(err)
	}
	file.WriteString("Unsafe html:")
	file.WriteString(rawString + "\n")
	file.WriteString("Safe html:")
	file.WriteString(safeString)
}
