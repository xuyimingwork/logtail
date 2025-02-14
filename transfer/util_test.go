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

package transfer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vogo/logtail/transfer"
)

func TestEscapeLimitJsonBytes(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []byte(`ab`), transfer.EscapeLimitJSONBytes([]byte(`abcd`), 2))
	assert.Equal(t, []byte(`abcd`), transfer.EscapeLimitJSONBytes([]byte(`abcd`), 4))
	assert.Equal(t, []byte(`你好`), transfer.EscapeLimitJSONBytes([]byte(`你好世界`), 8))
	assert.Equal(t, []byte(`你好`), transfer.EscapeLimitJSONBytes([]byte(`你好世界`), 9))

	assert.Equal(t, []byte(`ab\"cd`), transfer.EscapeLimitJSONBytes([]byte(`ab"cd`), 6))
	assert.Equal(t, []byte(`ab\tcd`), transfer.EscapeLimitJSONBytes([]byte(`ab	cd`), 8))
	assert.Equal(t, []byte(`ab\ncd`), transfer.EscapeLimitJSONBytes([]byte("ab\ncd"), 8))
	assert.Equal(t, []byte(`abc\n`), transfer.EscapeLimitJSONBytes([]byte("abc\nd"), 4))
	assert.Equal(t, []byte(`abc\n`), transfer.EscapeLimitJSONBytes([]byte("abc\nd"), 4))

	assert.Equal(t, []byte(`test 操作异常`), transfer.EscapeLimitJSONBytes([]byte("test 操作异常"), 1024))
}
