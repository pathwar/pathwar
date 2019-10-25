import {
  GET_ALL_TOURNAMENTS_SUCCESS,
  SET_ACTIVE_TOURNAMENT,
  GET_ALL_TOURNAMENT_TEAMS_SUCCESS,
  SET_CHALLENGES_LIST,
  GET_CHALLENGE_DETAILS_SUCCESS
} from '../constants/actionTypes';

const initialState = {
  tournaments: {
    error: null,
    allTournaments: null,
    allTeamTournaments: null,
    activeTournament: null,
    allTeamsOnTournament: null,
    activeChallenges: null,
    challengeInDetail: null
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

    case GET_CHALLENGE_DETAILS_SUCCESS:
      return {
        ...state,
        challengeInDetail: action.payload.challenge
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
