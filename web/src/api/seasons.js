import { baseApi } from "./index"

export function postPreferences(seasonID) {
  return baseApi.post(`/user/preferences`, { active_season_id: seasonID })
}

export function getAllSeasonTeams(seasonID) {
  const urlIdParam = encodeURIComponent(seasonID)
  return baseApi.get(`/teams?season_id=${urlIdParam}`)
}

//TODO: Verify endpoint to return all seasons
export function getAllSeasons() {
  return baseApi.get(`/seasons`)
}

export function getChallenges(seasonID) {
  return baseApi.get(`/season-challenges?season_id=${seasonID}`)
}

export function getTeamDetails(teamID) {
  const urlIdParam = encodeURIComponent(teamID)
  return baseApi.get(`/team?team_id=${urlIdParam}`)
}

export function getChallengeDetails(challengeID) {
  const urlIdParam = encodeURIComponent(challengeID)
  return baseApi.get(`/season-challenge?season_challenge_id=${urlIdParam}`)
}

export function postBuyChallenge(challengeID, teamID) {
  return baseApi.post(`/season-challenge/buy`, { "season_challenge_id": challengeID, "team_id": teamID })
}

export function postValidateChallenge(validateData) {
  const { subscriptionID, passphrase, comment } = validateData;
  return baseApi.post(`/season-challenge/validate`, { "challenge_subscription_id": subscriptionID, "passphrase": passphrase, "comment": comment }  )
}

export function postCloseChallenge(subscriptionID) {
  return baseApi.post(`/season-challenge/close`, { "challenge_subscription_id": subscriptionID }  )
}
