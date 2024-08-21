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

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dubbogo/protoc-gen-go-dubbo/generator"
	"github.com/dubbogo/protoc-gen-go-dubbo/internal/version"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	usage = "Flags:\n  -h, --help\tPrint this help and exit.\n      --version\tPrint the version and exit."
)

func main() {

	var flags flag.FlagSet
	protocolSpecFlag := flags.Bool("in_file_protocol_option", false, "enable the in-file transport protocol option")

	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Fprintln(os.Stdout, version.Version)
		os.Exit(0)
	}
	if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Fprintln(os.Stdout, usage)
		os.Exit(0)
	}
	if len(os.Args) != 1 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	flags.Init("protoc-gen-go-dubbo", flag.ContinueOnError)

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, file := range gen.Files {
			if file.Generate {
				filename := file.GeneratedFilenamePrefix + ".dubbo.go"
				g := gen.NewGeneratedFile(filename, file.GoImportPath)
				dubboGo, err := generator.ProcessProtoFile(g, file, *protocolSpecFlag)
				if err != nil {
					return err
				}
				generator.GenDubbo(g, dubboGo)
			}
		}
		return nil
	})
}
