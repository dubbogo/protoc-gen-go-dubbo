/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package generator

import (
	"fmt"

	"github.com/dubbogo/protoc-gen-go-dubbo/proto/unified_idl_extend"
	"github.com/dubbogo/protoc-gen-go-dubbo/util"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

var (
	ErrStreamMethod               = errors.New("dubbo doesn't support stream method")
	ErrMoreExtendArgsRespFieldNum = errors.New("extend args for response message should only has 1 field")
	ErrNoExtendArgsRespFieldNum   = errors.New("extend args for response message should has a field")
)

type Dubbogo struct {
	*protogen.File

	Source       string
	ProtoPackage string
	Services     []*Service
}

type Service struct {
	ServiceName   string
	InterfaceName string
	Methods       []*Method
}

type Method struct {
	MethodName string
	InvokeName string

	// empty when RequestExtendArgs is true
	RequestType       string
	RequestExtendArgs bool
	ArgsType          []string
	ArgsName          []string

	ResponseExtendArgs bool
	ResponseIsWrapper  bool
	ReturnType         string
}

func ProcessProtoFile(g *protogen.GeneratedFile, file *protogen.File, methodProtocolSpecFlag, serviceProtocolSpecFlag bool) (*Dubbogo, error) {
	desc := file.Proto
	dubboGo := &Dubbogo{
		File:         file,
		Source:       desc.GetName(),
		ProtoPackage: desc.GetPackage(),
		Services:     make([]*Service, 0),
	}

	for _, service := range file.Services {
		serviceProtocolOpt, ok := proto.GetExtension(service.Desc.Options(), unified_idl_extend.E_ServiceProtocol).(*unified_idl_extend.ServiceProtocolTypeOption)
		if serviceProtocolSpecFlag && ok {
			if serviceProtocolOpt.GetProtocolName() != unified_idl_extend.ProtocolType_DUBBO.String() || serviceProtocolOpt == nil {
				// skip the service which is not dubbo protocol or does not have a service option
				continue
			}
		}
		serviceMethods := make([]*Method, 0)
		skipMethodFlag := false
		for _, method := range service.Methods {
			if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
				return nil, ErrStreamMethod
			}
			m := &Method{
				MethodName:  method.GoName,
				RequestType: g.QualifiedGoIdent(method.Input.GoIdent),
				ReturnType:  g.QualifiedGoIdent(method.Output.GoIdent),
			}

			methodProtocolOpt, ok := proto.GetExtension(method.Desc.Options(), unified_idl_extend.E_MethodProtocols).(*unified_idl_extend.MethodProtocolTypeOption)
			if methodProtocolSpecFlag && ok {
				flagDubboProtocol := false
				for _, protoName := range methodProtocolOpt.GetProtocolNames() {
					if protoName == unified_idl_extend.ProtocolType_DUBBO.String() {
						flagDubboProtocol = true
						break
					}
				}
				if !flagDubboProtocol || methodProtocolOpt == nil {
					// skip the method which is not dubbo protocol or does not have a method option
					skipMethodFlag = true
					continue
				}
			}
			skipMethodFlag = false

			methodOpt, ok := proto.GetExtension(method.Desc.Options(), unified_idl_extend.E_MethodExtend).(*unified_idl_extend.Hessian2MethodOptions)

			invokeName := util.ToLower(method.GoName)
			if ok && methodOpt != nil {
				invokeName = methodOpt.MethodName
			}
			m.InvokeName = invokeName

			inputOpt, ok := proto.GetExtension(method.Input.Desc.Options(), unified_idl_extend.E_MessageExtend).(*unified_idl_extend.Hessian2MessageOptions)
			if ok && inputOpt.ExtendArgs {
				m.RequestExtendArgs = true
				for _, field := range method.Input.Fields {
					goType, _ := util.FieldGoType(g, field)

					opt, ok := proto.GetExtension(field.Desc.Options(), unified_idl_extend.E_FieldExtend).(*unified_idl_extend.Hessian2FieldOptions)
					if ok && opt != nil && opt.IsWrapper && goType != "bool" {
						goType = "*" + goType
					}

					m.ArgsType = append(m.ArgsType, goType)
					m.ArgsName = append(m.ArgsName, util.ToLower(field.GoName))
				}
			}

			outputOpt, ok := proto.GetExtension(method.Output.Desc.Options(), unified_idl_extend.E_MessageExtend).(*unified_idl_extend.Hessian2MessageOptions)
			if ok && outputOpt != nil && outputOpt.ExtendArgs {
				m.ResponseExtendArgs = true
				if len(method.Output.Fields) == 0 {
					return nil, ErrNoExtendArgsRespFieldNum
				}
				if len(method.Output.Fields) != 1 {
					return nil, ErrMoreExtendArgsRespFieldNum
				}

				field := method.Output.Fields[0]
				goType, _ := util.FieldGoType(g, field)

				opt, ok := proto.GetExtension(field.Desc.Options(), unified_idl_extend.E_FieldExtend).(*unified_idl_extend.Hessian2FieldOptions)
				if ok && opt != nil && opt.IsWrapper {
					goType = "*" + goType
					m.ResponseIsWrapper = true
				}

				m.ReturnType = goType
			}
			if !skipMethodFlag {
				serviceMethods = append(serviceMethods, m)
			}
		}
		if len(serviceMethods) == 0 {
			continue
		}

		serviceOpt, ok := proto.GetExtension(service.Desc.Options(), unified_idl_extend.E_ServiceExtend).(*unified_idl_extend.Hessian2ServiceOptions)
		interfaceName := fmt.Sprintf("%s.%s", dubboGo.ProtoPackage, service.GoName)
		if ok && serviceOpt != nil {
			interfaceName = serviceOpt.InterfaceName
		}
		dubboGo.Services = append(dubboGo.Services, &Service{
			ServiceName:   service.GoName,
			Methods:       serviceMethods,
			InterfaceName: interfaceName,
		})
	}

	if len(dubboGo.Services) == 0 {
		return nil, fmt.Errorf("no service found in %s", dubboGo.Source)
	}

	return dubboGo, nil
}
