package matterclient

import (
	"errors"
	"time"

	"github.com/mattermost/mattermost-server/model"
)

func (m *MMClient) GetNickName(userId string) string { //nolint:golint
	user := m.GetUser(userId)
	if user != nil {
		return user.Nickname
	}
	return ""
}

func (m *MMClient) GetStatus(userId string) string { //nolint:golint
	res, resp := m.Client.GetUserStatus(userId, "")
	if resp.Error != nil {
		return ""
	}
	if res.Status == model.STATUS_AWAY {
		return "away"
	}
	if res.Status == model.STATUS_ONLINE {
		return "online"
	}
	return "offline"
}

func (m *MMClient) GetTeamId() string { //nolint:golint
	return m.Team.Id
}

// GetTeamName returns the name of the specified teamId
func (m *MMClient) GetTeamName(teamId string) string { //nolint:golint
	m.RLock()
	defer m.RUnlock()
	for _, t := range m.OtherTeams {
		if t.Id == teamId {
			return t.Team.Name
		}
	}
	return ""
}

func (m *MMClient) GetUser(userId string) *model.User { //nolint:golint
	res, resp := m.Client.GetUser(userId, "")
	if resp.Error != nil {
		return nil
	}
	return res
}

func (m *MMClient) GetUserName(userId string) string { //nolint:golint
	user := m.GetUser(userId)
	if user != nil {
		return user.Username
	}
	return ""
}

func (m *MMClient) UpdateUserNick(nick string) error {
	user := m.User
	user.Nickname = nick
	_, resp := m.Client.UpdateUser(user)
	if resp.Error != nil {
		return resp.Error
	}
	return nil
}

func (m *MMClient) UsernamesInChannel(channelId string) ([]string, error) { //nolint:golint
	usernames := []string{}
	for page := 0; true; page++ {
		mmusers, resp := m.Client.GetUsersInChannel(channelId, page, 200, "")
		if resp.Error != nil {
			m.logger.Errorf("UsernamesInChannel(%s) failed: %s", channelId, resp.Error)
			return nil, errors.New(resp.Error.DetailedError)
		}

		// end of users
		if len(mmusers) == 0 {
			break
		}

		for _, user := range mmusers {
			usernames = append(usernames, user.Username)
		}

		time.Sleep(time.Millisecond * 200)
	}

	return usernames, nil
}

func (m *MMClient) UpdateStatus(userId string, status string) error { //nolint:golint
	_, resp := m.Client.UpdateUserStatus(userId, &model.Status{Status: status})
	if resp.Error != nil {
		return resp.Error
	}
	return nil
}
