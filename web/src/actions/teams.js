import { 
	GET_USER_TEAMS,
	GET_USER_TEAMS_SUCCESS,
	GET_USER_TEAMS_FAILED, 
	SET_TEAMS_LIST, 
	SET_TEAMS_LIST_FAILED 
} from "../constants/actionTypes"
import { getAllTeams, getUserTeams } from "../api/teams"

export const fetchUserTeams = (userID) => async dispatch => {
	dispatch({
		type: GET_USER_TEAMS
	})

	try {
		const response = await getUserTeams(userID);
		const teams = response.data.items;
		const lastActiveTeam = teams.find((team) => team.lastActive);

		dispatch({
			type: GET_USER_TEAMS_SUCCESS,
			payload: { 
				userTeamsList: teams,
				lastActiveTeam: lastActiveTeam
			}
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
			payload: { allTeamsList: response.data.items }
		});
	} catch (error) {
		dispatch({ type: SET_TEAMS_LIST_FAILED, payload: { error } });
	}
};