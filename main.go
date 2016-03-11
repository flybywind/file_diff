package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	var (
		op int
		sep,
		keys string
	)
	flag.IntVar(&op, "o", 0, "operation, `0` for intersect, `1` for subfrom first, `2` for subfrom second")
	flag.StringVar(&sep, "s", " ", "seperator, space for default")
	flag.StringVar(&keys, "k", "1", "keys for operation, default is 1, the first field. use comma to represent multiple keys such as : '0,1'")

	flag.Parse()
	argc := len(os.Args)
	if argc < 3 || !flag.Parsed() {
		flag.Usage()
		os.Exit(-1)
	}

	var (
		f1, f2 *os.File
		err    error
		sc     *bufio.Scanner
		key, line0,
		line string
		seg    []string
		exists bool
	)

	f1, err = os.Open(os.Args[argc-2])
	if err != nil {
		log.Fatalln("open", os.Args[argc-2], "failed:", err)
	}
	f2, err = os.Open(os.Args[argc-1])
	if err != nil {
		log.Fatalln("open", os.Args[argc-1], "failed:", err)
	}

	buf_map := map[string]string{}
	output := []string{}
	key_list_str := strings.Split(keys, ",")
	key_list := []int64{}
	var idx int64
	for _, i := range key_list_str {
		if idx, err = strconv.ParseInt(i, 10, 64); err != nil {
			log.Fatalln("parse key list failed:", err)
		}
		key_list = append(key_list, idx)
	}

	sc = bufio.NewScanner(f1)
	for sc.Scan() {
		line = sc.Text()
		seg = strings.Split(line, sep)
		if key, err = concatField(seg, key_list); err != nil {
			log.Fatalf("process %s failed: %v\n", os.Args[argc-2], err)
		} else {
			buf_map[key] = line
		}
	}
	sc = bufio.NewScanner(f2)
	for sc.Scan() {
		line = sc.Text()
		seg = strings.Split(line, sep)
		if key, err = concatField(seg, key_list); err != nil {
			log.Fatalf("process %s failed: %v\n", os.Args[argc-1], err)
		}
		switch op {
		case 0:
			if _, exists = buf_map[key]; exists {
				output = append(output, line)
			}
		case 1:
			// 从file1中减去，就是输出file2中剩下的行
			if _, exists = buf_map[key]; !exists {
				output = append(output, line)
			}
		case 2:
			// 从file2中减去，就是输出file2中不包含的file1中的行:
			if line0, exists = buf_map[key]; !exists {
				output = append(output, line0)
			}
		}
	}

	for _, l := range output {
		fmt.Println(l)
	}
}

func concatField(seg []string, idx []int64) (key string, err error) {
	buf := bytes.Buffer{}
	N := int64(len(seg))
	for _, i := range idx {
		if i < N {
			buf.WriteString(seg[i] + "|")
		} else {
			return "", fmt.Errorf("%s doesn't containt enough key", strings.Join(seg, ", "))
		}
	}
	return buf.String(), nil
}
