import { 
  PERFORM_LOGIN, 
  LOGIN_SUCCESS,
  LOGIN_FAILED,
  SET_USER_SESSION,
  SET_USER_SESSION_FAILED
} from '../constants/actionTypes';

const initialState = {
  session: {
    error: null,
    fetching: false,
    activeUser: null
  }
};

export default function sessionReducer(state = initialState.session, action) {

  switch (action.type) {
    case PERFORM_LOGIN:
      return {
        ...state,
        fetching: true
      };

    case LOGIN_SUCCESS:
      return {
        ...state,
        fetching: false,
        error: null
    };

    case LOGIN_FAILED:
      return {
        ...state,
        fetching: false,
        error: action.payload
    } ;
        

    case SET_USER_SESSION:
      
      return {
        ...state,
        activeUser: action.payload
      };

    case SET_USER_SESSION_FAILED:
      
      return {
        ...state,
        activeUser: null
      };

    default:
      return state;
  }
}
