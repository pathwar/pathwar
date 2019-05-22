import { 
	GET_USER_TEAMS_SUCCESS,
	GET_USER_TEAMS_FAILED, 
	SET_TEAMS_LIST, 
	SET_TEAMS_LIST_FAILED,
	JOIN_TEAM_SUCCESS,
	JOIN_TEAM_FAILED
} from "../constants/actionTypes"
import { getAllTeams, getUserTeams, joinTeam as joinTeamCall } from "../api/teams"
import { fetchTeamTournaments as fetchTeamTournamentsAction } from "./tournaments";

export const fetchUserTeams = (userID) => async dispatch => {

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

		dispatch(fetchTeamTournamentsAction(lastActiveTeam.metadata.id))
		
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

export const joinTeam = (userID, teamID) => async dispatch => {
	try {
		const response = await joinTeamCall(userID, teamID);
		dispatch({
			type: JOIN_TEAM_SUCCESS,
			payload: response.data
		});

		dispatch(fetchTeamsList())
		dispatch(fetchUserTeams(userID))
	}
	catch (error) {
		dispatch({ type: JOIN_TEAM_FAILED, payload: { error } });
		alert("Join team failed, please try again!")
	}
}