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

package tasks

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_COMMITTER_TABLE = "argocd_committer"

var _ plugin.SubTaskEntryPoint = CollectCommitter

func CollectCommitter(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMITTER_TABLE)
	//logger := taskCtx.GetLogger()

	collectorWithState, err := helper.NewStatefulApiCollector(*rawDataSubTaskArgs, data.CreatedDateAfter)
	if err != nil {
		return err
	}
	incremental := collectorWithState.IsIncremental()

	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		Incremental: incremental,
		ApiClient:   data.ApiClient,
		// PageSize:    100,
		// TODO write which api would you want request
		UrlTemplate: "https://people.apache.org/",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			input := reqData.Input.(*helper.DatePair)
			query.Set("start_time", strconv.FormatInt(input.PairStartTime.Unix(), 10))
			query.Set("end_time", strconv.FormatInt(input.PairEndTime.Unix(), 10))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			// TODO decode result from api request
			return []json.RawMessage{}, nil
		},
	})
	if err != nil {
		return err
	}
	return collectorWithState.Execute()
}

var CollectCommitterMeta = plugin.SubTaskMeta{
	Name:             "CollectCommitter",
	EntryPoint:       CollectCommitter,
	EnabledByDefault: true,
	Description:      "Collect Committer data from Argocd api",
	DomainTypes:      []string{},
}
