/* eslint-disable no-unused-vars */
import Cookies from "js-cookie";
import {
  LOGIN_FAILED,
  SET_USER_SESSION,
  SET_USER_SESSION_FAILED,
  SET_KEYCLOAK_SESSION,
  LOGOUT
} from "../constants/actionTypes"
import { USER_SESSION_TOKEN_NAME } from "../constants/userSession";
import { getUserSession } from "../api/userSession"
import { setActiveTeam as setActiveTeamAction } from "./teams";
import {
  setActiveTournament as setActiveTournamentAction,
  fetchPreferences as fetchPreferencesAction
} from "./tournaments"

export const logoutUser = () => async dispatch => {
  dispatch({
    type: LOGOUT
  })
}

export const setUserSession = (activeUserSession) => async dispatch => {
  dispatch({
    type: SET_USER_SESSION,
    payload: { activeUserSession }
  })
}

export const fetchUserSession = (postPreferences) => async dispatch => {

  try {
    const userSessionResponse = await getUserSession();
    const userSessionData = userSessionResponse.data;
    const defaultTournamentSet = userSessionData.tournaments.find((item) => item.tournament.is_default);
    const defaultTeamSet = userSessionData.tournaments.find((item) => item.team.is_default);

    const defaultTournament = defaultTournamentSet.tournament;
    const defaultTeam = defaultTeamSet.team;

    const activeTournamentId = userSessionData.user.active_tournament_id

    dispatch(setUserSession(userSessionData))

    if (postPreferences) {
      dispatch(fetchPreferencesAction(defaultTournament.id))
    }

    if (activeTournamentId) {
      const activeTournament = userSessionData.tournaments.find((item) => item.tournament.id === activeTournamentId);
      dispatch(setActiveTournamentAction(activeTournament.tournament));
      dispatch(setActiveTeamAction(defaultTeam));
    }

  }
  catch(error) {
    dispatch({ type: SET_USER_SESSION_FAILED, payload: { error } });
  }
}

export const setKeycloakSession = (keycloakInstance, authenticated) => async dispatch => {

  try {

    dispatch({
      type: SET_KEYCLOAK_SESSION,
      payload: { keycloakInstance: keycloakInstance, authenticated: authenticated }
    });

    Cookies.set(USER_SESSION_TOKEN_NAME, keycloakInstance.token);
    dispatch(fetchUserSession(true))

  } catch (error) {
    dispatch({ type: LOGIN_FAILED, payload: { error } });
  }
};
