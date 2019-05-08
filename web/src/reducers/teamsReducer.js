import { SET_TEAMS_LIST } from '../constants/actionTypes';

const initialState = {
  teams: {
      teamsList: null
  }
};

export default function sessionReducer(state = initialState.teams, action) {

  switch (action.type) {
    case SET_TEAMS_LIST:
      
      return {
        ...state,
        teamsList: action.payload
      };

    default:
      return state;
  }
}
