package director

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry-incubator/bosh-backup-and-restore/system"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Restores a director", func() {
	var restorePath = "/var/vcap/store/test-backup-and-restore"
	var restoredArtifactPath = restorePath + "/backup"
	var artifactName = "artifactToRestore"

	AfterEach(func() {
		By("cleaning up the jump box")
		Eventually(JumpboxInstance.RunCommandAs("vcap",
			fmt.Sprintf(
				`sudo rm -rf %s/%s`,
				workspaceDir,
				artifactName,
			),
		)).Should(gexec.Exit(0))

		By("cleaning up the director")
		Eventually(JumpboxInstance.RunCommandAs("vcap",
			fmt.Sprintf(`cd %s; ssh %s vcap@%s -i key.pem "sudo rm -rf %s"`,
				workspaceDir,
				skipSSHFingerprintCheckOpts,
				directorHost,
				restorePath,
			),
		)).Should(gexec.Exit(0))
	})

	It("restores", func() {
		By("setting up the jump box")
		Eventually(JumpboxInstance.RunCommandAs("vcap",
			fmt.Sprintf("sudo mkdir -p %s && sudo chmod -R 0777 %s",
				workspaceDir+"/"+artifactName, workspaceDir+"/"+artifactName))).Should(gexec.Exit(0))
		JumpboxInstance.Copy(fixturesPath+"bosh-0-amazing-backup-and-restore.tar", workspaceDir+"/"+artifactName)
		JumpboxInstance.Copy(fixturesPath+"bosh-0-remarkable-backup-and-restore.tar", workspaceDir+"/"+artifactName)
		JumpboxInstance.Copy(fixturesPath+"bosh-0-test-backup-and-restore.tar", workspaceDir+"/"+artifactName)
		JumpboxInstance.Copy(fixturesPath+"metadata", workspaceDir+"/"+artifactName)

		By("running the restore command")
		restoreCommand := JumpboxInstance.RunCommandAs("vcap",
			fmt.Sprintf(`cd %s; ./bbr director --username vcap --private-key-path ./key.pem --host %s restore --artifact-path %s`,
				workspaceDir,
				directorHost,
				artifactName,
			))
		Eventually(restoreCommand).Should(gexec.Exit(0))

		By("ensuring data is restored")
		Eventually(JumpboxInstance.RunCommandAs("vcap",
			fmt.Sprintf(`cd %s; ssh %s vcap@%s -i key.pem "stat %s"`,
				workspaceDir,
				skipSSHFingerprintCheckOpts,
				directorHost,
				restoredArtifactPath,
			),
		)).Should(gexec.Exit(0))
	})
})
