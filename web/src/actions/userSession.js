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
  VALIDATE_COUPON_SUCCESS,
  VALIDATE_COUPON_FAILED,
  SET_AUTH_SESSION,
} from "../constants/actionTypes";
import {USER_AUTH_SESSION_TOKEN, USER_SESSION_TOKEN_NAME} from "../constants/userSession";
import {
  getUserSession,
  deleteUserAccount,
  postCouponValidation,
  fetchAccessToken,
} from "../api/userSession";
import {
  fetchUserOrganizationsInvitations,
  setActiveOrganization as setActiveOrganizationAction,
  setUserOrganizationsList
} from "./organizations";
import {
  setActiveSeason as setActiveSeasonAction,
  fetchPreferences as fetchPreferencesAction,
  setActiveTeam as setActiveTeamAction,
  fetchUserTeamsInvitations,
  fetchUserSeasons,
} from "./seasons";
import dispatchFireworks from "../utils/fireworks-dispatcher";
import {apiErrorsCode} from "../constants/apiErrorsCode";

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
  const browser = typeof window !== "undefined" && window;

  try {
    const userSessionResponse = await getUserSession();
    const { data: userSessionData } = userSessionResponse;
    const { user, organizations, organization_invites, team_invites, seasons } = userSessionData;
    const activeSeasonId = user.active_season_id;

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
      dispatch(fetchUserSeasons(seasons));
      dispatch(setActiveOrganizationAction(activeSeason.team.organization));
      dispatch(setUserOrganizationsList(organizations));
      dispatch(fetchUserOrganizationsInvitations(organization_invites));
      dispatch(fetchUserTeamsInvitations(team_invites));
    }

    if (browser) {
      setTimeout(() => {
        window.$crisp.push(["set", "user:email", [user.email]]);
        window.$crisp.push(["set", "user:nickname", [user.slug]]);
        window.$crisp.push(["set", "user:avatar", [user.gravatar_url]]);
        window.$crisp.push([
          "set",
          "session:data",
          [
            [
              ["user_id", user.id],
              ["active_team_member_id", user.active_team_member_id],
              ["active_season_id", user.active_season_id],
            ],
          ],
        ]);
      }, 2000);
    }
  } catch (error) {
    dispatch({ type: SET_USER_SESSION_FAILED, payload: { error } });
  }
};

//TODO: receive authenticated parameter from the provider
export const setAuthSession = () => async dispatch => {
  try {
    const response = await fetchAccessToken();
    dispatch({
      type: SET_AUTH_SESSION,
      payload: {
        token: response.data.token,
        authenticated: true,
      }
    })
    Cookies.set(USER_AUTH_SESSION_TOKEN, response.data.token);
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

//Coupon Actions
export const fetchCouponValidation = (hash, teamID) => async dispatch => {
  try {
    const response = await postCouponValidation(hash, teamID);
    dispatch({
      type: VALIDATE_COUPON_SUCCESS,
      payload: { team: response.data.coupon_validation.team },
    });
    toast.success(`Coupon validation success!`);
    dispatchFireworks();
  } catch (error) {
    dispatch({
      type: VALIDATE_COUPON_FAILED,
      payload: { error: error.response },
    });
    const errorText = apiErrorsCode.get(error.response.data.message.split(':')[0]) ?
      apiErrorsCode.get(error.response.data.message.split(':')[0]) : 'An error occurred';
    toast.error(`${errorText}`);
  }
};
