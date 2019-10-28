import { baseApi } from "./index";

export function postPreferences(seasonID) {
  return baseApi.post(`/user/preferences`, {"active_season_id": seasonID});
}

export function getAllSeasonTeams(seasonID) {
  const urlIdParam = encodeURIComponent(seasonID);
  return baseApi.get(`/teams?season_id=${urlIdParam}`)
}

export function getAllSeasons() {
    return baseApi.get(`/seasons`);
}

export function getChallenges(seasonID) {
    return baseApi.get(`/challenges`, { id: seasonID });
}

export function getTeamDetails(teamID) {
  const urlIdParam = encodeURIComponent(teamID);
  return baseApi.get(`/season/team?season_team_id=${urlIdParam}`);
}

export function getChallengeDetails(challengeID) {
  const urlIdParam = encodeURIComponent(challengeID);
  return baseApi.get(`/challenge?challenge_id=${urlIdParam}`);
}
