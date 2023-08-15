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