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

package sqlcommon

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/kaleido-io/firefly/internal/database"
	"github.com/kaleido-io/firefly/internal/fftypes"
	"github.com/kaleido-io/firefly/internal/i18n"
	"github.com/kaleido-io/firefly/internal/log"
)

var (
	namespaceColumns = []string{
		"id",
		"ntype",
		"name",
		"description",
		"created",
		"confirmed",
	}
	namespaceFilterTypeMap = map[string]string{
		"type": "ntype",
	}
)

func (s *SQLCommon) UpsertNamespace(ctx context.Context, namespace *fftypes.Namespace) (err error) {
	ctx, tx, autoCommit, err := s.beginOrUseTx(ctx)
	if err != nil {
		return err
	}
	defer s.rollbackTx(ctx, tx, autoCommit)

	// Do a select within the transaction to detemine if the UUID already exists
	namespaceRows, err := s.queryTx(ctx, tx,
		sq.Select("id").
			From("namespaces").
			Where(sq.Eq{"name": namespace.Name}),
	)
	if err != nil {
		return err
	}

	if namespaceRows.Next() {

		var id uuid.UUID
		_ = namespaceRows.Scan(&id)
		if namespace.ID != nil {
			if *namespace.ID != id {
				namespaceRows.Close()
				return database.IDMismatch
			}
		}
		namespace.ID = &id // Update on returned object
		namespaceRows.Close()

		// Update the namespace
		if _, err = s.updateTx(ctx, tx,
			sq.Update("namespaces").
				// Note we do not update ID
				Set("ntype", string(namespace.Type)).
				Set("name", namespace.Name).
				Set("description", namespace.Description).
				Set("created", namespace.Created).
				Set("confirmed", namespace.Confirmed).
				Where(sq.Eq{"name": namespace.Name}),
		); err != nil {
			return err
		}
	} else {
		namespaceRows.Close()

		if namespace.ID == nil {
			namespace.ID = fftypes.NewUUID()
		}

		if _, err = s.insertTx(ctx, tx,
			sq.Insert("namespaces").
				Columns(namespaceColumns...).
				Values(
					namespace.ID,
					string(namespace.Type),
					namespace.Name,
					namespace.Description,
					namespace.Created,
					namespace.Confirmed,
				),
		); err != nil {
			return err
		}
	}

	return s.commitTx(ctx, tx, autoCommit)
}

func (s *SQLCommon) namespaceResult(ctx context.Context, row *sql.Rows) (*fftypes.Namespace, error) {
	namespace := fftypes.Namespace{}
	err := row.Scan(
		&namespace.ID,
		&namespace.Type,
		&namespace.Name,
		&namespace.Description,
		&namespace.Created,
		&namespace.Confirmed,
	)
	if err != nil {
		return nil, i18n.WrapError(ctx, err, i18n.MsgDBReadErr, "namespaces")
	}
	return &namespace, nil
}

func (s *SQLCommon) GetNamespace(ctx context.Context, name string) (message *fftypes.Namespace, err error) {

	rows, err := s.query(ctx,
		sq.Select(namespaceColumns...).
			From("namespaces").
			Where(sq.Eq{"name": name}),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		log.L(ctx).Debugf("Namespace '%s' not found", name)
		return nil, nil
	}

	namespace, err := s.namespaceResult(ctx, rows)
	if err != nil {
		return nil, err
	}

	return namespace, nil
}

func (s *SQLCommon) GetNamespaces(ctx context.Context, filter database.Filter) (message []*fftypes.Namespace, err error) {

	query, err := s.filterSelect(ctx, "", sq.Select(namespaceColumns...).From("namespaces"), filter, namespaceFilterTypeMap)
	if err != nil {
		return nil, err
	}

	rows, err := s.query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	namespace := []*fftypes.Namespace{}
	for rows.Next() {
		d, err := s.namespaceResult(ctx, rows)
		if err != nil {
			return nil, err
		}
		namespace = append(namespace, d)
	}

	return namespace, err

}

func (s *SQLCommon) UpdateNamespace(ctx context.Context, name string, update database.Update) (err error) {

	ctx, tx, autoCommit, err := s.beginOrUseTx(ctx)
	if err != nil {
		return err
	}
	defer s.rollbackTx(ctx, tx, autoCommit)

	query, err := s.buildUpdate(ctx, sq.Update("namespaces"), update, namespaceFilterTypeMap)
	if err != nil {
		return err
	}
	query = query.Where(sq.Eq{"name": name})

	_, err = s.updateTx(ctx, tx, query)
	if err != nil {
		return err
	}

	return s.commitTx(ctx, tx, autoCommit)
}