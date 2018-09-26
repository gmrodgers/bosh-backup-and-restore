package all_deployments_tests

import (
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	. "github.com/cloudfoundry-incubator/bosh-backup-and-restore/system"
)

var _ = Describe("All deployments", func() {
	It("Can run pre-backup-check on all deployments", func() {
		cmd := exec.Command(
			commandPath,
			"deployment",
			"--ca-cert", MustHaveEnv("BOSH_CA_CERT"),
			"--username", MustHaveEnv("BOSH_CLIENT"),
			"--password", MustHaveEnv("BOSH_CLIENT_SECRET"),
			"--target", MustHaveEnv("BOSH_ENVIRONMENT"),
			"--all-deployments",
			"pre-backup-check",
		)
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		Expect(session.Out).To(gbytes.Say("Deployment 'redis-1' can be backed up."))
		Expect(session.Out).To(gbytes.Say("Deployment 'redis-2' can be backed up."))
		Expect(session.Out).To(gbytes.Say("Deployment 'redis-3' can be backed up."))
		Expect(session.Out).To(gbytes.Say("All 3 deployments can be backed up"))
	})

	FIt("Can run backup on all deployments", func() {
		cmd := exec.Command(
			commandPath,
			"deployment",
			"--ca-cert", MustHaveEnv("BOSH_CA_CERT"),
			"--username", MustHaveEnv("BOSH_CLIENT"),
			"--password", MustHaveEnv("BOSH_CLIENT_SECRET"),
			"--target", MustHaveEnv("BOSH_ENVIRONMENT"),
			"--all-deployments",
			"backup",
		)
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		redisDeployment1 := NewDeployment("redis-1", "")
		redisInstance1 := redisDeployment1.Instance("redis", "0")
		//redisDeployment2 := NewDeployment("redis-2", "")
		//redisDeployment3 := NewDeployment("redis-3", "")

		By("providing debug output", func() {
			Expect(session.Out).To(gbytes.Say("Starting backup of redis-1"))
			Expect(session.Out).To(gbytes.Say("Backup created of redis-1"))
			Expect(session.Out).To(gbytes.Say("Starting backup of redis-2"))
			Expect(session.Out).To(gbytes.Say("Backup created of redis-2"))
			Expect(session.Out).To(gbytes.Say("Starting backup of redis-3"))
			Expect(session.Out).To(gbytes.Say("Backup created of redis-3"))
			Expect(session.Out).To(gbytes.Say("All 3 deployments backed up."))
		})

		// TODO do we need this? If yes, move it before backup
		By("populating data in redis", func() {
			dataFixture := "../../fixtures/redis_test_commands"
			redisInstance1.Copy(dataFixture, "/tmp")
			Eventually(
				redisInstance1.RunCommand(
					"cat /tmp/redis_test_commands | /var/vcap/packages/redis/bin/redis-cli > /dev/null",
				),
			).Should(gexec.Exit(0))
		})

		By("running the pre-backup lock script", func() {
			session := redisInstance1.RunCommand(
				"cat /tmp/pre-backup-lock.out",
			)

			Eventually(session).Should(gexec.Exit(0))
			Expect(session.Out).To(gbytes.Say("output from pre-backup-lock"))
		})

		By("running the post backup unlock script", func() {
			session := redisInstance1.RunCommand(
				"cat /tmp/post-backup-unlock.out",
			)
			Eventually(session).Should(gexec.Exit(0))

			Expect(session.Out).To(gbytes.Say("output from post-backup-unlock"))
		})

		By("creating a timestamped directory for holding the artifacts locally", func() {
			cmd := exec.Command("ls", ".")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))
			Expect(session.Out).To(gbytes.Say(`\b` + redisDeployment1.Name + `_(\d){8}T(\d){6}Z\b`))
		})

		By("creating the backup artifacts locally", func() {
			// TODO use glob for this
			cmd := exec.Command("stat", fmt.Sprintf("%s/redis-0-redis-server.tar", BackupDirWithTimestamp(redisDeployment1.Name)))
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))
		})

		By("cleaning up artifacts from the remote instances", func() {

			session := redisInstance1.RunCommand(
				"ls -l /var/vcap/store/bbr-backup",
			)
			Eventually(session).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(1))
			Expect(session.Out).To(gbytes.Say("No such file or directory"))
		})
	})
})
