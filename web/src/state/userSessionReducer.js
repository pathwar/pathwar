import {
  LOGOUT,
  LOGIN_FAILED,
  SET_USER_SESSION,
  SET_KEYCLOAK_SESSION,
} from "../constants/actionTypes";

const initialState = {
  session: {
    error: undefined,
    fetching: false,
    isAuthenticated: false,
    activeUserSession: undefined,
    activeKeycloakSession: undefined,
  },
};

export default function userSessionReducer(
  state = initialState.session,
  action
) {
  switch (action.type) {
    case LOGIN_FAILED:
      return {
        ...state,
        fetching: false,
        activeKeycloakSession: undefined,
        isAuthenticated: false,
        error: action.payload.error,
      };

    case LOGOUT:
      return {
        ...state,
        fetching: false,
        activeKeycloakSession: undefined,
        activeUserSession: undefined,
        isAuthenticated: false,
        error: undefined,
      };

    case SET_KEYCLOAK_SESSION:
      return {
        ...state,
        fetching: false,
        activeKeycloakSession: action.payload.keycloakInstance,
        isAuthenticated: action.payload.authenticated,
      };

    case SET_USER_SESSION:
      return {
        ...state,
        fetching: false,
        activeUserSession: action.payload.activeUserSession,
      };

    default:
      return state;
  }
}
