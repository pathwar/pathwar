import { 
  LOGIN_FAILED,
  SET_USER_SESSION,
  PING_USER_SUCCESS,
  PING_USER_FAILED
} from '../constants/actionTypes';

const initialState = {
  session: {
    error: null,
    fetching: false,
    activeSession: null,
    isAuthenticated: false
  }
};

export default function userSessionReducer(state = initialState.session, action) {

  switch (action.type) {

    case LOGIN_FAILED:
      return {
        ...state,
        fetching: false,
        activeSession: null,
        isAuthenticated: false,
        error: action.payload.error
    } ;

    case SET_USER_SESSION:
      return {
        ...state,
        fetching: false,
        activeSession: action.payload.activeSession,
        isAuthenticated: action.payload.authenticated
      };

    case PING_USER_SUCCESS:
      return {
        ...state,
        isAuthenticated: action.payload.authenticated,
        activeSession: action.payload.activeSession
      }

    case PING_USER_FAILED:
      return {
        ...state,
        isAuthenticated: false,
        error: action.payload.error
      }

    default:
      return state;
  }
}
