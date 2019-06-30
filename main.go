package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

var do struct {
	Targets []string
	User    string
	UID     string
	GID     string
}

func init() {

	flag.StringVar(&do.User, "u", "", "user")
	flag.StringVar(&do.UID, "U", "", "UID")
	flag.StringVar(&do.GID, "G", "", "GID")
	flag.Parse()

	pFlags := flag.Args()

	if len(pFlags) >= 1 {
		for _, fl := range pFlags {
			do.Targets = append(do.Targets, getAbsolutePath(fl))
		}
	} else {
		errExit(127, "missing operand")
	}

	var u *user.User

	if do.UID == "" {
		if do.User == "" {
			x, err := user.Current()
			if err != nil {
				errExit(2, err.Error())
			}
			u = x
		} else {
			x, err := user.Lookup(do.User)
			if err != nil {
				errExit(2, err.Error())
			}
			u = x
		}
		do.UID = u.Uid
	} else {
		x, err := user.LookupId(do.UID)
		if err != nil {
			errExit(2, err.Error())
		}
		u = x
		do.UID = u.Uid
	}

	if do.GID == "" {
		do.GID = u.Gid
	}

}

func main() {

	u := strToInt(do.UID)
	g := strToInt(do.GID)

	for _, target := range do.Targets {

		f, err := os.Stat(target)
		if err != nil {
			errExit(2, err.Error())
		}

		if f.IsDir() {
			fmt.Println("\ndir:\t" + target)
			files := scanFolder(target)
			for _, file := range files {
				setOwner(file, u, g)
				fmt.Println(file)
			}
		} else {
			fmt.Println("\nfile:\t" + target)
			setOwner(target, u, g)
		}

	}
}

func scanFolder(d string) (files []string) {
	err := filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		errExit(1, err.Error())
		return nil
	}
	return files
}

func getAbsolutePath(t string) string {
	absolute, err := filepath.Abs(t)
	if err != nil {
		errExit(2, err.Error())
		return ""
	}
	return absolute
}

func setOwner(t string, u int, g int) {
	err := os.Chown(t, u, g)
	if err != nil {
		errExit(2, err.Error())
	}
}

func strToInt(s string) int {
	u, err := strconv.Atoi(s)
	if err != nil {
		errExit(2, err.Error())
	}
	return u
}

func errExit(code int, in interface{}) {
	fmt.Println(in)
	os.Exit(code)
}
