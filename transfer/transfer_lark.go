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

package transfer

import (
	"sync/atomic"
	"time"

	"github.com/vogo/logger"
)

// LarkTransfer transfer to support lark.
type LarkTransfer struct {
	id           string
	url          string
	transferring int32 // whether transferring message
}

// TypeLark transfer type lark.
const TypeLark = "lark"

const (
	larkMessageDataFixedBytesNum = 4
	larkMessageDataMaxLength     = 1024
	larkMessageTransferInterval  = time.Second * 5
)

var (
	// nolint:gochecknoglobals // ignore this
	larkTextMessageDataPrefix = []byte(`{"msg_type":"text","content":{"text":"[logtail-`)

	// nolint:gochecknoglobals // ignore this
	larkTextMessageDataSuffix = []byte(`"}}`)
)

func (d *LarkTransfer) Name() string {
	return d.id
}
func (d *LarkTransfer) Start() error { return nil }

func (d *LarkTransfer) Stop() error { return nil }

// Trans transfer data to Lark.
// nolint:dupl // ignore duplicated code for easy maintenance for diff transfers.
func (d *LarkTransfer) Trans(serverID string, data ...[]byte) error {
	if !atomic.CompareAndSwapInt32(&d.transferring, 0, 1) {
		// ignore message to
		return nil
	}

	go func() {
		<-time.After(larkMessageTransferInterval)
		atomic.StoreInt32(&d.transferring, 0)
	}()

	size := larkMessageDataFixedBytesNum + len(data)
	list := make([][]byte, size)
	list[0] = larkTextMessageDataPrefix
	list[1] = []byte(serverID)
	list[2] = messageTitleContentSplit

	idx := 3
	messageRemainCapacity := larkMessageDataMaxLength

	for _, bytes := range data {
		if messageRemainCapacity <= 0 {
			break
		}

		bytes = EscapeLimitJSONBytes(bytes, messageRemainCapacity)

		list[idx] = bytes
		idx++

		messageRemainCapacity -= len(bytes)
	}

	list[idx] = larkTextMessageDataSuffix

	if err := httpTrans(d.url, list[:idx+1]...); err != nil {
		logger.Errorf("lark error: %v", err)
	}

	return nil
}

// NewLarkTransfer initialize a lark transfer.
func NewLarkTransfer(id, url string) *LarkTransfer {
	return &LarkTransfer{
		id:           id,
		url:          url,
		transferring: 0,
	}
}
