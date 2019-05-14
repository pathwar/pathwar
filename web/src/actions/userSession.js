/* eslint-disable no-unused-vars */
import Cookies from "js-cookie";
import { 
	PERFORM_LOGIN, 
	LOGIN_FAILED,
	SET_USER_SESSION
} from "../constants/actionTypes"
import { USER_SESSION_TOKEN_NAME } from "../constants/userSession";
import { history } from "../store/configureStore";
import { performLogin } from "../api/userSession"

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

		// Cookies.set(USER_SESSION_TOKEN_NAME, response.data.token)

		history.push("/dashboard");
	} catch (error) {
		dispatch({ type: LOGIN_FAILED, payload: { error } });
	}
};