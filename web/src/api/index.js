import axios from "axios";

const withToken = function(config) {
	const token = "jwtToken";
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