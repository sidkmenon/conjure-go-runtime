// Copyright (c) 2020 Palantir Technologies. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	"github.com/palantir/conjure-go-runtime/v2/conjure-go-contract/codecs"
	werror "github.com/palantir/witchcraft-go-error"
)

// UnmarshalError attempts to deserialize the message to a known implementation of Error.
// Custom error types should be registered using RegisterErrorType.
// If the ErrorName is not recognized, a genericError is returned with all params marked unsafe.
// If we fail to unmarshal to a generic SerializableError or to the type specified by ErrorName, an error is returned.
func UnmarshalError(body []byte) (Error, error) {
	return UnmarshalErrorWithDecoder(globalRegistry, body)
}

// UnmarshalErrorWithDecoder attempts to deserialize the message to a known implementation of Error
// using the provided ConjureErrorDecoder.
func UnmarshalErrorWithDecoder(ced ConjureErrorDecoder, body []byte) (Error, error) {
	var name struct {
		Name string `json:"errorName"`
	}
	if err := codecs.JSON.Unmarshal(body, &name); err != nil {
		return nil, werror.Wrap(err, "failed to unmarshal body as conjure error")
	}
	cErr, err := ced.DecodeConjureError(name.Name, body)
	if err != nil {
		return nil, werror.Convert(err)
	}
	return cErr, nil
}
