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

package impl

import (
	"fmt"
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/argocd/api"
	"github.com/apache/incubator-devlake/plugins/argocd/models"
	"github.com/apache/incubator-devlake/plugins/argocd/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/argocd/tasks"
	"github.com/sirupsen/logrus"
)

// make sure interface is implemented
var _ plugin.PluginMeta = (*Argocd)(nil)
var _ plugin.PluginInit = (*Argocd)(nil)
var _ plugin.PluginTask = (*Argocd)(nil)
var _ plugin.PluginApi = (*Argocd)(nil)
var _ plugin.PluginMigration = (*Argocd)(nil)
var _ plugin.CloseablePluginTask = (*Argocd)(nil)

// var _ plugin.BlPluginBlueprintV100 = (*Argocd)(nil)
var _ plugin.CloseablePluginTask = (*Argocd)(nil)

type Argocd struct{}

func (p Argocd) Description() string {
	return "collect some Argocd data"
}

func (p Argocd) Init(br context.BasicRes) errors.Error {
	api.Init(br)
	return nil
}

func (p Argocd) SubTaskMetas() []plugin.SubTaskMeta {
	// TODO add your sub task here
	return []plugin.SubTaskMeta{}
}

func (p Argocd) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, err
	}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	connection := &models.ArgocdConnection{}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Argocd connection by the given connection ID")
	}

	apiClient, err := tasks.NewArgocdApiClient(taskCtx, connection)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Argocd API client instance")
	}
	taskData := &tasks.ArgocdTaskData{
		Options:   op,
		ApiClient: apiClient,
	}
	var createdDateAfter time.Time
	if op.CreatedDateAfter != "" {
		createdDateAfter, err = errors.Convert01(time.Parse(time.RFC3339, op.CreatedDateAfter))
		if err != nil {
			return nil, errors.BadInput.Wrap(err, "invalid value for `createdDateAfter`")
		}
	}
	if !createdDateAfter.IsZero() {
		taskData.CreatedDateAfter = &createdDateAfter
		logrus.Debug("collect data updated createdDateAfter %s", createdDateAfter)
	}
	return taskData, nil
}

// PkgPath information lost when compiled as plugin(.so)
func (p Argocd) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/argocd"
}

func (p Argocd) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Argocd) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"connections": {
			"POST": api.PostConnections,
			"GET":  api.ListConnections,
		},
		"connections/:connectionId": {
			"GET":    api.GetConnection,
			"PATCH":  api.PatchConnection,
			"DELETE": api.DeleteConnection,
		},
	}
}

func (p Argocd) MakePipelinePlan(connectionId uint64, scope []*common.ScopeConfig, subtasks []plugin.SubTaskMeta) (plugin.PipelinePlan, errors.Error) {
	return api.MakePipelinePlan(subtasks, connectionId, scope)
}

func (p Argocd) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.ArgocdTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
