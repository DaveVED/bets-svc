package modals

type Bet struct {
    InitiatorUserName string `json:"initiatorUserName"`
    OpponentUserName  string `json:"opponentUserName"`
    Game string `json:"game"`
    Amount string `json:"amount"`
}

type BetUpdate struct {
    Initiator      string    `json:"initiator"`
    Opponent       string    `json:"opponent"`
    Status         string    `json:"status"`
    Accepted       string    `json:"accepted"`
    NeedsToAccept  string    `json:"needsToAccept"`
    Winner         string    `json:"winner"`
    UpdatedBy      string    `json:"updatedBy"`
}