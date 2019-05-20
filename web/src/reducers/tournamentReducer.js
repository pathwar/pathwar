import { 
  GET_TOURNAMENTS_SUCCESS, 
  SET_LEVELS_LIST 
} from '../constants/actionTypes';

const initialState = {
  tournaments: {
    error: null,
    allTournaments: null,
    activeTournament: null,
    activeLevels: null
  }
};

export default function tournamentReducer(state = initialState.tournaments, action) {

  switch (action.type) {
    case GET_TOURNAMENTS_SUCCESS:

      return {
        ...state,
        error: null,
        allTournaments: action.payload.allTournaments
      }

    case SET_LEVELS_LIST:
      
      return {
        ...state,
        activeLevels: action.payload.activeLevels
      };

    default:
      return state;
  }
}
