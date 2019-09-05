import {
  GET_ALL_TOURNAMENTS_SUCCESS,
  GET_ALL_TOURNAMENTS_FAILED,
  GET_TEAM_TOURNAMENTS_SUCCESS,
  GET_TEAM_TOURNAMENTS_FAILED,
  SET_DEFAULT_TOURNAMENT,
  SET_ACTIVE_TOURNAMENT,
  SET_ACTIVE_TOURNAMENT_FAILED,
  SET_LEVELS_LIST,
  SET_LEVELS_LIST_FAILED
} from "../constants/actionTypes"
import {
  getAllTournaments,
  getTeamTournaments,
  getLevels
} from "../api/tournaments"

export const setActiveTournament = (tournamentData) => async dispatch => {
  try {
    // const response = await setTournamentActiveCall(tournamentData.metadata.id)
    // if (response.status === 200) {
      dispatch({
        type: SET_ACTIVE_TOURNAMENT,
        payload: { activeTournament: tournamentData }
      });
    // }
  }
  catch(error) {
    dispatch({ type: SET_ACTIVE_TOURNAMENT_FAILED, payload: { error }});
    alert("Set tournament active failed, please try again!")
  }

}

export const setDefaultTournament = (tournamentData) => async dispatch => {
  dispatch({
    type: SET_DEFAULT_TOURNAMENT,
    payload: { defaultTournament: tournamentData }
  });
}

export const fetchAllTournaments = () => async dispatch => {
  try {
    const response = await getAllTournaments();
    const allTournaments = response.data.items;

    dispatch({
      type: GET_ALL_TOURNAMENTS_SUCCESS,
      payload: { allTournaments: allTournaments }
    })
  } catch (error) {
    dispatch({ type: GET_ALL_TOURNAMENTS_FAILED, payload: { error } });
  }
}

export const fetchTeamTournaments = (teamID) => async dispatch => {
  try {
    const response = await getTeamTournaments(teamID);
    const allTeamTournaments = response.data.items;
    const lastActiveTournament = allTeamTournaments.find((tournament) => tournament.last_active)
    const defaultTournament = allTeamTournaments.find((tournament) => tournament.is_default)

    dispatch({
      type: GET_TEAM_TOURNAMENTS_SUCCESS,
      payload: { allTeamTournaments: allTeamTournaments }
    });

    if (lastActiveTournament === defaultTournament) {
      dispatch(setActiveTournament(lastActiveTournament));
    } else if (!lastActiveTournament && defaultTournament) {
      dispatch(setDefaultTournament(defaultTournament));
      dispatch(setActiveTournament(defaultTournament));
    }

  } catch (error) {
    dispatch({
      type: GET_TEAM_TOURNAMENTS_FAILED,
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
