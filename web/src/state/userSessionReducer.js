import {
  LOGOUT,
  LOGIN_FAILED,
  SET_USER_SESSION,
  VALIDATE_COUPON_SUCCESS,
  VALIDATE_CHALLENGE_SUCCESS,
  BUY_CHALLENGE_SUCCESS,
  SET_AUTH_SESSION,
} from "../constants/actionTypes";

const initialState = {
  session: {
    error: undefined,
    fetching: false,
    isAuthenticated: false,
    activeUserSession: undefined,
    cash: undefined,
    accessToken: undefined,
  },
};

export default function userSessionReducer(
  state = initialState.session,
  action
) {
  const { activeUserSession } = action.payload || {};

  switch (action.type) {
    case LOGIN_FAILED:
      return {
        ...state,
        fetching: false,
        isAuthenticated: false,
        error: action.payload.error,
        accessToken: undefined,
      };

    case LOGOUT:
      return {
        ...state,
        fetching: false,
        activeUserSession: undefined,
        isAuthenticated: false,
        error: undefined,
        accessToken: undefined,
      };

      case SET_AUTH_SESSION:
      return {
        ...state,
        accessToken: action.payload.token,
        isAuthenticated: action.payload.authenticated,
      }

    case SET_USER_SESSION:
      return {
        ...state,
        fetching: false,
        activeUserSession: activeUserSession,
        cash: activeUserSession.user.active_team_member.team.cash,
      };

    case VALIDATE_COUPON_SUCCESS:
      return {
        ...state,
        cash: action.payload.team.cash,
      };

    case VALIDATE_CHALLENGE_SUCCESS:
      return {
        ...state,
        cash: action.payload.challengeSubscription.team.cash,
      };

    case BUY_CHALLENGE_SUCCESS:
      return {
        ...state,
        cash: action.payload.challengeSubscription.team.cash,
      };

    default:
      return state;
  }
}
