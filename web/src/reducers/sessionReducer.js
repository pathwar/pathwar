import { SET_SESSION } from '../constants/actionTypes';

const initialState = {
  session: {
    activeSession: null
  }
};

export default function sessionReducer(state = initialState.session, action) {

  switch (action.type) {
    case SET_SESSION:
      
      return {
        ...state,
        activeSession: action.payload
      };

    default:
      return state;
  }
}
