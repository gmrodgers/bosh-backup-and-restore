package deployment

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/testcluster"
	"github.com/pivotal-cf-experimental/cf-webmock/mockbosh"
	"github.com/pivotal-cf-experimental/cf-webmock/mockhttp"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lock", func() {
	var director *mockhttp.Server
	var backupWorkspace string
	var session *gexec.Session
	var stdin io.WriteCloser
	var deploymentName string
	var downloadManifest bool
	var waitForBackupToFinish bool
	var verifyMocks bool
	var instance1 *testcluster.Instance
	manifest := `---
instance_groups:
- name: redis-dedicated-node
  instances: 1
  jobs:
  - name: redis
    release: redis
  - name: redis-writer
    release: redis
  - name: redis-broker
    release: redis
- name: redis-broker
  instances: 1
  jobs:
  - name: redis
    release: redis
  - name: redis-writer
    release: redis
  - name: redis-broker
    release: redis
`

	BeforeEach(func() {
		deploymentName = "my-little-deployment"
		downloadManifest = false
		waitForBackupToFinish = true
		verifyMocks = true
		director = mockbosh.NewTLS()
		director.ExpectedBasicAuth("admin", "admin")
		var err error
		backupWorkspace, err = ioutil.TempDir(".", "backup-workspace-")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if verifyMocks {
			director.VerifyMocks()
		}
		director.Close()

		instance1.DieInBackground()
		Expect(os.RemoveAll(backupWorkspace)).To(Succeed())
	})

	JustBeforeEach(func() {
		env := []string{"BOSH_CLIENT_SECRET=admin"}

		params := []string{
			"deployment",
			"--ca-cert", sslCertPath,
			"--username", "admin",
			"--target", director.URL,
			"--deployment", deploymentName,
			"--debug",
			"lock"}

		if waitForBackupToFinish {
			session = binary.Run(backupWorkspace, env, params...)
		} else {
			session, stdin = binary.Start(backupWorkspace, env, params...)
			Eventually(session).Should(gbytes.Say(".+"))
		}
	})

	Context("When there is a deployment which has one instance", func() {
		singleInstanceResponse := func(instanceGroupName string) []mockbosh.VMsOutput {
			return []mockbosh.VMsOutput{
				{
					IPs:     []string{"10.0.0.1"},
					JobName: instanceGroupName,
				},
			}
		}

		Context("and there is a plausible backup script", func() {
			BeforeEach(func() {
				instance1 = testcluster.NewInstance()
				By("creating a dummy backup script")
				instance1.CreateScript("/var/vcap/jobs/redis/bin/bbr/backup", `#!/usr/bin/env sh
set -u
touch /tmp/backup-script-was-run
printf "backupcontent1" > $BBR_ARTIFACT_DIRECTORY/backupdump1
printf "backupcontent2" > $BBR_ARTIFACT_DIRECTORY/backupdump2
`)
			})

			Context("and we don't ask for the manifest to be downloaded", func() {
				BeforeEach(func() {
					MockDirectorWithoutCleanupWith(director,
						mockbosh.Info().WithAuthTypeBasic(),
						VmsForDeployment(deploymentName, singleInstanceResponse("redis-dedicated-node")),
						DownloadManifest(deploymentName, manifest),
						SetupSSH(deploymentName, "redis-dedicated-node", "fake-uuid", 0, instance1))
				})

				Context("and the pre-backup-lock script is present", func() {
					BeforeEach(func() {
						instance1.CreateScript("/var/vcap/jobs/redis/bin/bbr/pre-backup-lock", `#!/usr/bin/env sh
touch /tmp/pre-backup-lock-script-was-run
`)
						instance1.CreateScript("/var/vcap/jobs/redis-broker/bin/bbr/pre-backup-lock", ``)
					})

					It("executes and logs the locks", func() {
						By("running the pre-backup-lock script", func() {
							Expect(instance1.FileExists("/tmp/pre-backup-lock-script-was-run")).To(BeTrue())
						})

						By("logging that it is locking the instance, and listing the scripts", func() {
							assertOutput(session, []string{
								`Locking redis on redis-dedicated-node/fake-uuid for backup`,
								"> /var/vcap/jobs/redis/bin/bbr/pre-backup-lock",
								"> /var/vcap/jobs/redis-broker/bin/bbr/pre-backup-lock",
							})
						})
					})

				})

				Context("when the pre-backup-lock script fails", func() {
					BeforeEach(func() {
						instance1.CreateScript("/var/vcap/jobs/redis/bin/bbr/pre-backup-lock", `#!/usr/bin/env sh
echo 'ultra-bar'
(>&2 echo 'ultra-baz')
touch /tmp/pre-backup-lock-output
exit 1
`)
						instance1.CreateScript("/var/vcap/jobs/redis-broker/bin/bbr/pre-backup-lock", ``)
						instance1.CreateScript("/var/vcap/jobs/redis/bin/bbr/post-backup-unlock", `#!/usr/bin/env sh
touch /tmp/post-backup-unlock-output
`)
					})

					It("logs the failure, and unlocks the system", func() {
						By("running the pre-backup-lock scripts", func() {
							Expect(instance1.FileExists("/tmp/pre-backup-lock-output")).To(BeTrue())
						})

						By("not running the backup script", func() {
							Expect(instance1.FileExists("/tmp/backup-script-was-run")).NotTo(BeTrue())
						})

						By("exiting with a non-zero error code", func() {
							Expect(session.ExitCode()).NotTo(Equal(0))
						})

						By("logging the error", func() {
							Expect(session.Err.Contents()).To(ContainSubstring(
								"Error attempting to run pre-backup-lock for job redis on redis-dedicated-node/fake-uuid"))
						})

						By("logging stderr", func() {
							Expect(session.Err.Contents()).To(ContainSubstring("ultra-baz"))
						})

						By("also running the post-backup-unlock scripts", func() {
							Expect(instance1.FileExists("/tmp/post-backup-unlock-output")).To(BeTrue())
						})
					})
				})

				Context("but /var/vcap/store is not world-accessible", func() {
					BeforeEach(func() {
						instance1.Run("sudo", "chmod", "700", "/var/vcap/store")
					})

					It("successfully backs up the deployment", func() {
						Expect(session.ExitCode()).To(BeZero())
					})
				})
			})
		})
	})

	Context("When there is a deployment which has two instances", func() {
		twoInstancesResponse := func(firstInstanceGroupName, secondInstanceGroupName string) []mockbosh.VMsOutput {

			return []mockbosh.VMsOutput{
				{
					IPs:     []string{"10.0.0.1"},
					JobName: firstInstanceGroupName,
				},
				{
					IPs:     []string{"10.0.0.2"},
					JobName: secondInstanceGroupName,
				},
			}
		}

		Context("one backupable", func() {
			var firstReturnedInstance, secondReturnedInstance *testcluster.Instance

			BeforeEach(func() {
				deploymentName = "my-bigger-deployment"
				firstReturnedInstance = testcluster.NewInstance()
				secondReturnedInstance = testcluster.NewInstance()
				MockDirectorWithoutCleanupWith(director,
					mockbosh.Info().WithAuthTypeBasic(),
					VmsForDeployment(deploymentName, twoInstancesResponse("redis-dedicated-node", "redis-broker")),
					DownloadManifest(deploymentName, manifest),
					append(SetupSSH(deploymentName, "redis-dedicated-node", "fake-uuid", 0, firstReturnedInstance),
						SetupSSH(deploymentName, "redis-broker", "fake-uuid-2", 0, secondReturnedInstance)...),
				)
				firstReturnedInstance.CreateExecutableFiles(
					"/var/vcap/jobs/redis/bin/bbr/backup",
				)
			})

			AfterEach(func() {
				firstReturnedInstance.DieInBackground()
				secondReturnedInstance.DieInBackground()
			})

			It("succeeds", func() {
				Expect(session.ExitCode()).To(BeZero())
			})

			Context("with ordering on pre-backup-lock specified", func() {
				BeforeEach(func() {
					firstReturnedInstance.CreateScript(
						"/var/vcap/jobs/redis/bin/bbr/pre-backup-lock", `#!/usr/bin/env sh
touch /tmp/redis-pre-backup-lock-called
exit 0`)
					secondReturnedInstance.CreateScript(
						"/var/vcap/jobs/redis-writer/bin/bbr/pre-backup-lock", `#!/usr/bin/env sh
touch /tmp/redis-writer-pre-backup-lock-called
exit 0`)
					secondReturnedInstance.CreateScript("/var/vcap/jobs/redis-writer/bin/bbr/metadata",
						`#!/usr/bin/env sh
echo "---
backup_should_be_locked_before:
- job_name: redis
  release: redis
"`)
				})

				It("locks in the specified order", func() {
					redisLockTime := firstReturnedInstance.GetCreatedTime("/tmp/redis-pre-backup-lock-called")
					redisWriterLockTime := secondReturnedInstance.GetCreatedTime("/tmp/redis-writer-pre-backup-lock-called")

					Expect(string(session.Out.Contents())).To(ContainSubstring("Detected order: redis-writer should be locked before redis/redis during backup"))

					Expect(redisWriterLockTime < redisLockTime).To(BeTrue(), fmt.Sprintf(
						"Writer locked at %s, which is after the server locked (%s)",
						strings.TrimSuffix(redisWriterLockTime, "\n"),
						strings.TrimSuffix(redisLockTime, "\n")))

				})
			})

			Context("with ordering on pre-backup-lock (where the default ordering would unlock in the wrong order) and a lock lock fails",
				func() {
					BeforeEach(func() {
						secondReturnedInstance.CreateScript(
							"/var/vcap/jobs/redis/bin/bbr/pre-backup-lock", `#!/usr/bin/env sh
touch /tmp/redis-pre-backup-lock-called
exit 0`)
						firstReturnedInstance.CreateScript(
							"/var/vcap/jobs/redis-writer/bin/bbr/pre-backup-lock", `#!/usr/bin/env sh
touch /tmp/redis-writer-pre-backup-lock-called
exit 1`)
						secondReturnedInstance.CreateScript(
							"/var/vcap/jobs/redis/bin/bbr/post-backup-unlock", `#!/usr/bin/env sh
touch /tmp/redis-post-backup-unlock-called
exit 0`)
						firstReturnedInstance.CreateScript(
							"/var/vcap/jobs/redis-writer/bin/bbr/post-backup-unlock", `#!/usr/bin/env sh
touch /tmp/redis-writer-post-backup-unlock-called
exit 0`)
						firstReturnedInstance.CreateScript("/var/vcap/jobs/redis-writer/bin/bbr/metadata",
							`#!/usr/bin/env sh
echo "---
backup_should_be_locked_before:
- job_name: redis
  release: redis
"`)
					})

					It("unlocks in the right order", func() {
						By("unlocking the redis job before unlocking the redis-writer job")
						redisUnlockTime := secondReturnedInstance.GetCreatedTime("/tmp/redis-post-backup-unlock-called")
						redisWriterUnlockTime := firstReturnedInstance.GetCreatedTime("/tmp/redis-writer-post-backup-unlock-called")

						Expect(redisUnlockTime < redisWriterUnlockTime).To(BeTrue(), fmt.Sprintf(
							"Writer unlocked at %s, which is before the server unlocked (%s)",
							strings.TrimSuffix(redisWriterUnlockTime, "\n"),
							strings.TrimSuffix(redisUnlockTime, "\n")))
					})
				})

			Context("but the pre-backup-lock ordering is cyclic", func() {
				BeforeEach(func() {
					firstReturnedInstance.CreateScript(
						"/var/vcap/jobs/redis/bin/bbr/pre-backup-lock", `#!/usr/bin/env sh
touch /tmp/redis-pre-backup-lock-called
exit 0`)
					firstReturnedInstance.CreateScript(
						"/var/vcap/jobs/redis-writer/bin/bbr/pre-backup-lock", `#!/usr/bin/env sh
touch /tmp/redis-writer-pre-backup-lock-called
exit 0`)
					firstReturnedInstance.CreateScript("/var/vcap/jobs/redis-writer/bin/bbr/metadata",
						`#!/usr/bin/env sh
echo "---
backup_should_be_locked_before:
- job_name: redis
  release: redis
"`)
					firstReturnedInstance.CreateScript("/var/vcap/jobs/redis/bin/bbr/metadata",
						`#!/usr/bin/env sh
echo "---
backup_should_be_locked_before:
- job_name: redis-writer
  release: redis
"`)
				})

				It("Should fail", func() {
					By("exiting with an error", func() {
						Expect(session).To(gexec.Exit(1))
					})

					By("printing a helpful error message", func() {
						Expect(string(session.Err.Contents())).To(ContainSubstring("job locking dependency graph is cyclic"))
					})
				})
			})
		})
	})

	Context("When deployment does not exist", func() {
		BeforeEach(func() {
			deploymentName = "my-non-existent-deployment"
			director.VerifyAndMock(
				mockbosh.Info().WithAuthTypeBasic(),
				mockbosh.VMsForDeployment(deploymentName).NotFound(),
			)
		})

		It("errors and exits", func() {
			By("returning exit code 1", func() {
				Expect(session.ExitCode()).To(Equal(1))
			})

			By("printing an error", func() {
				Expect(string(session.Err.Contents())).To(ContainSubstring("Director responded with non-successful status code"))
			})

			By("not printing a recommendation to run bbr backup-cleanup", func() {
				Expect(string(session.Err.Contents())).NotTo(ContainSubstring("It is recommended that you run `bbr backup-cleanup`"))
			})
		})

	})
})
