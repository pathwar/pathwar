import { 
	GET_USER_COMPETITIONS_SUCCESS, 
	GET_USER_COMPETITIONS_FAILED,
	SET_LEVELS_LIST, 
	SET_LEVELS_LIST_FAILED 
} from "../constants/actionTypes"
import { getCompetitions, getLevels } from "../api/competitions"

export const fetchCompetitions = () => async dispatch => {
	try {
		const response = await getCompetitions();
		dispatch({
			type: GET_USER_COMPETITIONS_SUCCESS,
			payload: { allCompetitions: response.data.items }
		});

	} catch (error) {
		dispatch({
			type: GET_USER_COMPETITIONS_FAILED,
			payload: { error }
		});
	}
}

export const fetchLevels = (competitionId) => async dispatch => {
	try {
		const response = await getLevels(competitionId);
		dispatch({
			type: SET_LEVELS_LIST,
			payload: response.data.items
		});
	} catch (error) {
		dispatch({ type: SET_LEVELS_LIST_FAILED, payload: { error } });
	}
};