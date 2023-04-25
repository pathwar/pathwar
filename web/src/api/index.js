import axios from "axios";
import Cookies from "js-cookie";
import {USER_AUTH_SESSION_TOKEN, USER_SESSION_TOKEN_NAME} from "../constants/userSession";

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

export const baseAuthApi = axios.create({
  baseURL: process.env.AUTH_SERVICE_URL,
});

// Authenticated routes
baseApi.interceptors.request.use(withToken);
