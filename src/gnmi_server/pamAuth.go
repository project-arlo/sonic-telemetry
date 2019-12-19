package gnmi_server

import (
	"golang.org/x/crypto/ssh"
	"os/user"
)


func UserPwAuth(username string, passwd string) (bool, error) {
	/*
	 * mgmt-framework container does not have access to /etc/passwd, /etc/group,
	 * /etc/shadow and /etc/tacplus_conf files of host. One option is to share
	 * /etc of host with /etc of container. For now disable this and use ssh
	 * for authentication.
	 */
	/* err := PAMAuthUser(username, passwd)
	    if err != nil {
			log.Printf("Authentication failed. user=%s, error:%s", username, err.Error())
	        return err
	    }*/

	//Use ssh for authentication.
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(passwd),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	_, err := ssh.Dial("tcp", "127.0.0.1:22", config)
	if err != nil {
		return false, err
	}

	return true, nil
}

func IsAdminGroup(username string) bool {

	usr, err := user.Lookup(username)
	if err != nil {
		return false
	}
	gids, err := usr.GroupIds()
	if err != nil {
		return false
	}
	
	admin, err := user.Lookup("admin")
	if err != nil {
		return false
	}
	for _, x := range gids {
		if x == admin.Gid {
			return true
		}
	}
	return false
}

func DoesUserExist(username string) bool {
	_, err := user.Lookup(username)
	if err != nil {
		return false
	}
	return true
}