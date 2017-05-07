package backup

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"github.com/stefanprodan/mgob/config"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"path/filepath"
)

func sftpUpload(file string, plan config.Plan) (string, error) {
	sshConf := &ssh.ClientConfig{
		User: plan.SFTP.Username,
		Auth: []ssh.AuthMethod{ssh.Password(plan.SFTP.Password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	sshCon, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", plan.SFTP.Host, plan.SFTP.Port), sshConf)
	if err != nil {
		return "", errors.Wrapf(err, "SSH dial to %v:%v failed", plan.SFTP.Host, plan.SFTP.Port)
	}
	defer sshCon.Close()

	sftpClient, err := sftp.NewClient(sshCon)
	if err != nil {
		return "", errors.Wrapf(err, "SFTP client init %v:%v failed", plan.SFTP.Host, plan.SFTP.Port)
	}
	defer sftpClient.Close()

	f, err := os.Open(file)
	if err != nil {
		return "", errors.Wrapf(err, "Opening file %v failed", file)
	}
	defer f.Close()

	_, fname := filepath.Split(file)
	dstPath := filepath.Join(plan.SFTP.Dir, fname)
	sf, err := sftpClient.Create(dstPath)
	if err != nil {
		return "", errors.Wrapf(err, "SFTP %v:%v creating file %v failed", plan.SFTP.Host, plan.SFTP.Port, dstPath)
	}

	_, err = io.Copy(sf, f)
	if err != nil {
		return "", errors.Wrapf(err, "SFTP %v:%v upload file %v failed", plan.SFTP.Host, plan.SFTP.Port, dstPath)
	}
	sf.Close()

	//listSftpBackups(sftpClient, plan.SFTP.Dir)

	return fmt.Sprintf("SFTP %v:%v upload %v finished", plan.SFTP.Host, plan.SFTP.Port, dstPath), nil
}

func listSftpBackups(client *sftp.Client, dir string) error {
	list, err := client.ReadDir(fmt.Sprintf("/%v", dir))
	if err != nil {
		return errors.Wrapf(err, "SFTP reading %v dir failed", dir)
	}

	for _, item := range list {
		fmt.Println(item.Name())
	}

	return nil
}
