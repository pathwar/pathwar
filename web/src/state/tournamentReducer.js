import {
  GET_ALL_TOURNAMENTS_SUCCESS,
  SET_ACTIVE_TOURNAMENT,
  GET_ALL_TOURNAMENT_TEAMS_SUCCESS,
  SET_CHALLENGES_LIST,
  SET_ACTIVE_TOURNAMENT
} from '../constants/actionTypes';

const initialState = {
  tournaments: {
    error: null,
    allTournaments: null,
    allTeamTournaments: null,
    activeTournament: null,
    allTeamsOnTournament: null,
    activeChallenges: null
  }
};

export default function tournamentReducer(state = initialState.tournaments, action) {

  switch (action.type) {
    case GET_ALL_TOURNAMENTS_SUCCESS:

      return {
        ...state,
        allTournaments: action.payload.allTournaments
      }

    case GET_ALL_TOURNAMENT_TEAMS_SUCCESS:
      return {
        ...state,
        allTeamsOnTournament: action.payload.allTeams
      }

    case SET_ACTIVE_TOURNAMENT:
      return {
        ...state,
        activeTournament: action.payload.activeTournament
      }

    case SET_CHALLENGES_LIST:

      return {
        ...state,
        activeChallenges: action.payload.activeChallenges
      };

    default:
      return state;
  }
}
