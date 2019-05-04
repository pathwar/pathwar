import { SET_SESSION, SET_SESSION_FAILED } from "../constants/actionTypes"
import { getUserSession } from "../api/session"

export const fetchUserSession = () => async dispatch => {
	try {
		const response = await getUserSession();
		dispatch({
            type: SET_SESSION,
            payload: response.data
		});
	} catch (error) {
		dispatch({ type: SET_SESSION_FAILED, payload: { error } });
	}
};