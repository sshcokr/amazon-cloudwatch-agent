package main

import (
	//"github.com/stretchr/testify/require"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
)

func TestAmiLatest(t *testing.T) {
	cfg, _ := config.LoadDefaultConfig(context.TODO())

	imng := CreateNewInstanceManager(cfg)
	// is it consistent
	previous := *imng.GetLatestAMIVersions().ImageId
	for i := 0; i < 5; i++ {
		current := *imng.GetLatestAMIVersions().ImageId
		require.Equalf(t, current, previous, "AMI is inconsistent %s | %s", current, previous)
	}
	fmt.Println(imng.GetLatestAMIVersions().ImageId)

}
func TestEc2Generation(t *testing.T) {
	rbm := CreateRemoteBuildManager(DEFAULT_INSTANCE_GUIDE)
	fmt.Println(rbm.ssmClient)
	rbm.Close()
}
func TestOnSpecificInstance(t *testing.T) {
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	imng := CreateNewInstanceManager(cfg)
	testInstance := GetInstanceFromID(imng.ec2Client, "i-0dd926b8dcf5884b6")
	ssmClient := ssm.NewFromConfig(cfg)
	RunCmdRemotely(ssmClient, testInstance, mergeCommands(
		"aws --version",
	),
		"Manual Testing")
}
func TestEnviorment(t *testing.T) {
	guide := map[string]OS{
		"MainBuildEnv": LINUX,
	}
	rbm := CreateRemoteBuildManager(guide)
	defer rbm.Close()
	func() {
		require.NoError(t,
			rbm.RunCommand(mergeCommands(
				"go version",
				"go env",
			),
				"MainBuildEnv",
				"test env go version"),
		)
	}()
	require.NoError(t,
		rbm.RunCommand(mergeCommands(
			"aws --version",
		),
			"MainBuildEnv",
			"test env aws"),
	)
	require.NoError(t,
		rbm.RunCommand(mergeCommands(
			"make --version",
		),
			"MainBuildEnv",

			"make"),
	)
}
func TestPublicRepoBuild(t *testing.T) {
	REPO_NAME := "https://github.com/aws/amazon-cloudwatch-agent.git"
	BRANCH_NAME := "main"
	rbm := CreateRemoteBuildManager(DEFAULT_INSTANCE_GUIDE)
	defer rbm.Close()
	rbm.BuildCWAAgent(REPO_NAME, BRANCH_NAME, fmt.Sprintf("PUBLIC_REPO_TEST-%d", time.Now().Unix()), "MainBuildEnv")
	//rbm.RunCommand(RemoveFolder("ccwa"), "clean the repo folder")
}

func TestPrivateFork(t *testing.T) {
	//REPO_NAME := "https://github.com/aws/amazon-cloudwatch-agent.git"
	//BRANCH_NAME := "main"
	//rbm := CreateRemoteBuildManager()
	//rbm.CloneGitRepo(REPO_NAME, BRANCH_NAME)
}
func TestMakeMsiZip(t *testing.T) {
	//TestPublicRepoBuild(t)
	guide := map[string]OS{
		"WindowsMSIPacker": LINUX,
	}
	rbm := CreateRemoteBuildManager(guide)
	defer rbm.Close()
	require.NoError(t, rbm.MakeMsiZip("WindowsMSIPacker", "PUBLIC_REPO_TEST-1695063642"))
}
func TestBuildMsi(t *testing.T) {

}
func TestMakeMacPkg(t *testing.T) {

}
