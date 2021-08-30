package main

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
)

func main() {
	git, err := gitlab.NewClient("JSqmyRboQpNGvWRzxsDx", gitlab.WithBaseURL("http://192.168.31.220"))
	fmt.Println(err)
	//projects, _, err := git.Projects.ListProjects(nil)
	//git.Projects.GetProject(2, nil)
	//fmt.Println(projects, err)
	var username = "shein"
	users, _, err := git.Users.ListUsers(&gitlab.ListUsersOptions{Username: &username})
	fmt.Println(users[0])
}
