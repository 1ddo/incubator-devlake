/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/trello/models"
	"github.com/go-playground/validator/v10"
)

var vld *validator.Validate
var connectionHelper *api.ConnectionApiHelper
var scopeHelper *api.ScopeApiHelper[models.TrelloConnection, models.TrelloBoard, models.TrelloScopeConfig]
var basicRes context.BasicRes
var scHelper *api.ScopeConfigHelper[models.TrelloScopeConfig]

func Init(br context.BasicRes) {
	basicRes = br
	vld = validator.New()
	connectionHelper = api.NewConnectionHelper(
		basicRes,
		vld,
	)
	params := &api.ReflectionParameters{
		ScopeIdFieldName:  "BoardId",
		ScopeIdColumnName: "board_id",
	}
	scopeHelper = api.NewScopeHelper[models.TrelloConnection, models.TrelloBoard, models.TrelloScopeConfig](
		basicRes,
		vld,
		connectionHelper,
		api.NewScopeDatabaseHelperImpl[models.TrelloConnection, models.TrelloBoard, models.TrelloScopeConfig](
			basicRes, connectionHelper, params),
		params,
		nil,
	)
	scHelper = api.NewScopeConfigHelper[models.TrelloScopeConfig](
		basicRes,
		vld,
	)
}
