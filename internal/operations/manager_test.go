// Copyright © 2021 Kaleido, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in comdiliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or imdilied.
// See the License for the specific language governing permissions and
// limitations under the License.

package operations

import (
	"context"
	"testing"

	"github.com/hyperledger/firefly/internal/config"
	"github.com/hyperledger/firefly/internal/txcommon"
	"github.com/hyperledger/firefly/mocks/databasemocks"
	"github.com/hyperledger/firefly/mocks/tokenmocks"
	"github.com/hyperledger/firefly/pkg/fftypes"
	"github.com/hyperledger/firefly/pkg/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestOperations(t *testing.T) (*operationsManager, func()) {
	config.Reset()
	mdi := &databasemocks.Plugin{}
	mti := &tokenmocks.Plugin{}
	mti.On("Name").Return("ut_tokens").Maybe()
	ctx, cancel := context.WithCancel(context.Background())
	om, err := NewOperationsManager(ctx, mdi, map[string]tokens.Plugin{"magic-tokens": mti})
	assert.NoError(t, err)
	return om.(*operationsManager), cancel
}

func TestInitFail(t *testing.T) {
	_, err := NewOperationsManager(context.Background(), nil, nil)
	assert.Regexp(t, "FF10128", err)
}

func TestStartOperationNotSupported(t *testing.T) {
	om, cancel := newTestOperations(t)
	defer cancel()

	op := &fftypes.Operation{}

	err := om.StartOperation(context.Background(), op)
	assert.Regexp(t, "FF10346", err)
}

func TestStartOperationCreatePool(t *testing.T) {
	om, cancel := newTestOperations(t)
	defer cancel()

	op := &fftypes.Operation{
		ID:   fftypes.NewUUID(),
		Type: fftypes.OpTypeTokenCreatePool,
	}
	pool := &fftypes.TokenPool{
		Connector: "magic-tokens",
	}
	err := txcommon.AddTokenPoolCreateInputs(op, pool)
	assert.NoError(t, err)

	mti := om.tokens["magic-tokens"].(*tokenmocks.Plugin)
	mti.On("CreateTokenPool", context.Background(), op.ID, mock.AnythingOfType("*fftypes.TokenPool")).Return(false, nil)

	err = om.StartOperation(context.Background(), op)
	assert.NoError(t, err)

	mti.AssertExpectations(t)
}

func TestStartOperationCreatePoolBadInput(t *testing.T) {
	om, cancel := newTestOperations(t)
	defer cancel()

	op := &fftypes.Operation{
		ID:   fftypes.NewUUID(),
		Type: fftypes.OpTypeTokenCreatePool,
		Input: fftypes.JSONObject{
			"test": map[bool]bool{true: false},
		},
	}

	err := om.StartOperation(context.Background(), op)
	assert.Regexp(t, "FF10151", err)
}

func TestStartOperationActivatePool(t *testing.T) {
	om, cancel := newTestOperations(t)
	defer cancel()

	op := &fftypes.Operation{
		ID:   fftypes.NewUUID(),
		Type: fftypes.OpTypeTokenActivatePool,
	}
	info := fftypes.JSONObject{
		"some": "info",
	}
	pool := &fftypes.TokenPool{
		ID:        fftypes.NewUUID(),
		Connector: "magic-tokens",
	}
	txcommon.AddTokenPoolActivateInputs(op, pool.ID, info)

	mdi := om.database.(*databasemocks.Plugin)
	mti := om.tokens["magic-tokens"].(*tokenmocks.Plugin)
	mdi.On("GetTokenPoolByID", context.Background(), pool.ID).Return(pool, nil)
	mti.On("ActivateTokenPool", context.Background(), op.ID, mock.AnythingOfType("*fftypes.TokenPool"), info).Return(false, nil)

	err := om.StartOperation(context.Background(), op)
	assert.NoError(t, err)

	mdi.AssertExpectations(t)
	mti.AssertExpectations(t)
}

func TestStartOperationActivatePoolBadInput(t *testing.T) {
	om, cancel := newTestOperations(t)
	defer cancel()

	op := &fftypes.Operation{
		ID:   fftypes.NewUUID(),
		Type: fftypes.OpTypeTokenActivatePool,
		Input: fftypes.JSONObject{
			"id": "bad",
		},
	}

	err := om.StartOperation(context.Background(), op)
	assert.Regexp(t, "FF10142", err)
}

func TestStartOperationTokenTransfer(t *testing.T) {
	om, cancel := newTestOperations(t)
	defer cancel()

	op := &fftypes.Operation{
		ID:   fftypes.NewUUID(),
		Type: fftypes.OpTypeTokenTransfer,
	}
	pool := &fftypes.TokenPool{
		ID:         fftypes.NewUUID(),
		Connector:  "magic-tokens",
		ProtocolID: "F1",
	}
	transfer := &fftypes.TokenTransfer{
		Pool:      pool.ID,
		Connector: "magic-tokens",
		Type:      fftypes.TokenTransferTypeTransfer,
	}
	txcommon.AddTokenTransferInputs(op, transfer)

	mdi := om.database.(*databasemocks.Plugin)
	mti := om.tokens["magic-tokens"].(*tokenmocks.Plugin)
	mdi.On("GetTokenPoolByID", context.Background(), pool.ID).Return(pool, nil)
	mti.On("TransferTokens", context.Background(), op.ID, "F1", mock.AnythingOfType("*fftypes.TokenTransfer")).Return(nil)

	err := om.StartOperation(context.Background(), op)
	assert.NoError(t, err)

	mdi.AssertExpectations(t)
	mti.AssertExpectations(t)
}

func TestStartOperationTokenTransferBadInput(t *testing.T) {
	om, cancel := newTestOperations(t)
	defer cancel()

	op := &fftypes.Operation{
		ID:   fftypes.NewUUID(),
		Type: fftypes.OpTypeTokenTransfer,
		Input: fftypes.JSONObject{
			"test": map[bool]bool{true: false},
		},
	}

	err := om.StartOperation(context.Background(), op)
	assert.Regexp(t, "FF10151", err)
}
