/* eslint-disable no-unused-vars */
import Cookies from "js-cookie";
import { 
	PERFORM_LOGIN, 
	LOGIN_FAILED,
	SET_USER_SESSION,
	PING_USER_SUCCESS,
	PING_USER_FAILED
} from "../constants/actionTypes"
import { USER_SESSION_TOKEN_NAME } from "../constants/userSession";
import { performLogin, pingUser } from "../api/userSession"

export const performLoginAction = (email, password) => async dispatch => {
	dispatch({
		type: PERFORM_LOGIN
	});

	try {
		const response = await performLogin(email, password);
		
		dispatch({
			type: SET_USER_SESSION,
			payload: { activeUser: response.data }
		});

		Cookies.set(USER_SESSION_TOKEN_NAME, response.data.token)
	} catch (error) {
		dispatch({ type: LOGIN_FAILED, payload: { error } });
	}
};

export const pingUserAction = () => async dispatch => {

	try {
		const response = await pingUser();
		dispatch({
			type: PING_USER_SUCCESS,
			payload: { 
				isAuthenticated: response.data.isAuthenticated,
				activeUser: response.data.user
			}
		});

		Cookies.set(USER_SESSION_TOKEN_NAME, response.data.token)
	} catch (error) {
		dispatch({ type: PING_USER_FAILED, payload: { error } });
	}
};