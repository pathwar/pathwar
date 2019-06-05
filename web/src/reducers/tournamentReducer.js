import { 
  GET_ALL_TOURNAMENTS_SUCCESS,
  GET_TEAM_TOURNAMENTS_SUCCESS, 
  SET_LEVELS_LIST, 
  SET_ACTIVE_TOURNAMENT
} from '../constants/actionTypes';

const initialState = {
  tournaments: {
    error: null,
    allTournaments: null,
    allTeamTournaments: null,
    activeTournament: null,
    activeLevels: null
  }
};

export default function tournamentReducer(state = initialState.tournaments, action) {

  switch (action.type) {
    case GET_ALL_TOURNAMENTS_SUCCESS:

      return {
        ...state,
        allTournaments: action.payload.allTournaments
      }

    case GET_TEAM_TOURNAMENTS_SUCCESS:

      return {
        ...state,
        error: null,
        allTeamTournaments: action.payload.allTeamTournaments
      }

    case SET_ACTIVE_TOURNAMENT: 
      return {
        ...state,
        activeTournament: action.payload.activeTournament
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
