package main

import (
	"github.com/mailru/easyjson"
	"fmt"
	"io"
	"bufio"
	"os"
	"strings"
	// "log"
)

//easyjson:json
type User struct {
	Name string
	Browsers []string
	Email string
}

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	seenBrowsers := []string{}
	uniqueBrowsers := 0
	fmt.Fprintf(out, "found users:\n")

	i := -1
	user := User{}
	reader := bufio.NewReader(file)
	for {
		i += 1
		user = User{}
		// fmt.Printf("%v %v\n", err, line)
		line, isPrefix, err := reader.ReadLine()
		if err != nil || isPrefix == true {
			break
		}

		err = easyjson.Unmarshal(line, &user)

		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		for _, browserRaw := range user.Browsers {
			browser := browserRaw

			notSeenBefore := false
			if strings.Contains(browser, "Android") {
				notSeenBefore = true
				isAndroid = true
			}

			if strings.Contains(browser, "MSIE") {
				notSeenBefore = true
				isMSIE = true
			}

			for _, item := range seenBrowsers {
				if item == browser {
					notSeenBefore = false
				}
			}
			if notSeenBefore {
				// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
				seenBrowsers = append(seenBrowsers, browser)
				uniqueBrowsers++
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := strings.Replace(user.Email, "@", " [at] ", -1)
		fmt.Fprintf(out, "[%d] %s <%s>\n", i, user.Name, email)
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
