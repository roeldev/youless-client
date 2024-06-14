// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youless

import (
	"github.com/go-pogo/errors"
	"time"
)

const (
	ErrInvalidBaseURL errors.Msg = "invalid base url"
	ErrInvalidConfig  errors.Msg = "invalid config"
)

// Config is the configuration for a Client. It can be unmarshalled from json,
// yaml, env or flag values.
type Config struct {
	// BaseURL of the device.
	BaseURL string `json:"base_url" yaml:"baseUrl" default:"http://youless"`
	// Name of the device, is optional and used for logging/debugging.
	Name string `json:"name" yaml:"name" default:"YouLess"`
	// Timeout specifies a time limit for requests made by the http.Client used
	// by Client.
	Timeout time.Duration `json:"timeout" yaml:"timeout" default:"5s"`
	// Password used to connect with the device.
	Password string `json:"password" yaml:"password"`
	// PasswordFile contains the password used to connect with the device. When
	// both Password and PasswordFile are set, PasswordFile takes precedence.
	PasswordFile string `json:"password_file" yaml:"passwordFile"`
}

func (c Config) Validate() error {
	if c.BaseURL == "" {
		return errors.Wrap(ErrInvalidBaseURL, ErrInvalidConfig)
	}
	return nil
}
