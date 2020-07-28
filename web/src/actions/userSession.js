/* eslint-disable no-unused-vars */
import Cookies from "js-cookie";
import { toast } from "react-toastify";
import {
  LOGIN_FAILED,
  SET_USER_SESSION,
  SET_USER_SESSION_FAILED,
  SET_KEYCLOAK_SESSION,
  LOGOUT,
  DELETE_ACCOUNT_FAILED,
  DELETE_ACCOUNT_SUCCESS,
} from "../constants/actionTypes";
import { USER_SESSION_TOKEN_NAME } from "../constants/userSession";
import { getUserSession, deleteUserAccount } from "../api/userSession";
import { setActiveOrganization as setActiveOrganizationAction } from "./organizations";
import {
  setActiveSeason as setActiveSeasonAction,
  fetchPreferences as fetchPreferencesAction,
  setActiveTeam as setActiveTeamAction,
} from "./seasons";

export const logoutUser = () => async dispatch => {
  dispatch({
    type: LOGOUT,
  });
};

export const setUserSession = activeUserSession => async dispatch => {
  dispatch({
    type: SET_USER_SESSION,
    payload: { activeUserSession },
  });
};

export const fetchUserSession = postPreferences => async dispatch => {
  try {
    const userSessionResponse = await getUserSession();
    const { data: userSessionData } = userSessionResponse;

    const activeSeasonId = userSessionData.user.active_season_id;

    dispatch(setUserSession(userSessionData));

    if (postPreferences) {
      dispatch(fetchPreferencesAction(activeSeasonId));
    }

    if (activeSeasonId) {
      const activeSeason = userSessionData.seasons.find(
        item => item.season.id === activeSeasonId
      );
      dispatch(setActiveSeasonAction(activeSeason.season));
      dispatch(setActiveTeamAction(activeSeason.team));
      dispatch(setActiveOrganizationAction(activeSeason.team.organization));
    }
  } catch (error) {
    dispatch({ type: SET_USER_SESSION_FAILED, payload: { error } });
  }
};

export const setKeycloakSession = (
  keycloakInstance,
  authenticated
) => async dispatch => {
  try {
    dispatch({
      type: SET_KEYCLOAK_SESSION,
      payload: {
        keycloakInstance: keycloakInstance,
        authenticated: authenticated,
      },
    });

    Cookies.set(USER_SESSION_TOKEN_NAME, keycloakInstance.token);
    dispatch(fetchUserSession(true));
  } catch (error) {
    dispatch({ type: LOGIN_FAILED, payload: { error } });
  }
};

export const deleteAccount = reason => async dispatch => {
  try {
    const response = await deleteUserAccount(reason);
    dispatch({
      type: DELETE_ACCOUNT_SUCCESS,
      payload: { activeChallenges: response.data.items },
    });
    toast.success("Delete account SUCCESS!");
  } catch (error) {
    dispatch({ type: DELETE_ACCOUNT_FAILED, payload: { error } });
    toast.error("Delete account FAILED!");
  }
};
