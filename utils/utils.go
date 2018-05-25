package utils

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/fatih/color"

	"github.com/iawia002/annie/config"
	"github.com/iawia002/annie/request"
)

// MAXLENGTH Maximum length of file name
const MAXLENGTH = 80

// MatchOneOf match one of the patterns
func MatchOneOf(text string, patterns ...string) []string {
	var (
		re    *regexp.Regexp
		value []string
	)
	for _, pattern := range patterns {
		// (?flags): set flags within current group; non-capturing
		// s: let . match \n (default false)
		// https://github.com/google/re2/wiki/Syntax
		re = regexp.MustCompile(pattern)
		value = re.FindStringSubmatch(text)
		if len(value) > 0 {
			return value
		}
	}
	return nil
}

// MatchAll return all matching results
func MatchAll(text, pattern string) [][]string {
	re := regexp.MustCompile(pattern)
	value := re.FindAllStringSubmatch(text, -1)
	return value
}

// FileSize return the file size of the specified path file
func FileSize(filePath string) (int64, bool) {
	file, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, false
		}
		log.Fatal(err)
	}
	return file.Size(), true
}

// Domain get the domain of given URL
func Domain(url string) string {
	domainPattern := `([a-z0-9][-a-z0-9]{0,62})\.` +
		`(com\.cn|com\.hk|` +
		`cn|com|net|edu|gov|biz|org|info|pro|name|xxx|xyz|be|` +
		`me|top|cc|tv|tt)`
	domain := MatchOneOf(url, domainPattern)
	if domain != nil {
		return domain[1]
	}
	return "Universal"
}

// LimitLength Handle overly long strings
func LimitLength(s string, length int) string {
	const ELLIPSES = "..."
	str := []rune(s)
	if len(str) > length {
		return string(str[:length-len(ELLIPSES)]) + ELLIPSES
	}
	return s
}

// FileName Converts a string to a valid filename
func FileName(name string) string {
	rep := strings.NewReplacer("\n", " ", "/", " ", "|", "-", ": ", "：", ":", "：", "'", "’")
	name = rep.Replace(name)
	if runtime.GOOS == "windows" {
		rep = strings.NewReplacer("\"", " ", "?", " ", "*", " ", "\\", " ", "<", " ", ">", " ")
		name = rep.Replace(name)
	}
	return LimitLength(name, MAXLENGTH)
}

// FilePath gen valid file path
func FilePath(name, ext string, escape bool) string {
	var outputPath string
	if config.OutputPath != "" {
		_, err := os.Stat(config.OutputPath)
		if err != nil && os.IsNotExist(err) {
			log.Fatal(err)
		}
	}
	fileName := fmt.Sprintf("%s.%s", name, ext)
	if escape {
		fileName = FileName(fileName)
	}
	outputPath = filepath.Join(config.OutputPath, fileName)
	return outputPath
}

// ItemInSlice if a item is in the list
func ItemInSlice(item, list interface{}) bool {
	v1 := reflect.ValueOf(item)
	v2 := reflect.ValueOf(list)
	for i := 0; i < v2.Len(); i++ {
		indexType := v2.Index(i).Type().String()
		if v1.Type().String() != indexType {
			continue
		}
		switch indexType {
		case "int":
			if v1.Int() == v2.Index(i).Int() {
				return true
			}
		case "string":
			if v1.String() == v2.Index(i).String() {
				return true
			}
		}
	}
	return false
}

// GetNameAndExt return the name and ext of the URL
// https://img9.bcyimg.com/drawer/15294/post/1799t/1f5a87801a0711e898b12b640777720f.jpg ->
// 1f5a87801a0711e898b12b640777720f, jpg
func GetNameAndExt(uri string) (string, string) {
	u, _ := url.ParseRequestURI(uri)
	s := strings.Split(u.Path, "/")
	filename := strings.Split(s[len(s)-1], ".")
	if len(filename) > 1 {
		return filename[0], filename[1]
	}
	// Image url like this
	// https://img9.bcyimg.com/drawer/15294/post/1799t/1f5a87801a0711e898b12b640777720f.jpg/w650
	// has no suffix
	contentType := request.ContentType(uri, uri)
	return filename[0], strings.Split(contentType, "/")[1]
}

// Md5 md5 hash
func Md5(text string) string {
	sign := md5.New()
	sign.Write([]byte(text))
	return fmt.Sprintf("%x", sign.Sum(nil))
}

// M3u8URLs get all urls from m3u8 url
func M3u8URLs(uri string) []string {
	html := request.Get(uri, "", nil)
	lines := strings.Split(html, "\n")
	var urls []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "http") {
				urls = append(urls, line)
			} else {
				base, _ := url.Parse(uri)
				u, _ := url.Parse(line)
				urls = append(urls, fmt.Sprintf("%s", base.ResolveReference(u)))
			}
		}
	}
	return urls
}

// PrintVersion print version information
func PrintVersion() {
	blue := color.New(color.FgBlue)
	cyan := color.New(color.FgCyan)
	fmt.Printf(
		"\n%s: version %s, A fast, simple and clean video downloader.\n\n",
		cyan.Sprintf("annie"),
		blue.Sprintf(config.VERSION),
	)
}

// Reverse Reverse a string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Range generate a sequence of numbers by range
func Range(min, max int) []int {
	items := make([]int, max-min+1)
	for index := range items {
		items[index] = min + index
	}
	return items
}
