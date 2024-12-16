/*
*  Copyright (c) WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 LLC. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package impl

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/wso2/product-mi-tooling/cmd/formatter"
	"github.com/wso2/product-mi-tooling/cmd/utils"
	"github.com/wso2/product-mi-tooling/cmd/utils/artifactUtils"
)

const (
	defaultInboundEndpointListTableFormat = "table {{.Name}}\t{{.Type}}\t{{.Status}}"
	defaultInboundEndpointDetailedFormat  = "detail Name - {{.Name}}\n" +
		"Type - {{.Type}}\n" +
		"Stats - {{.Stats}}\n" +
		"Tracing - {{.Tracing}}\n" +
		"Status - {{.Status}}\n" +
		"Parameters :\n" +
		"NAME\tVALUE\n" +
		"{{range .Parameters}}{{.Name}}\t{{.Value}}\n{{end}}"
)

// GetInboundEndpointList returns a list of inbound endpoints deployed in the micro integrator in a given environment
func GetInboundEndpointList(env string) (*artifactUtils.InboundEndpointList, error) {
	resp, err := getArtifactList(utils.MiManagementInboundEndpointResource, env, &artifactUtils.InboundEndpointList{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactUtils.InboundEndpointList), nil
}

// PrintInboundEndpointList print a list of inbound endpoints according to the given format
func PrintInboundEndpointList(inboundEPList *artifactUtils.InboundEndpointList, format string) {
	if inboundEPList.Count > 0 {
		inboundEPs := inboundEPList.InboundEndpoints
		inboundEPListContext := getContextWithFormat(format, defaultInboundEndpointListTableFormat)

		renderer := func(w io.Writer, t *template.Template) error {
			for _, inboundEP := range inboundEPs {
				if err := t.Execute(w, inboundEP); err != nil {
					return err
				}
				_, _ = w.Write([]byte{'\n'})
			}
			return nil
		}
		inboundEPListTableHeaders := map[string]string{
			"Name": nameHeader,
			"Type": typeHeader,
			"Status": statusHeader,
		}
		if err := inboundEPListContext.Write(renderer, inboundEPListTableHeaders); err != nil {
			fmt.Println("Error executing template:", err.Error())
		}
	} else {
		fmt.Println("No Inbound Endpoints found")
	}
}

// GetInboundEndpoint returns a information about a specific inbound endpoint deployed in the micro integrator in a given environment
func GetInboundEndpoint(env, inboundEPName string) (*artifactUtils.InboundEndpoint, error) {
	resp, err := getArtifactInfo(utils.MiManagementInboundEndpointResource, "inboundEndpointName", inboundEPName, env,
		&artifactUtils.InboundEndpoint{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactUtils.InboundEndpoint), nil
}

// PrintInboundEndpointDetails prints details about an inbound endpoint according to the given format
func PrintInboundEndpointDetails(inboundEP *artifactUtils.InboundEndpoint, format string) {
	if format == "" || strings.HasPrefix(format, formatter.TableFormatKey) {
		format = defaultInboundEndpointDetailedFormat
	}

	inboundEPContext := formatter.NewContext(os.Stdout, format)
	renderer := getItemRenderer(inboundEP)

	if err := inboundEPContext.Write(renderer, nil); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}
