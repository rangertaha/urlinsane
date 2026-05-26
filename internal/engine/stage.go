// Copyright 2024 Rangertaha. All Rights Reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package engine

import (
	"context"

	"github.com/rangertaha/urlinsane/internal/db"
)

// Stage is one phase of the scan pipeline. It consumes a stream of domains and
// produces a stream of domains; ctx cancels the whole pipeline. The pipeline is
// the ordered composition of stages — see Urlinsane.stages / Execute.
type Stage interface {
	Name() string
	Run(ctx context.Context, in <-chan *db.Domain) <-chan *db.Domain
}

// stageFunc adapts a named closure to a Stage. Most stages are stateless
// wrappers over an Urlinsane method, so a closure is all they need.
type stageFunc struct {
	name string
	run  func(ctx context.Context, in <-chan *db.Domain) <-chan *db.Domain
}

func (s stageFunc) Name() string { return s.name }

func (s stageFunc) Run(ctx context.Context, in <-chan *db.Domain) <-chan *db.Domain {
	return s.run(ctx, in)
}
