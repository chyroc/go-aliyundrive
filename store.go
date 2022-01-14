/**
 * Copyright 2022 chyroc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package aliyundrive

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"time"
)

type Token struct {
	AccessToken  string    `json:"access_token"`
	ExpiredAt    time.Time `json:"expired_at"` // access-token 的过期时间，秒级
	RefreshToken string    `json:"refresh_token"`
}

// Store 定义一个存储的接口，使用者可以按照自己的需求来实现
type Store interface {
	Get(ctx context.Context, key string) (*Token, error)
	Set(ctx context.Context, token *Token) error
}

type FileStore struct {
	file string
}

func (r *FileStore) Get(ctx context.Context, key string) (*Token, error) {
	bs, err := ioutil.ReadFile(r.file)
	if err != nil {
		return nil, err
	}
	token := new(Token)
	err = json.Unmarshal(bs, token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (r *FileStore) Set(ctx context.Context, token *Token) error {
	bs, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(r.file, bs, 0o666)
}

func NewFileStore(file string) Store {
	return &FileStore{
		file: file,
	}
}
