package gnmi_server

import (
	"os/user"
	"common_utils"
	"github.com/msteinert/pam"
	"errors"
	"github.com/golang/glog"
)
type UserCredential struct {
	Username string
	Password string
}

//PAM conversation handler.
func (u UserCredential) PAMConvHandler(s pam.Style, msg string) (string, error) {

	switch s {
	case pam.PromptEchoOff:
		return u.Password, nil
	case pam.PromptEchoOn:
		return u.Password, nil
	case pam.ErrorMsg:
		return "", nil
	case pam.TextInfo:
		return "", nil
	default:
		return "", errors.New("unrecognized conversation message style")
	}
}

// PAMAuthenticate performs PAM authentication for the user credentials provided
func (u UserCredential) PAMAuthenticate() error {
	tx, err := pam.StartFunc("login", u.Username, u.PAMConvHandler)
	if err != nil {
		return err
	}
	return tx.Authenticate(0)
}

func PAMAuthUser(u string, p string) error {

	cred := UserCredential{u, p}
	err := cred.PAMAuthenticate()
	return err
}
func PopulateAuthStruct(username string, auth *common_utils.AuthInfo) error {
	usr, err := user.Lookup(username)
	if err != nil {
		return err
	}

	auth.User = username

	// Get primary group
	group, err := user.LookupGroupId(usr.Gid)
	if err != nil {
		return err
	}
	auth.Group = group.Name

	// Lookup remaining groups
	gids, err := usr.GroupIds()
	if err != nil {
		return err
	}
	auth.Groups = make([]string, len(gids))
	for idx, gid := range gids {
		group, err := user.LookupGroupId(gid)
		if err != nil {
			return err
		}
		auth.Groups[idx] = group.Name
	}

	// TODO: Populate roles list
	return nil
}

func UserPwAuth(username string, passwd string) (bool, error) {
	/*
	 * mgmt-framework container does not have access to /etc/passwd, /etc/group,
	 * /etc/shadow and /etc/tacplus_conf files of host. One option is to share
	 * /etc of host with /etc of container. For now disable this and use ssh
	 * for authentication.
	 */
	err := PAMAuthUser(username, passwd)
	if err != nil {
		glog.Infof("Authentication failed. user=%s, error:%s", username, err.Error())
	    return false, err
	}

	return true, nil
}


func DoesUserExist(username string) bool {
	_, err := user.Lookup(username)
	if err != nil {
		return false
	}
	return true
}