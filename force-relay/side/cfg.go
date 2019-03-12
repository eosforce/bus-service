package side

import eos "github.com/eosforce/goforceio"

type actionToCommit struct {
	account    eos.AccountName
	actionName eos.ActionName
}

type relayCfg struct {
	Chain           eos.Name
	TransferAccount eos.AccountName
	ActionToCommit  []actionToCommit
}

var cfg *relayCfg

// AppendActionInfo append a action info to commit
func (r *relayCfg) AppendActionInfo(account eos.AccountName, actionName eos.ActionName) {
	r.ActionToCommit = append(r.ActionToCommit, actionToCommit{
		account:    account,
		actionName: actionName,
	})
}

// IsActionNeedToCommit if a action need to commit to relay chain
func (r *relayCfg) IsActionNeedToCommit(account eos.AccountName, actionName eos.ActionName) bool {
	for _, act := range r.ActionToCommit {
		if act.actionName == actionName && act.account == account {
			return true
		}
	}

	return false
}

// NewCfg new cfg to relay info
func NewCfg(chain, transfer string) *relayCfg {
	return &relayCfg{
		Chain:           eos.Name(chain),
		TransferAccount: eos.AN(transfer),
		ActionToCommit:  make([]actionToCommit, 0, 32),
	}
}

// SetCfg set cfg from NewCfg
func SetCfg(cc *relayCfg) {
	cfg = cc
}
