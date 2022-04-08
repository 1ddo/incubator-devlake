package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type GithubPullRequestComment struct {
	GithubId        int `gorm:"primaryKey"`
	PullRequestId   int `gorm:"index"`
	Body            string
	AuthorUsername  string `gorm:"type:varchar(255)"`
	AuthorUserId    int
	GithubCreatedAt time.Time
	GithubUpdatedAt time.Time `gorm:"index"`
	common.NoPKModel
}
