import { 
	PERFORM_LOGIN, 
	LOGIN_FAILED,
	SET_USER_SESSION
} from "../constants/actionTypes"
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

		history.push("/dashboard");
	} catch (error) {
		dispatch({ type: LOGIN_FAILED, payload: { error } });
	}
};