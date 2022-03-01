package tasks

import (
	"context"
	"fmt"
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/errors"
	lakeModels "github.com/merico-dev/lake/models"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
	"regexp"
	"strconv"
	"strings"
)

var prBodyCloseRegex *regexp.Regexp
var prBodyClosePattern string
var numberPrefix string

func init() {
	prBodyClosePattern = config.GetConfig().GetString("GITHUB_PR_BODY_CLOSE_PATTERN")
	numberPrefix = config.GetConfig().GetString("GITHUB_PR_BODY_NUMBER_PREFIX")
}

func EnrichPullRequestIssues(ctx context.Context, repoId int, owner string, repo string) (err error) {
	numberPattern := fmt.Sprintf(numberPrefix+`\d+[ ]*)+)`, owner, repo)
	if len(prBodyClosePattern) > 0 {
		prPattern := prBodyClosePattern + numberPattern
		prBodyCloseRegex = regexp.MustCompile(prPattern)
	}

	githubIssueUrlPattern := fmt.Sprintf(config.GetConfig().GetString("GITHUB_ISSUE_URL_PATTERN"), owner, repo)
	githubPullRequst := &githubModels.GithubPullRequest{}
	cursor, err := lakeModels.Db.Model(&githubPullRequst).
		Where("repo_id = ?", repoId).
		Rows()
	if err != nil {
		return err
	}
	resList := make([]string, 0)

	defer cursor.Close()
	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return errors.TaskCanceled
		default:
		}
		err = lakeModels.Db.ScanRows(cursor, githubPullRequst)
		if err != nil {
			return err
		}

		issueNumberListStr := getCloseIssueId(githubPullRequst.Body)

		if issueNumberListStr == "" {
			continue
		}
		//replace https:// to #, then we can deal with later
		if strings.Contains(issueNumberListStr, "https") {
			if !strings.Contains(issueNumberListStr, githubIssueUrlPattern) {
				continue
			}
			numberPrefixRegex := regexp.MustCompile(numberPrefix)
			issueNumberListStr = numberPrefixRegex.ReplaceAllString(issueNumberListStr, "#")
		}
		charPattern := regexp.MustCompile(`[a-zA-Z\s,]+`)
		issueNumberListStr = charPattern.ReplaceAllString(issueNumberListStr, "#")
		//split the string by '#'
		issueNumberList := strings.Split(issueNumberListStr, "#")

		for _, issueNumberStr := range issueNumberList {
			issue := &githubModels.GithubIssue{}

			issueNumberStr = strings.TrimSpace(issueNumberStr)
			issueNumber, numFormatErr := strconv.Atoi(issueNumberStr)
			if numFormatErr != nil {
				continue
			}
			err = lakeModels.Db.Where("number = ? and repo_id = ?", issueNumber, repoId).Limit(1).Find(issue).Error
			if err != nil {
				return err
			}
			if issue == nil {
				continue
			}
			githubPullRequstIssue := &githubModels.GithubPullRequestIssue{
				PullRequestId: githubPullRequst.GithubId,
				IssueId:       issue.GithubId,
				PullNumber:    githubPullRequst.Number,
				IssueNumber:   issue.Number,
			}

			err = lakeModels.Db.Clauses(
				clause.OnConflict{UpdateAll: true}).Create(githubPullRequstIssue).Error
			if err != nil {
				return err
			}
		}
	}
	for _, v := range resList {
		fmt.Println(v)
	}

	return nil
}

func getCloseIssueId(body string) string {
	if prBodyCloseRegex != nil {
		matchString := prBodyCloseRegex.FindString(body)
		return matchString
	}
	return ""
}