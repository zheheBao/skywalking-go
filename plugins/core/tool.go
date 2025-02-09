// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package core

import (
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	commonv3 "skywalking.apache.org/repo/goapi/collect/common/v3"
)

// Millisecond converts time to unix millisecond
func Millisecond(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// UUID generate UUID
func UUID() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(id.String(), "-", ""), nil
}

// GenerateGlobalID generates global unique id
func GenerateGlobalID() (globalID string, err error) {
	return UUID()
}

func ProcessNo() string {
	if os.Getpid() > 0 {
		return strconv.Itoa(os.Getpid())
	}
	return ""
}

func HostName() string {
	if hs, err := os.Hostname(); err == nil {
		return hs
	}
	return "unknown"
}

func OSName() string {
	return runtime.GOOS
}

func AllIPV4() (ipv4s []string) {
	adders, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	for _, addr := range adders {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ipv4 := ipNet.IP.String()
				if ipv4 == "127.0.0.1" || ipv4 == "localhost" {
					continue
				}
				ipv4s = append(ipv4s, ipv4)
			}
		}
	}
	return
}

func IPV4() string {
	ipv4s := AllIPV4()
	if len(ipv4s) > 0 {
		return ipv4s[0]
	}
	return "no-hostname"
}

func buildOSInfo() (props []*commonv3.KeyStringValuePair) {
	processNo := ProcessNo()
	if processNo != "" {
		kv := &commonv3.KeyStringValuePair{
			Key:   "Process No.",
			Value: processNo,
		}
		props = append(props, kv)
	}

	hostname := &commonv3.KeyStringValuePair{
		Key:   "hostname",
		Value: HostName(),
	}
	props = append(props, hostname)

	language := &commonv3.KeyStringValuePair{
		Key:   "language",
		Value: "go",
	}
	props = append(props, language)

	osName := &commonv3.KeyStringValuePair{
		Key:   "OS Name",
		Value: OSName(),
	}
	props = append(props, osName)

	ipv4s := AllIPV4()
	if len(ipv4s) > 0 {
		for _, ipv4 := range ipv4s {
			kv := &commonv3.KeyStringValuePair{
				Key:   "ipv4",
				Value: ipv4,
			}
			props = append(props, kv)
		}
	}
	return props
}
