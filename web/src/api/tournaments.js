/* eslint-disable no-unused-vars */
// import { baseApi } from "./index";
import axios from "axios";

export function setTournamentActive(activeTeamID, tournamentID) {
        // return baseApi.post(`/user-tournaments`, { activeTeamID, tournamentID });
}
export function getAllTournaments() {
    return axios.get("https://gist.githubusercontent.com/moul/826ef89d52651570a396ef3210a72e40/raw/e95d0e0391abca995949ab1258d5569e0b5ec356/GET%2520tournaments.json")
    // Uncomment line to use base api with auth token
    // return baseApi.get(`/tournaments`);
}

export function getTeamTournaments(teamID) {
    return axios.get("https://gist.githubusercontent.com/moul/826ef89d52651570a396ef3210a72e40/raw/e95d0e0391abca995949ab1258d5569e0b5ec356/GET%2520tournaments.json")
    // Uncomment line to use base api with auth token
    // return baseApi.get(`/team-tournaments`, { teamID });
}

export function getLevels(tournamentID) {
    return axios.get("https://gist.githubusercontent.com/moul/826ef89d52651570a396ef3210a72e40/raw/e95d0e0391abca995949ab1258d5569e0b5ec356/GET%2520levels.json")
    // Uncomment line to use base api with auth token
    // return baseApi.get(`/levels`, { id: tournamentID });
}