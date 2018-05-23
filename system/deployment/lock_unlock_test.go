package deployment

import (
	"fmt"

	. "github.com/cloudfoundry-incubator/bosh-backup-and-restore/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var workspaceDir = "/var/vcap/store/bbr-backup_workspace"
var instanceCollection = map[string][]string{
	"redis":       {"0", "1"},
	"other-redis": {"0"},
}

var _ = Describe("lock", func() {
	It("runs the pre-backup lock script", func() {
		lockCommand := fmt.Sprintf(
			`cd %s; BOSH_CLIENT_SECRET=%s ./db-lock deployment --ca-cert bosh.crt --username %s --target %s --deployment %s lock`,
			workspaceDir,
			MustHaveEnv("BOSH_CLIENT_SECRET"),
			MustHaveEnv("BOSH_CLIENT"),
			MustHaveEnv("BOSH_URL"),
			RedisDeployment.Name,
		)

		Eventually(JumpboxInstance.RunCommandAs("vcap", lockCommand)).Should(gexec.Exit(0))

		runOnInstances(instanceCollection, func(instName, instIndex string) {
			session := RedisDeployment.Instance(instName, instIndex).RunCommand(
				"cat /tmp/pre-backup-lock.out",
			)

			Eventually(session).Should(gexec.Exit(0))
			Expect(session.Out.Contents()).Should(ContainSubstring("output from pre-backup-lock"))
		})
	})
})

var _ = Describe("unlock", func() {
	It("runs the post-backup unlock script", func() {
		unlockCommand := fmt.Sprintf(
			`cd %s; BOSH_CLIENT_SECRET=%s ./db-lock deployment --ca-cert bosh.crt --username %s --target %s --deployment %s unlock`,
			workspaceDir,
			MustHaveEnv("BOSH_CLIENT_SECRET"),
			MustHaveEnv("BOSH_CLIENT"),
			MustHaveEnv("BOSH_URL"),
			RedisDeployment.Name,
		)

		Eventually(JumpboxInstance.RunCommandAs("vcap", unlockCommand)).Should(gexec.Exit(0))

		runOnInstances(instanceCollection, func(instName, instIndex string) {
			session := RedisDeployment.Instance(instName, instIndex).RunCommand(
				"cat /tmp/post-backup-unlock.out",
			)
			Eventually(session).Should(gexec.Exit(0))

			Expect(session.Out.Contents()).Should(ContainSubstring("output from post-backup-unlock"))
		})
	})
})
