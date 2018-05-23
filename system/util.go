package system

import (
	"os"
)

var RedisDeployment = DeploymentWithNameAndFixture("redis-db-lock", "redis")
var JumpboxDeployment = DeploymentWithNameAndFixture("jumpbox-db-lock", "jumpbox")
var JumpboxInstance = JumpboxDeployment.Instance("jumpbox", "0")

func MustHaveEnv(keyname string) string {
	val := os.Getenv(keyname)

	if val == "" {
		panic("Need " + keyname + " for the test")
	}

	return val
}

func DeploymentWithNameAndFixture(deploymentName, fixtureName string) Deployment {
	return NewDeployment(deploymentName+"-"+MustHaveEnv("TEST_ENV"), "../../fixtures/"+fixtureName+".yml")
}
