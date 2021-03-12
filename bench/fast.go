package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	seenBrowsers := map[string]bool{}
	foundUsers := strings.Builder{}

	var user User
	id := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err = user.Init(id, scanner.Bytes()); err != nil {
			panic(err)
		}

		id += 1

		isAndroid := false
		isMSIE := false

		for _, browser := range user.Browsers {
			if strings.Contains(browser, "Android") {
				isAndroid = true
				seenBrowsers[browser] = true
			}

			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				seenBrowsers[browser] = true
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		foundUsers.WriteString(user.String())
		foundUsers.WriteString("\n")
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers.String())
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}

type User struct {
	id       int
	Browsers []string `json:"browsers"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
}

func (u *User) Init(id int, data []byte) error {
	u.id = id
	return json.Unmarshal(data, u)
}

func (u *User) String() string {
	email := strings.Replace(u.Email, "@", " [at] ", 1)
	return fmt.Sprintf("[%d] %s <%s>", u.id, u.Name, email)
}
