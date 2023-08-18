package fmt

import (
	"github.com/fundrick/bets-svc/internal/modals"

    "strings"
)

func FormatGetUserResponse(user []map[string]interface{}) modals.UserResponse {
	var userResponse modals.UserResponse

	for _, v := range user {
		userAspect := v["userAspects"].(string)
		splitAspectType := strings.Split(userAspect, ":")
		aspectType := splitAspectType[0]
		aspectValue := splitAspectType[1]
		switch aspectType {
		case "PROFILE":
			userResponse.UserName = strings.Split(v["userName"].(string), ":")[1]
			userResponse.UserProfileAttributes = modals.UserProfileAttributesResponse{
				DiscordUserName: v["discordUsername"].(string),
				DiscordImageUrl: v["discordImageUrl"].(string),
				SteamUserName: v["steamUserName"].(string),
				SteamUserId: v["steamUserId"].(string),
				Wins: v["wins"].(string),
				Losses: v["losses"].(string),
				TopGame: v["topGame"].(string),
				BestGame: v["bestGame"].(string),
				WorstGame: v["worstGame"].(string),
			}
		case "BET":
			bet := modals.BetResponse{
				BetId:     aspectValue,
				Opponent:  v["opponent"].(string),
				Game:      v["game"].(string),
				Amount:    v["amount"].(string),
				Status:    v["status"].(string),
				Winner:    v["winner"].(string),
				NeedsToAccept: v["needsToAccept"].(string),
				Initiator: v["initiator"].(string),
				Accepted: v["accepted"].(string),
				CreatedOn: v["createdOn"].(string),
				CreatedBy: v["createdBy"].(string),
				UpdatedOn: v["updatedOn"].(string),
				UpdatedBy: v["updatedBy"].(string),
			}
			userResponse.Bets = append(userResponse.Bets, bet)
		}
	}

	return userResponse
}

func FormatUersResponse(users []map[string]interface{}) modals.UsersResponse {
	var usersResponse modals.UsersResponse

	for _, userData := range users {
		user := userData["userName"].(string)
		formatedUser := strings.Split(user, ":")[1]
		usersResponse.UserNames = append(usersResponse.UserNames, formatedUser)
	}

	return usersResponse
}