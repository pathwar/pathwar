/* eslint-disable no-unused-vars */
import { baseApi } from "./index";

export function getAllOrganizations() {
    return baseApi.get(`/organizations`);
}
export function getUserOrganizations(userID) {
    return baseApi.get(`/user-teams`, { userID });
}

export function joinOrganization(userID, teamID) {
  // return axios.post("/join-teams", {userID, teamID})
}
