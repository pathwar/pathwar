import { SET_LEVELS_LIST, SET_LEVELS_LIST_FAILED } from "../constants/actionTypes"
import { getLevels } from "../api/levels"

export const fetchLevels = () => async dispatch => {
	try {
		const response = await getLevels();
		dispatch({
            type: SET_LEVELS_LIST,
            payload: response.data.items
		});
	} catch (error) {
		dispatch({ type: SET_LEVELS_LIST_FAILED, payload: { error } });
	}
};