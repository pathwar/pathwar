/* eslint-disable no-unused-vars */
import { baseApi } from "./index";

export function getAllTeams() {
    return baseApi.get(`/teams`);
}
export function getUserTeams(userID) {
    return baseApi.get(`/user-teams`, { userID });
}

export function joinTeam(userID, teamID) {
  // return axios.post("/join-teams", {userID, teamID})
}
