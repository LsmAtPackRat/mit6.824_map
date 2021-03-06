package main

import (
	"fmt"
	"mapreduce"
	"os"
    "strconv"
    "strings"
    "unicode"
)

//
// The map function is called once for each file of input. The first
// argument is the name of the input file, and the second is the
// file's complete contents. You should ignore the input file name,
// and look only at the contents argument. The return value is a slice
// of key/value pairs.
//
func mapF(filename string, contents string) []mapreduce.KeyValue {
	// Your code here (Part II).
    // split contents into a string slice
    //fmt.Println(contents)
    //words := strings.Fields(contents)
    f := func(c rune) bool {
        return !unicode.IsLetter(c)
    }
    words := strings.FieldsFunc(contents, f)
    var result []mapreduce.KeyValue
    word_count := make(map[string]int)
    for _, word := range words {
        // count the number of word in contents
        count, ok := word_count[word]
        if !ok {
            word_count[word] = 1
        } else {
            word_count[word] = count + 1
        }
    }
    //fmt.Println(result)
    // word_count -> result
    for k, v := range word_count {
        var kv mapreduce.KeyValue
        kv.Key = k
        kv.Value = strconv.Itoa(v)
        result = append(result, kv)
    }
    return result
}

//
// The reduce function is called once for each key generated by the
// map tasks, with a list of all the values created for that key by
// any map task.
//
func reduceF(key string, values []string) string {
	// Your code here (Part II).
    count := 0
    for _, str_num := range values {
        part, _ := strconv.Atoi(str_num)
        count += part
    }
    //fmt.Println("count = %s", string(count))
    return strconv.Itoa(count)
}

// Can be run in 3 ways:
// 1) Sequential (e.g., go run wc.go master sequential x1.txt .. xN.txt)
// 2) Master (e.g., go run wc.go master localhost:7777 x1.txt .. xN.txt)
// 3) Worker (e.g., go run wc.go worker localhost:7777 localhost:7778 &)
func main() {
	if len(os.Args) < 4 {
		fmt.Printf("%s: see usage comments in file\n", os.Args[0])
	} else if os.Args[1] == "master" {
		var mr *mapreduce.Master
		if os.Args[2] == "sequential" {
			mr = mapreduce.Sequential("wcseq", os.Args[3:], 3, mapF, reduceF)
		} else {
			mr = mapreduce.Distributed("wcseq", os.Args[3:], 3, os.Args[2])
		}
		mr.Wait()
	} else {
		mapreduce.RunWorker(os.Args[2], os.Args[3], mapF, reduceF, 100, nil)
	}
}
