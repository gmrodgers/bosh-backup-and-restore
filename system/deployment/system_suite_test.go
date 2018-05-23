package deployment

import (
	"fmt"
	"time"

	. "github.com/cloudfoundry-incubator/bosh-backup-and-restore/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"path"
	"sync"
	"testing"
)

func TestSystem(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "System Suite")
}

var (
	commandPath string
	err         error
)

var _ = BeforeSuite(func() {
	SetDefaultEventuallyTimeout(15 * time.Minute)

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer GinkgoRecover()
		defer wg.Done()

		By("deploying the Redis test release")
		RedisDeployment.Deploy()
	}()

	go func() {
		defer GinkgoRecover()
		defer wg.Done()

		By("deploying the jump box")
		JumpboxDeployment.Deploy()
	}()

	wg.Wait()

	By("building db-lock")
	commandPath, err = gexec.BuildWithEnvironment("github.com/cloudfoundry-incubator/bosh-backup-and-restore/cmd/bbr", []string{"GOOS=linux", "GOARCH=amd64"})
	Expect(err).NotTo(HaveOccurred())

	By("setting up the jump box")
	Eventually(JumpboxInstance.RunCommand(
		fmt.Sprintf("sudo mkdir %s && sudo chown vcap:vcap %s && sudo chmod 0777 %s", workspaceDir, workspaceDir, workspaceDir))).Should(gexec.Exit(0))

	JumpboxInstance.Copy(commandPath, path.Join(workspaceDir, "db-lock"))
	JumpboxInstance.Copy(MustHaveEnv("BOSH_CERT_PATH"), workspaceDir+"/bosh.crt")
})

var _ = AfterSuite(func() {
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer GinkgoRecover()
		defer wg.Done()

		By("tearing down the redis release")
		RedisDeployment.Delete()
	}()

	go func() {
		defer GinkgoRecover()
		defer wg.Done()

		By("tearing down the jump box")
		JumpboxDeployment.Delete()
	}()

	wg.Wait()
})

func runOnInstances(instanceCollection map[string][]string, f func(string, string)) {
	for instanceGroup, instances := range instanceCollection {
		for _, instanceIndex := range instances {
			f(instanceGroup, instanceIndex)
		}
	}
}
