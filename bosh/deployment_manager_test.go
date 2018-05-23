package bosh_test

import (
	"fmt"

	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/bosh"
	"github.com/cloudfoundry-incubator/bosh-backup-and-restore/bosh/fakes"
	orchestrator_fakes "github.com/cloudfoundry-incubator/bosh-backup-and-restore/orchestrator/fakes"
)

var _ = Describe("DeploymentManager", func() {
	var boshClient *fakes.FakeBoshClient
	var logger *fakes.FakeLogger
	var deploymentName = "brownie"

	var deploymentManager *bosh.DeploymentManager
	BeforeEach(func() {
		boshClient = new(fakes.FakeBoshClient)
		logger = new(fakes.FakeLogger)
	})
	JustBeforeEach(func() {
		deploymentManager = bosh.NewDeploymentManager(boshClient, logger)
	})

	Context("Find", func() {
		var findError error
		var deployment orchestrator.Deployment
		var instances []orchestrator.Instance
		BeforeEach(func() {
			instances = []orchestrator.Instance{new(orchestrator_fakes.FakeInstance)}
			boshClient.FindInstancesReturns(instances, nil)
		})
		JustBeforeEach(func() {
			deployment, findError = deploymentManager.Find(deploymentName)
		})
		It("asks the bosh director for instances", func() {
			Expect(boshClient.FindInstancesCallCount()).To(Equal(1))
			Expect(boshClient.FindInstancesArgsForCall(0)).To(Equal(deploymentName))
		})
		It("returns the deployment manager with instances", func() {
			Expect(deployment).To(Equal(orchestrator.NewDeployment(logger, instances)))
		})

		Context("error finding instances", func() {
			var expectedFindError = fmt.Errorf("a tuna sandwich")
			BeforeEach(func() {
				boshClient.FindInstancesReturns(nil, expectedFindError)
			})

			It("returns an error", func() {
				Expect(findError).To(MatchError(ContainSubstring("failed to find instances")))
			})
		})
	})
})
