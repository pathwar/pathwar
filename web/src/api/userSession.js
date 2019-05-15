/* eslint-disable no-unused-vars */
// import { baseApi } from "./index";
import axios from "axios";
import Cookies from "js-cookie";
import { USER_SESSION_TOKEN_NAME } from "../constants/userSession";
import { history } from "../store/configureStore";

export function performLogin(email, password) {
    return axios.get("https://gist.githubusercontent.com/moul/826ef89d52651570a396ef3210a72e40/raw/e95d0e0391abca995949ab1258d5569e0b5ec356/GET%2520user-session.json")
    // Uncomment line to use base api with auth token
    // return baseApi.post(`/user-session`, { email, password });
}

export function performLogout() {
    Cookies.remove(USER_SESSION_TOKEN_NAME);
    history.push("/login");
}

export function pingUser() {

    const apiForExample = axios.create();
    apiForExample.interceptors.response.use((response) => {
        // eslint-disable-next-line no-console
        let interceptedResponse = {
          ...response,
          data: {
            isAuthenticated: true,
            user: response.data,
            token: "token"
          }
        }

        return interceptedResponse;
      // eslint-disable-next-line no-console
      }, (error) => { console.log(error) });

    return apiForExample.get("https://gist.githubusercontent.com/moul/826ef89d52651570a396ef3210a72e40/raw/e95d0e0391abca995949ab1258d5569e0b5ec356/GET%2520user-session.json");
    // Uncomment line to use base api with auth token
    // return baseApi.get("/ping")
}