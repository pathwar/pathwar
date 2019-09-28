/* eslint-disable no-unused-vars */
import Cookies from "js-cookie";
import {
  LOGIN_FAILED,
  SET_USER_SESSION,
  SET_KEYCLOAK_SESSION,
  PING_USER_SUCCESS,
  PING_USER_FAILED,
  LOGOUT
} from "../constants/actionTypes"
import { USER_SESSION_TOKEN_NAME } from "../constants/userSession";
import { getUserSession } from "../api/userSession"
import { setActiveTeam as setActiveTeamAction } from "./teams";
import { setActiveTournament as setActiveTournamentAction } from "./tournaments"

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

export const setKeycloakSession = (keycloakInstance, authenticated) => async dispatch => {

  try {

    dispatch({
      type: SET_KEYCLOAK_SESSION,
      payload: { keycloakInstance: keycloakInstance, authenticated: authenticated }
    });

    Cookies.set(USER_SESSION_TOKEN_NAME, keycloakInstance.token);

    const userSessionResponse = await getUserSession();

    dispatch(setUserSession(userSessionResponse.data))

    console.log("AI ELA >>>", userSessionResponse)

    // dispatch(setActiveTeamAction(lastActiveTeam))
    // dispatch(setActiveTournamentAction(defaultTournament));


    //TODO: Verify how we can retrieve the team and tournament to be set on first load after login.

    // const lastActiveTeam = {
    //   "metadata": {
    //     "id": "2dNuq9SuRvKcVHDOfNJSFQ==",
    //     "created_at": "2019-04-25T11:41:24Z",
    //     "updated_at": "2019-04-25T11:41:24Z"
    //   },
    //   "name": "chartreuse",
    //   "gravatar_url": "http://www.internationalmonetize.io/harness/communities",
    //   "locale": "fr_FR",
    //   "last_active": true
    //   }
    // const defaultTournament = {
    //   "metadata": {
    //     "id": "6ugtOJHGRrumcmufKZmTTQ==",
    //     "created_at": "2019-04-25T12:32:36Z",
    //     "updated_at": "2019-04-25T12:32:36Z"
    //   },
    //   "name": "tumblr",
    //   "status": "Started",
    //   "visibility": "Public",
    //   "is_default": true
    // }

  } catch (error) {
    dispatch({ type: LOGIN_FAILED, payload: { error } });
  }
};
