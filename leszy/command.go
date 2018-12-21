// Copyright (c) 2018 Samsung Electronics Co., Ltd All Rights Reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

// command.go defines Command interface and types shared between Leszy Commands

package leszy

import (
	"github.com/spf13/cobra"

	b "github.com/SamsungSLAV/boruta/http/client"
)

// BaseCmd is the common (base) part of all types implementing Command
// interface.
type BaseCmd struct {
	Command *cobra.Command
	Clients *Clients
}

func (b *BaseCmd) Cmd() *cobra.Command {
	return b.Command
}

// Clients holds SLAV stack clients.
// TODO: add Weles
type Clients struct {
	Boruta *b.BorutaClient
}
