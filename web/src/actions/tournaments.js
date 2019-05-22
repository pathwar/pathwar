import { 
	GET_TOURNAMENTS_SUCCESS, 
	GET_TOURNAMENTS_FAILED,
	SET_DEFAULT_TOURNAMENT,
	SET_ACTIVE_TOURNAMENT,
	SET_LEVELS_LIST, 
	SET_LEVELS_LIST_FAILED 
} from "../constants/actionTypes"
import { getTeamTournaments, getLevels } from "../api/tournaments"

export const fetchTeamTournaments = (teamID) => async dispatch => {
	try {
		const response = await getTeamTournaments(teamID);
		const allTournaments = response.data.items;
		const defaultTournament = allTournaments.find((tournament) => tournament.is_default)
		
		dispatch({
			type: GET_TOURNAMENTS_SUCCESS,
			payload: { allTournaments: allTournaments }
		});

		if (defaultTournament) {
			dispatch({
				type: SET_DEFAULT_TOURNAMENT,
				payload: { defaultTournament: defaultTournament }
			});
			
			dispatch({
				type: SET_ACTIVE_TOURNAMENT,
				payload: { activeTournament: defaultTournament }
			});

			dispatch(fetchLevels(defaultTournament.metadata.id))
		}

	} catch (error) {
		dispatch({
			type: GET_TOURNAMENTS_FAILED,
			payload: { error }
		});
	}
}

export const fetchLevels = (tournamentID) => async dispatch => {
	try {
		const response = await getLevels(tournamentID);
		dispatch({
			type: SET_LEVELS_LIST,
			payload: { activeLevels: response.data.items }
		});
	} catch (error) {
		dispatch({ type: SET_LEVELS_LIST_FAILED, payload: { error } });
	}
};