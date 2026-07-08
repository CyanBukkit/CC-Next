package stuck

// RecoveryAction selects how to recover when Claude appears stuck.
type RecoveryAction int

const (
	ActionAutoContinue RecoveryAction = iota
	ActionAlertUser
)

// RecoveryHandler executes the configured action when a stuck state is detected.
type RecoveryHandler struct {
	action  RecoveryAction
	onAuto  func()
	onAlert func()
}

// NewRecoveryHandler creates a handler with the selected action callbacks.
func NewRecoveryHandler(action RecoveryAction, onAuto, onAlert func()) *RecoveryHandler {
	return &RecoveryHandler{
		action:  action,
		onAuto:  onAuto,
		onAlert: onAlert,
	}
}

// Handle executes the appropriate recovery action.
func (rh *RecoveryHandler) Handle() {
	switch rh.action {
	case ActionAutoContinue:
		if rh.onAuto != nil {
			rh.onAuto()
		}
	case ActionAlertUser:
		if rh.onAlert != nil {
			rh.onAlert()
		}
	}
}

// SetAction changes the recovery action for subsequent Handle calls.
func (rh *RecoveryHandler) SetAction(action RecoveryAction) {
	rh.action = action
}
