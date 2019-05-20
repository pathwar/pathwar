import { 
	GET_TOURNAMENTS_SUCCESS, 
	GET_TOURNAMENTS_FAILED,
	SET_LEVELS_LIST, 
	SET_LEVELS_LIST_FAILED 
} from "../constants/actionTypes"
import { getTournaments, getLevels } from "../api/tournaments"

export const fetchTournaments = () => async dispatch => {
	try {
		const response = await getTournaments();
		dispatch({
			type: GET_TOURNAMENTS_SUCCESS,
			payload: { allTournaments: response.data.items }
		});

	} catch (error) {
		dispatch({
			type: GET_TOURNAMENTS_FAILED,
			payload: { error }
		});
	}
}

export const fetchLevels = (competitionId) => async dispatch => {
	try {
		const response = await getLevels(competitionId);
		dispatch({
			type: SET_LEVELS_LIST,
			payload: { activeLevels: response.data.items }
		});
	} catch (error) {
		dispatch({ type: SET_LEVELS_LIST_FAILED, payload: { error } });
	}
};