/* eslint-disable no-unused-vars */
import { baseApi } from "./index";

export function getAllOrganizations() {
  return baseApi.get(`/organizations`);
}

//TODO: Verify  endpoint to join ang organization
export function joinOrganization(userID, teamID) {
  // return axios.post("/join-teams", {userID, teamID})
}
