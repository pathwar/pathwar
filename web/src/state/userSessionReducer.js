import { 
  PERFORM_LOGIN, 
  LOGIN_FAILED,
  SET_USER_SESSION,
  PING_USER_SUCCESS,
  PING_USER_FAILED
} from '../constants/actionTypes';

const initialState = {
  session: {
    error: null,
    fetching: false,
    activeUser: null,
    isAuthenticated: false
  }
};

export default function userSessionReducer(state = initialState.session, action) {

  switch (action.type) {
    case PERFORM_LOGIN:
      return {
        ...state,
        error: null,
        fetching: true
      };

    case LOGIN_FAILED:
      return {
        ...state,
        fetching: false,
        activeUser: null,
        isAuthenticated: false,
        error: action.payload.error
    } ;

    case SET_USER_SESSION:
      return {
        ...state,
        fetching: false,
        activeUser: action.payload.activeUser,
        isAuthenticated: true
      };

    case PING_USER_SUCCESS:
      return {
        ...state,
        isAuthenticated: action.payload.isAuthenticated,
        activeUser: action.payload.activeUser
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
