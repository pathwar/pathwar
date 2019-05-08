// import { baseApi } from "./index";
import axios from "axios";

export function getLevels() {
    return axios.get("https://gist.githubusercontent.com/moul/826ef89d52651570a396ef3210a72e40/raw/e95d0e0391abca995949ab1258d5569e0b5ec356/GET%2520levels.json")
    // Uncomment line to use base api with auth token
    // return baseApi.get(`/user-session`);
}