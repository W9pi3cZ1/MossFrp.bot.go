package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// 检测是否存在文件的函数
func doesFileExist(fileName string) (isExist bool) {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// 写入文件的函数
func WriteFileString(fileName string, Data string) {
	fa, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	f := base64.NewEncoder(base64.StdEncoding, fa)
	n, err := f.Write([]byte(Data))
	_ = n
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	f.Close()
}

// 读取文件的函数
func ReadFile(fileName string) (data string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	bytesread, err := file.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}

	_ = bytesread

	b, _ := base64.StdEncoding.DecodeString(string(buffer))

	return string(b)
}

func RequestAPI(pwdfilename string) {
	apiUrl := "https://https.ghs.wiki:7002/API?"
	var email string
	var pwd string
	if doesFileExist(pwdfilename) {
		email = strings.Split(ReadFile(pwdfilename), "\n")[0]
		pwd = strings.Split(ReadFile(pwdfilename), "\n")[1]
	} else {
		fmt.Printf("输入邮箱：")
		fmt.Scanln(&email)
		fmt.Printf("输入密码：")
		fmt.Scanln(&pwd)
		WriteFileString(pwdfilename, fmt.Sprintf("%s\n%s", email, pwd))
	}
	queryString := "type=login&loginType=email&account=%s&password=%s"
	queryString = fmt.Sprintf(queryString, email, pwd)
	response, err := http.Get(apiUrl + queryString)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	//签到准备
	Token := gjson.Get(string(body), "token")
	queryString = "type=signIn&token=%s"
	fmt.Scanln(&Token)
	queryString = fmt.Sprintf(queryString, Token)
	fmt.Println(queryString)

	//签到请求
	response, err = http.Get(apiUrl + queryString)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	defer response.Body.Close()

	body, err = io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	fmt.Println(string(body))
}

func main() {
	if doesFileExist("odate") {
		if ReadFile("odate") == time.Now().Format("2006-01-02") {
			fmt.Println("已签到！")
			time.Sleep(2 * time.Second)
		} else {
			os.Remove("odate")
			WriteFileString("odate", time.Now().Format("2006-01-02"))
			RequestAPI("pwd")
			RequestAPI("pwd1")
			RequestAPI("pwd2")
			time.Sleep(2 * time.Second)
		}
	} else {
		WriteFileString("odate", time.Now().Format("2006-01-02"))
		RequestAPI("pwd")
		RequestAPI("pwd1")
		RequestAPI("pwd2")
		time.Sleep(2 * time.Second)
	}
}
