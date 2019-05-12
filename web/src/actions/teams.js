import { 
	GET_USER_TEAMS,
	GET_USER_TEAMS_SUCCESS,
	GET_USER_TEAMS_FAILED, 
	SET_TEAMS_LIST, 
	SET_TEAMS_LIST_FAILED 
} from "../constants/actionTypes"
import { getAllTeams, getUserTeams } from "../api/teams"

export const fetchUserTeams = (userId) => async dispatch => {
	dispatch({
		type: GET_USER_TEAMS
	})

	try {
		const response = await getUserTeams(userId);
		dispatch({
			type: GET_USER_TEAMS_SUCCESS,
			payload: { userTeams: response.data.items }
		})
	} catch (error) {
		dispatch({
			type: GET_USER_TEAMS_FAILED,
			payload: { error }
		})
	}
}

export const fetchTeamsList = () => async dispatch => {
	try {
		const response = await getAllTeams();
		dispatch({
			type: SET_TEAMS_LIST,
			payload: { teamsList: response.data.items }
		});
	} catch (error) {
		dispatch({ type: SET_TEAMS_LIST_FAILED, payload: { error } });
	}
};