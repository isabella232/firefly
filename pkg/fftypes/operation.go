// Copyright © 2021 Kaleido, Inc.
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

package fftypes

// OpType describes mechanical steps in the process that have to be performed,
// might be asynchronous, and have results in the back-end systems that might need
// to be correlated with messages by operators.
type OpType string

const (
	// OpTypeBlockchainBatchPin is a blockchain transaction to pin a batch
	OpTypeBlockchainBatchPin OpType = "BlockchainBatchPin"
	// OpTypePublicStorageBatchBroadcast is a public storage operation to store broadcast data
	OpTypePublicStorageBatchBroadcast OpType = "PublicStorageBatchBroadcast"
)

// OpStatus is the current status of an operation
type OpStatus string

const (
	// OpStatusPending indicates the operation has been submitted, but is not yet confirmed as successful or failed
	OpStatusPending OpStatus = "pending"
	// OpStatusSucceeded the infrastructure runtime has returned success for the operation.
	OpStatusSucceeded OpStatus = "succeeded"
	// OpStatusFailed happens when an error is reported by the infrastructure runtime
	OpStatusFailed OpStatus = "failed"
)

type Named interface {
	Name() string
}

// NewMessageOp creates a new operation for a message
func NewMessageOp(plugin Named, backendID string, msg *Message, opType OpType, opStatus OpStatus, recipient string) *Operation {
	return &Operation{
		ID:        NewUUID(),
		Plugin:    plugin.Name(),
		BackendID: backendID,
		Namespace: msg.Header.Namespace,
		Message:   msg.Header.ID,
		Data:      nil,
		Type:      opType,
		Recipient: recipient,
		Status:    opStatus,
		Created:   Now(),
	}
}

// NewMessageDataOp creates a new operation for a data
func NewMessageDataOp(plugin Named, backendID string, msg *Message, dataIDx int, opType OpType, opStatus OpStatus, recipient string) *Operation {
	return &Operation{
		ID:        NewUUID(),
		Plugin:    plugin.Name(),
		BackendID: backendID,
		Namespace: msg.Header.Namespace,
		Message:   msg.Header.ID,
		Data:      msg.Data[dataIDx].ID,
		Type:      opType,
		Recipient: recipient,
		Status:    opStatus,
		Created:   Now(),
	}
}

// Operation is a description of an action performed in an infrastructure runtime, such as sending a batch of data
type Operation struct {
	ID        *UUID    `json:"id"`
	Namespace string   `json:"namespace,omitempty"`
	Message   *UUID    `json:"message"`
	Data      *UUID    `json:"data,omitempty"`
	Type      OpType   `json:"type"`
	Recipient string   `json:"recipient,omitempty"`
	Status    OpStatus `json:"status"`
	Error     string   `json:"error,omitempty"`
	Plugin    string   `json:"plugin"`
	BackendID string   `json:"backendID"`
	Created   *FFTime  `json:"created,omitempty"`
	Updated   *FFTime  `json:"updated,omitempty"`
}