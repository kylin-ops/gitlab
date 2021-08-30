// 接口来源https://docs.gitlab.com/ee/api/projects.html#project-visibility-level

package sdk

import (
	"errors"
	"github.com/xanzy/go-gitlab"
	"strings"
)

func NewClient(address, token string) (*Client, error) {
	cli, err := gitlab.NewClient(token, gitlab.WithBaseURL(address))
	return &Client{Client: cli}, err
}

type Client struct {
	Client *gitlab.Client
}

func (c *Client) usernameToUserId(username string) (userId int, err error) {
	users, _, err := c.Client.Users.ListUsers(&gitlab.ListUsersOptions{Username: &username})
	if err != nil {
		return userId, err
	}
	for _, user := range users {
		if user.Username == username {
			userId = user.ID
			break
		}
	}
	if userId == 0 {
		err = errors.New("The specified username \"" + username + "\" does not exist")
	}
	return userId, err
}

// 获取所有的项目
func (c *Client) ProjectListAll() (projects []*gitlab.Project, err error) {
	projects, _, err = c.Client.Projects.ListProjects(nil)
	return projects, err
}

func (c *Client) getVisibilityValue(visibilityStr string) (visibility gitlab.VisibilityValue, err error) {
	switch visibilityStr {
	case "public":
		visibility = gitlab.PublicVisibility
	case "private":
		visibility = gitlab.PrivateVisibility
	case "internal":
		visibility = gitlab.InternalVisibility
	default:
		err = errors.New("The value of \"projectType\" can only be \"public, private, internal\"")
	}
	return visibility, err
}

func (c *Client) getAccessLevel(accessLevel string) (level gitlab.AccessLevelValue, err error) {
	switch accessLevel {
	case "no":
		level = gitlab.NoPermissions
	case "minimal":
		level = gitlab.MinimalAccessPermissions
	case "guest":
		level = gitlab.GuestPermissions
	case "reporter":
		level = gitlab.ReporterPermissions
	case "developer":
		level = gitlab.DeveloperPermissions
	case "maintainer":
		level = gitlab.MaintainerPermissions
	case "owner":
		level = gitlab.OwnerPermission
	default:
		err = errors.New("The value of \"accessLevel\" can only be \"no, minimal, guest, reporter, developer, maintainer, owner\"")
	}
	return level, err
}

// 获取指定用于的项目
func (c *Client) ProjectListFromUsername(username string) (projects []*gitlab.Project, err error) {
	userId, err := c.usernameToUserId(username)
	if err != nil {
		return nil, err
	}
	projects, _, err = c.Client.Projects.ListUserProjects(userId, nil)
	return projects, err
}

// 创建基本项目
func (c *Client) ProjectCreate(projectName, visibilityStr string) (project *gitlab.Project, err error) {
	visibility, err := c.getVisibilityValue(visibilityStr)
	if err != nil {
		return nil, err
	}
	project, _, err = c.Client.Projects.CreateProject(&gitlab.CreateProjectOptions{
		Name:       &projectName,
		Path:       &projectName,
		Visibility: &visibility,
	})
	return project, err
}

// 为项目添加多个项目成员
func (c *Client) ProjectSetMember(projectId int, accessLevel string, users []string) error {
	level, err := c.getAccessLevel(accessLevel)
	var userIds []int
	var userMap = map[int]string{}
	var errUsers []string
	if err != nil {
		return err
	}
	for _, user := range users {
		userId, err := c.usernameToUserId(user)
		if err != nil {
			errUsers = append(errUsers, user)
			continue
		}
		userIds = append(userIds, userId)
		userMap[userId] = user
	}
	if len(errUsers) > 0 {
		return errors.New("\"" + strings.Join(errUsers, ",") + "\"" + " username does not exist in Gitlab")
	}
	for _, userId := range userIds {
		_, _, err = c.Client.ProjectMembers.AddProjectMember(projectId, &gitlab.AddProjectMemberOptions{
			UserID:      userId,
			AccessLevel: &level,
		})
		if err != nil {
			errUsers = append(errUsers, userMap[userId])
		}
	}
	if len(errUsers) > 0 {
		return errors.New("\"" + strings.Join(errUsers, ",") + "\"" + " Failed to add gitLab project member")
	}
	return nil
}

// 创建基本项目并指定成员
func (c *Client) ProjectCreateSetMember(projectName, visibilityStr, accessLevel string, users []string) (*gitlab.Project, error) {
	project, err := c.ProjectCreate(projectName, visibilityStr)
	if err != nil {
		return project, err
	}
	return project, c.ProjectSetMember(project.ID, accessLevel, users)
}
