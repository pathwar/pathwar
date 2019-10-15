import { baseApi } from "./index";

export function postPreferences(tournamentID) {
  return baseApi.post(`/preferences`, {"active_tournament_id": tournamentID});
}

export function getAllTournamentTeams(tournamentID) {
  return baseApi.get(`/tournament/teams?tournament_id=${tournamentID}`)
}

export function getAllTournaments() {
    return baseApi.get(`/tournaments`);
}

export function getChallenges(tournamentID) {
    return baseApi.get(`/challenges`, { id: tournamentID });
}
