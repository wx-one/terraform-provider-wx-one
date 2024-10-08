// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func isNotFoundError(err error) bool {
	if errList, ok := err.(gqlerror.List); ok {
		gqlerr := &gqlerror.Error{}
		if errList.As(&gqlerr) {
			if errorCode, ok := gqlerr.Extensions["code"].(string); ok {
				if errorCode == "NOT_FOUND" {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		} else {
			return false
		}
	} else {
		return false
	}
}
