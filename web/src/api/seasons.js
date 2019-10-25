import { baseApi } from "./index";

export function postPreferences(seasonID) {
  return baseApi.post(`/user/preferences`, {"active_season_id": seasonID});
}

export function getAllSeasonTeams(seasonID) {
  const urlIdParam = encodeURIComponent(seasonID);
  return baseApi.get(`/season/teams?season_id=${urlIdParam}`)
}

export function getAllSeasons() {
    return baseApi.get(`/seasons`);
}

export function getChallenges(seasonID) {
    return baseApi.get(`/challenges`, { id: seasonID });
}

export function getChallengeDetails(challengeID) {
  const urlIdParam = encodeURIComponent(challengeID);
  return baseApi.get(`/challenge?challenge_id=${urlIdParam}`);
}
