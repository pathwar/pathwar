import { baseApi } from "./index";

//Season main calls
export function postPreferencesByID(seasonID) {
  return baseApi.post(`/user/preferences`, { active_season_id: seasonID });
}

export function postPreferencesBySlug(seasonSlug) {
  return baseApi.post(`/user/preferences`, { active_season_slug: seasonSlug });
}

//TODO: Verify endpoint to return all seasons
export function getAllSeasons() {
  return baseApi.get(`/seasons`);
}

//Team calls
export function getTeamDetails(teamID) {
  const urlIdParam = encodeURIComponent(teamID);
  return baseApi.get(`/team?team_id=${urlIdParam}`);
}

export function getAllSeasonTeams(seasonID) {
  const urlIdParam = encodeURIComponent(seasonID);
  return baseApi.get(`/teams?season_id=${urlIdParam}`);
}

export function postCreateTeam(seasonID, name) {
  return baseApi.post(`/team`, { season_id: seasonID, name: name });
}

//Challenge calls
export function getChallenges(seasonID) {
  return baseApi.get(`/season-challenges?season_id=${seasonID}`);
}

export function getChallengeDetails(challengeID) {
  const urlIdParam = encodeURIComponent(challengeID);
  return baseApi.get(`/season-challenge?season_challenge_id=${urlIdParam}`);
}

export function postBuyChallenge(flavorChallengeID, seasonID) {
  return baseApi.post(`/season-challenge/buy`, {
    flavor_id: flavorChallengeID,
    season_id: seasonID,
  });
}

export function postValidateChallenge(validateData) {
  const { subscriptionID, passphrases, comment } = validateData;
  return baseApi.post(`/challenge-subscription/validate`, {
    challenge_subscription_id: subscriptionID,
    passphrases: passphrases,
    comment: comment,
  });
}

export function postCloseChallenge(subscriptionID) {
  return baseApi.post(`/challenge-subscription/close`, {
    challenge_subscription_id: subscriptionID,
  });
}
