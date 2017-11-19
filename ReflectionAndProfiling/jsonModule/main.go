package jsonModule

import (
	"encoding/json"
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

	users := []User{}
	scanner := bufio.NewScanner(file)
	
	user := User{}
	for scanner.Scan() {
		user = User{}
		// fmt.Printf("%v %v\n", err, line)
		line := scanner.Text()
		err := json.Unmarshal([]byte(line), &user)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}

	if err = scanner.Err(); err != nil {
		panic(err)
	}

	seenBrowsers := []string{}
	uniqueBrowsers := 0
	foundUsers := []string{"found users:\n",}

	for i, user := range users {
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
		foundUsers = append(foundUsers, fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email))
	}

	fmt.Fprintln(out, strings.Join(foundUsers,""))
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
