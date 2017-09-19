// Copyright (C) 2017 go-nebulas authors
//
// This file is part of the go-nebulas library.
//
// the go-nebulas library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// the go-nebulas library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-nebulas library.  If not, see <http://www.gnu.org/licenses/>.
//

package pow

import (
	"github.com/nebulasio/go-nebulas/consensus"
	log "github.com/sirupsen/logrus"
)

const (
	// Minted minted state key
	Minted = "minted"
)

// MintedState minted state, transite from @MiningState
type MintedState struct {
	p *Pow
}

// NewMintedState create MintedState instance.
func NewMintedState(p *Pow) *MintedState {
	state := &MintedState{p: p}
	return state
}

// Event handle event.
func (state *MintedState) Event(e consensus.Event) (bool, consensus.State) {
	return false, nil
}

// Enter called when transiting to this state.
func (state *MintedState) Enter(data interface{}) {
	log.Debug("MintedState.Enter: enter.")

	p := state.p
	bkPool := p.chain.BlockPool()

	// process minted block.
	if p.miningBlock.Nonce() > 0 {
		log.Info("MintedState.Enter: process minted block.")

		// sign.
		p.miningBlock.Sign()

		// send new block to network.
		p.nm.BroadcastBlock(p.miningBlock)

		bkPool.AddLocalBlock(p.miningBlock)
		if p.receivedBlock == nil {
			p.receivedBlock = p.miningBlock
		}
	}

	// process the received block.
	if p.receivedBlock != nil {
		p.receivedBlock = nil
		p.ForkChoice()
	}

	// move to prepare state.
	state.p.TransiteByKey(Prepare, nil)
}

// Leave called when leaving this state.
func (state *MintedState) Leave(data interface{}) {
	log.Debug("MintedState.Leave: leave.")
}
