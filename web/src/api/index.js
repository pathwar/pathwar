import axios from "axios";
import Cookies from "js-cookie";
import {USER_AUTH_SESSION_TOKEN} from "../constants/userSession";

const withToken = function(config) {
  const token = Cookies.get(USER_AUTH_SESSION_TOKEN);
  if (token) {
    config.headers.Authorization = token;
  }
  return config;
};

export const baseApi = axios.create({
  baseURL: process.env.GATSBY_API_URL,
});

// Authenticated routes
baseApi.interceptors.request.use(withToken);
