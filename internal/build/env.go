
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:38</date>
//</624450091620503552>


package build

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
//These flags override values in build env.
	GitCommitFlag   = flag.String("git-commit", "", `Overrides git commit hash embedded into executables`)
	GitBranchFlag   = flag.String("git-branch", "", `Overrides git branch being built`)
	GitTagFlag      = flag.String("git-tag", "", `Overrides git tag being built`)
	BuildnumFlag    = flag.String("buildnum", "", `Overrides CI build number`)
	PullRequestFlag = flag.Bool("pull-request", false, `Overrides pull request status of the build`)
	CronJobFlag     = flag.Bool("cron-job", false, `Overrides cron job status of the build`)
)

//环境包含由生成环境提供的元数据。
type Environment struct {
Name                string //环境名称
Repo                string //name of GitHub repo
Commit, Branch, Tag string //GIT信息
	Buildnum            string
	IsPullRequest       bool
	IsCronJob           bool
}

func (env Environment) String() string {
	return fmt.Sprintf("%s env (commit:%s branch:%s tag:%s buildnum:%s pr:%t)",
		env.Name, env.Commit, env.Branch, env.Tag, env.Buildnum, env.IsPullRequest)
}

//env返回有关当前CI环境的元数据，返回到localenv
//如果不在CI上运行。
func Env() Environment {
	switch {
	case os.Getenv("CI") == "true" && os.Getenv("TRAVIS") == "true":
		return Environment{
			Name:          "travis",
			Repo:          os.Getenv("TRAVIS_REPO_SLUG"),
			Commit:        os.Getenv("TRAVIS_COMMIT"),
			Branch:        os.Getenv("TRAVIS_BRANCH"),
			Tag:           os.Getenv("TRAVIS_TAG"),
			Buildnum:      os.Getenv("TRAVIS_BUILD_NUMBER"),
			IsPullRequest: os.Getenv("TRAVIS_PULL_REQUEST") != "false",
			IsCronJob:     os.Getenv("TRAVIS_EVENT_TYPE") == "cron",
		}
	case os.Getenv("CI") == "True" && os.Getenv("APPVEYOR") == "True":
		return Environment{
			Name:          "appveyor",
			Repo:          os.Getenv("APPVEYOR_REPO_NAME"),
			Commit:        os.Getenv("APPVEYOR_REPO_COMMIT"),
			Branch:        os.Getenv("APPVEYOR_REPO_BRANCH"),
			Tag:           os.Getenv("APPVEYOR_REPO_TAG_NAME"),
			Buildnum:      os.Getenv("APPVEYOR_BUILD_NUMBER"),
			IsPullRequest: os.Getenv("APPVEYOR_PULL_REQUEST_NUMBER") != "",
			IsCronJob:     os.Getenv("APPVEYOR_SCHEDULED_BUILD") == "True",
		}
	default:
		return LocalEnv()
	}
}

//localenv返回从git收集的生成环境元数据。
func LocalEnv() Environment {
	env := applyEnvFlags(Environment{Name: "local", Repo: "ethereum/go-ethereum"})

	head := readGitFile("HEAD")
	if splits := strings.Split(head, " "); len(splits) == 2 {
		head = splits[1]
	} else {
		return env
	}
	if env.Commit == "" {
		env.Commit = readGitFile(head)
	}
	if env.Branch == "" {
		if head != "HEAD" {
			env.Branch = strings.TrimPrefix(head, "refs/heads/")
		}
	}
	if info, err := os.Stat(".git/objects"); err == nil && info.IsDir() && env.Tag == "" {
		env.Tag = firstLine(RunGit("tag", "-l", "--points-at", "HEAD"))
	}
	return env
}

func firstLine(s string) string {
	return strings.Split(s, "\n")[0]
}

func applyEnvFlags(env Environment) Environment {
	if !flag.Parsed() {
		panic("you need to call flag.Parse before Env or LocalEnv")
	}
	if *GitCommitFlag != "" {
		env.Commit = *GitCommitFlag
	}
	if *GitBranchFlag != "" {
		env.Branch = *GitBranchFlag
	}
	if *GitTagFlag != "" {
		env.Tag = *GitTagFlag
	}
	if *BuildnumFlag != "" {
		env.Buildnum = *BuildnumFlag
	}
	if *PullRequestFlag {
		env.IsPullRequest = true
	}
	if *CronJobFlag {
		env.IsCronJob = true
	}
	return env
}

