package director

import (
	"io/ioutil"
	"os"

	. "github.com/cloudfoundry-incubator/bosh-backup-and-restore/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
	"time"
)

const (
	fixturesPath                = "../../fixtures/director-backup/"
	skipSSHFingerprintCheckOpts = "-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null"
)

var (
	workspaceDir              string
	commandPath               string
	directorHost              string
	directorSSHUsername       string
	directorSSHPrivateKeyPath string
	err                       error
)

func TestDirector(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Director Suite")
}

var _ = BeforeSuite(func() {
	SetDefaultEventuallyTimeout(4 * time.Minute)

	MustHaveEnv("BOSH_ALL_PROXY")

	directorHost = MustHaveEnv("DIRECTOR_HOST")
	directorSSHUsername = MustHaveEnv("DIRECTOR_SSH_USERNAME")
	directorSSHPrivateKeyPath = MustHaveEnv("DIRECTOR_SSH_KEY_PATH")

	commandPath, err = gexec.Build("github.com/cloudfoundry-incubator/bosh-backup-and-restore/cmd/bbr")
	Expect(err).NotTo(HaveOccurred())

	workspaceDir, err = ioutil.TempDir("", "bbr_system_test_director")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
	Expect(os.RemoveAll(workspaceDir)).To(Succeed())
})
