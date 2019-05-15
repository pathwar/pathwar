import axios from "axios";
import Cookies from "js-cookie";
import { USER_SESSION_TOKEN_NAME } from "../constants/userSession";

const withToken = function(config) {
	const token = Cookies.get(USER_SESSION_TOKEN_NAME);
	if (token) {
		config.headers.Authorization = token;
	}
	return config;
};

export const baseApi = axios.create({
	baseURL: process.env.API_URL
});

// Authenticated routes
baseApi.interceptors.request.use(withToken);