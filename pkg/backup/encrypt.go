package backup

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/codeskyblue/go-sh"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stefanprodan/mgob/pkg/config"
)

func encrypt(file string, encryptedFile string, plan config.Plan, conf *config.AppConfig) (string, error) {
	if plan.Encryption.Gpg != nil {
		if !conf.HasGpg {
			return "", errors.Errorf("GPG configuration is present, but no GPG binary is found! Uploading unencrypted backup.")
		}
		return gpgEncrypt(file, encryptedFile, plan)
	}

	return "", errors.Errorf("Encryption config is not valid!")
}

func removeUnencrypted(file string, encryptedFile string) {
	// Check if encrypted file exists and remove original
	stat, err := os.Stat(encryptedFile)
	if err == nil && stat.Size() > 0 {
		os.Remove(file)
	}
}

func gpgEncrypt(file string, encryptedFile string, plan config.Plan) (string, error) {
	output := ""
	recipient := ""

	recipients := plan.Encryption.Gpg.Recipients

	keyFile := plan.Encryption.Gpg.KeyFile
	if keyFile != "" {
		keyFileStat, err := os.Stat(keyFile)
		if err == nil && !keyFileStat.IsDir() {
			// import key from file
			importCmd := fmt.Sprintf("gpg --batch --import  %v", keyFile)

			result, err := sh.Command("/bin/sh", "-c", importCmd).CombinedOutput()
			if len(result) > 0 {
				output += strings.Replace(string(result), "\n", " ", -1)
			}
			if err != nil {
				return "", errors.Wrapf(err, "Importing encryption key for plan %v failed %s", plan.Name, output)
			}
			if !strings.Contains(output, "imported: 1") && !strings.Contains(output, "unchanged: 1") {
				return "", errors.Errorf("Importing encryption key failed %v", output)
			}

			re := regexp.MustCompile(`key ([0-9A-F]+):`)
			keyMatch := re.FindStringSubmatch(output)
			log.WithField("plan", plan.Name).Debugf("Import output: %v", output)
			log.WithField("plan", plan.Name).Debugf("Parsed key id: %v", keyMatch[1])
			if keyMatch != nil {
				recipients = append(recipients, keyMatch[1])
			}
		}
	}

	recipient = strings.Join(recipients, " -r ")

	if recipient == "" {
		return "", errors.Errorf("GPG configuration is present, but no encryption key is configured! %v", output)
	}

	keyServer := plan.Encryption.Gpg.KeyServer
	if keyServer == "" {
		keyServer = "hkps://keys.openpgp.org"
	}

	// encrypt file
	encryptCmd := fmt.Sprintf(
		"gpg -v --batch --yes --trust-model always --auto-key-locate local,%v -e -r %v -o %v %v",
		keyServer, recipient, encryptedFile, file)

	result, err := sh.Command("/bin/sh", "-c", encryptCmd).CombinedOutput()
	if len(result) > 0 {
		output += strings.Replace(string(result), "\n", " ", -1)
	}
	if err != nil {
		return "", errors.Wrapf(err, "Encryption for plan %v failed %s", plan.Name, output)
	}

	return output, nil
}
