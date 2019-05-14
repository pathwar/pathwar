import { 
  GET_USER_COMPETITIONS_SUCCESS, 
  SET_LEVELS_LIST 
} from '../constants/actionTypes';

const initialState = {
  competition: {
    error: null,
    fetchingCompetitions: null,
    allCompetitions: null,
    lastCompetition: null,
    activeCompetition: null,
    levels: null
  }
};

export default function competitionReducer(state = initialState.competition, action) {

  switch (action.type) {
    case GET_USER_COMPETITIONS_SUCCESS:

      return {
        ...state,
        error: null,
        allCompetitions: action.payload.allCompetitions
      }

    case SET_LEVELS_LIST:
      
      return {
        ...state,
        levels: action.payload
      };

    default:
      return state;
  }
}
