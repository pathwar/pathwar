/* eslint-disable no-unused-vars */
// import { baseApi } from "./index";
import axios from "axios";

export function getAllTeams() {
    return axios.get("https://gist.githubusercontent.com/moul/826ef89d52651570a396ef3210a72e40/raw/e95d0e0391abca995949ab1258d5569e0b5ec356/GET%2520teams.json")
    // Uncomment line to use base api with auth token
    // return baseApi.get(`/teams`, { userId });
}
export function getUserTeams(userID) {

    const apiForExample = axios.create();
    apiForExample.interceptors.response.use((response) => {
        // eslint-disable-next-line no-console
        let interceptedResponse = {
          ...response,
          data: {
            "items": [
                {
                  "metadata": {
                    "id": "2dNuq9SuRvKcVHDOfNJSFQ==",
                    "created_at": "2019-04-25T11:41:24Z",
                    "updated_at": "2019-04-25T11:41:24Z"
                  },
                  "name": "chartreuse",
                  "gravatar_url": "http://www.internationalmonetize.io/harness/communities",
                  "locale": "fr_FR",
                  "last_active": true
                }
            ]
          }
        }

        return interceptedResponse;
      // eslint-disable-next-line no-console
      }, (error) => { console.log(error) });

      return apiForExample.get("https://gist.githubusercontent.com/moul/826ef89d52651570a396ef3210a72e40/raw/e95d0e0391abca995949ab1258d5569e0b5ec356/GET%2520teams.json")
   
    // Uncomment line to use base api with auth token
    // return baseApi.get(`/user-teams`, { userID });
}

export function joinTeam(userID, teamID) {
  // return axios.post("/join-teams", {userID, teamID})
}