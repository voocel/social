package files

import (
	"encoding/base64"
	"os"
	"regexp"
	"social/pkg/util/snowflake"
)

func ImgFromBase64(prefix, img string) (path string, err error) {
	expr := `^data:\s*image\/(\w+);base64,`
	b, _ := regexp.MatchString(expr, img)
	if !b {
		return
	}

	re, _ := regexp.Compile(expr)
	allData := re.FindAllSubmatch([]byte(img), 2)
	fileType := string(allData[0][1])

	base64Str := re.ReplaceAllString(img, "")

	path = "static/avatar/" + GenFilename(fileType)

	// date := time.Now().Format("2006-01-02")
	// if ok := IsFileExist(path + date); !ok {
	// 	os.Mkdir(path+"/"+date, 0666)
	// }

	bytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(path, bytes, 0666)
	return prefix + "/" + path, nil
}

// IsFileExist 判断文件是否存在
func IsFileExist(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true

}

func GenFilename(ext string) string {
	return snowflake.IDBase32() + ext
}
