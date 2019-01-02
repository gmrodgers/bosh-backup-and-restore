package director

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = FDescribe("PreBackupCheck", func() {
	It("checks if the director is backupable", func() {
		By("running the pre-backup-check command")
		preBackupCheckCommand := exec.Command(commandPath,
			"director",
			"--host", directorHost,
			"--username", directorSSHUsername,
			"--private-key-path", directorSSHPrivateKeyPath,
			"pre-backup-check",
		)

		session, err := gexec.Start(preBackupCheckCommand, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))
		Expect(session.Out.Contents()).To(ContainSubstring("Director can be backed up"))
	})
})
