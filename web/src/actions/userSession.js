import { SET_USER_SESSION, SET_USER_SESSION_FAILED } from "../constants/actionTypes"
import { getUserSession } from "../api/userSession"

export const fetchUserSession = () => async dispatch => {
	try {
		const response = await getUserSession();

		dispatch({
			type: SET_USER_SESSION,
			payload: response.data
		});
		
	} catch (error) {
		dispatch({ type: SET_USER_SESSION_FAILED, payload: { error } });
	}
};