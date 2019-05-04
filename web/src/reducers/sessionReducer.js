import { SET_SESSION } from '../constants/actionTypes';

const initialState = {
  session: {
    activeSession: null
  }
};

export default function sessionReducer(state = initialState.session, {type, payload}) {

  switch (type) {
    case SET_SESSION:
      
      return {
        ...state,
        activeSession: payload.session
      };

    default:
      return state;
  }
}
