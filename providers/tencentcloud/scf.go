// Copyright 2022 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tencentcloud

import (
	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
)

type ScfGenerator struct {
	TencentCloudService
}

func (g *ScfGenerator) InitResources() error {
	args := g.GetArgs()
	region := args["region"].(string)
	credential := args["credential"].(common.Credential)
	profile := NewTencentCloudClientProfile()
	client, err := scf.NewClient(&credential, region, profile)
	if err != nil {
		return err
	}

	request := scf.NewListFunctionsRequest()

	var offset int64
	var pageSize int64 = 50
	allInstances := make([]*scf.Function, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		response, err := client.ListFunctions(request)
		if err != nil {
			return err
		}

		allInstances = append(allInstances, response.Response.Functions...)
		if len(response.Response.Functions) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	for _, instance := range allInstances {
		resource := terraformutils.NewResource(
			*instance.Namespace+"+"+*instance.FunctionName,
			*instance.Namespace+"_"+*instance.FunctionName,
			"tencentcloud_scf_function",
			"tencentcloud",
			map[string]string{},
			[]string{},
			map[string]interface{}{},
		)
		g.Resources = append(g.Resources, resource)
	}

	return nil
}
