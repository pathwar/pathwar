import { SET_TEAMS_LIST, SET_TEAMS_LIST_FAILED } from "../constants/actionTypes"
import { getTeams } from "../api/teams"

export const fetchTeamsList = () => async dispatch => {
	try {
		const response = await getTeams();
		dispatch({
            type: SET_TEAMS_LIST,
            payload: response.data.items
		});
	} catch (error) {
		dispatch({ type: SET_TEAMS_LIST_FAILED, payload: { error } });
	}
};