package modals

type User struct {
    UserName string `json:"userName"`
    UserProfileAttributes UserProfileAttributes `json:"userProfileAttributes"`
}

type UserProfileAttributes struct {
    DiscordUserName string `json:"discordUsername"`
    DiscordImageUrl string `json:"discordImageUrl"`
    SteamUserName   string `json:"steamUserName"`
    SteamUserId     string `json:"steamUserId"`
    Wins            string `json:"wins"`
    Losses          string `json:"losses"`
    TopGame         string `json:"topGame"`
    BestGame        string `json:"bestGame"`
    WorstGame       string `json:"worstGame"`
}

type BetResponse struct {
    BetId     string `json:"betId"`
    Initiator string `json:"initiator"`
    Opponent  string `json:"opponent"`
    Game      string `json:"game"`
    Amount    string `json:"amount"`
    Status    string `json:"status"`
    Accepted  string `json:"accepted"`
    NeedsToAccept string `json:"needsToAccept"`
    Winner    string `json:"winner"`
    CreatedOn string `json:"createdOn"`
    CreatedBy string `json:"createdBy"`
    UpdatedOn string `json:"updatedOn"`
    UpdatedBy string `json:"updatedBy"`
}

type UserResponse struct {
    UserName string `json:"userName"`
    UserProfileAttributes UserProfileAttributesResponse `json:"userProfileAttributesResponse,omitempty"`
    Bets []BetResponse `json:"bets,omitempty"`
}

type UserProfileAttributesResponse struct {
    DiscordUserName string `json:"discordUsername"`
    DiscordImageUrl string `json:"discordImageUrl"`
    SteamUserName   string `json:"steamUserName"`
    SteamUserId     string `json:"steamUserId"`
	Wins            string `json:"wins"`
    Losses          string `json:"losses"`
    TopGame         string `json:"topGame"`
    BestGame        string `json:"bestGame"`
    WorstGame       string `json:"worstGame"`
}