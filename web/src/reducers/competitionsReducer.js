import { SET_LEVELS_LIST } from '../constants/actionTypes';

const initialState = {
  competition: {
    levels: null
  }
};

export default function competitionsReducer(state = initialState.session, action) {

  switch (action.type) {
    case SET_LEVELS_LIST:
      
      return {
        ...state,
        levels: action.payload
      };

    default:
      return state;
  }
}
