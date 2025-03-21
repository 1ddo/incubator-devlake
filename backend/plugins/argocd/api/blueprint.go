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
	"bytes"
	"encoding/gob"
	"encoding/json"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/argocd/tasks"
)

func MakePipelinePlan(subtaskMetas []plugin.SubTaskMeta, connectionId uint64, scope []*common.ScopeConfig) (plugin.PipelinePlan, errors.Error) {
	var err error
	plan := make(plugin.PipelinePlan, len(scope))
	for i, scopeElem := range scope {
		taskOptions := make(map[string]interface{})
		buf := &bytes.Buffer{}
		gob.NewEncoder(buf).Encode(scopeElem.Entities)

		err = json.Unmarshal(buf.Bytes(), &taskOptions)
		if err != nil {
			return nil, errors.Default.Wrap(err, "error unmarshalling task options")
		}
		taskOptions["connectionId"] = connectionId

		//TODO Add transformation rules to task options

		/*
		   var transformationRules tasks.TransformationRules
		   if len(scopeElem.Transformation) > 0 {
		       err = json.Unmarshal(scopeElem.Transformation, &transformationRules)
		       if err != nil {
		           return nil, err
		       }
		   }
		*/
		//taskOptions["transformationRules"] = transformationRules
		_, err := tasks.DecodeAndValidateTaskOptions(taskOptions)
		if err != nil {
			return nil, err
		}
		// subtasks
		subtasks, err := helper.MakePipelinePlanSubtasks(subtaskMetas, scopeElem.Entities)
		if err != nil {
			return nil, err
		}
		plan[i] = plugin.PipelineStage{
			{
				Plugin:   "argocd",
				Subtasks: subtasks,
				Options:  taskOptions,
			},
		}
	}
	return plan, nil
}
