import {
  GET_ALL_SEASONS_SUCCESS,
  SET_ACTIVE_SEASON,
  GET_ALL_SEASON_TEAMS_SUCCESS,
  SET_CHALLENGES_LIST,
  GET_CHALLENGE_DETAILS_SUCCESS,
  GET_TEAM_DETAILS_SUCCESS,
  SET_ACTIVE_TEAM
} from '../constants/actionTypes';

const initialState = {
  seasons: {
    error: null,
    allSeasons: null,
    activeSeason: null,
    activeTeamInSeason: null,
    activeTeam: null,
    teamInDetail: null,
    allTeamsOnSeason: null,
    activeChallenges: null,
    challengeInDetail: undefined,
  }
};

export default function seasonReducer(state = initialState.seasons, action) {

  switch (action.type) {
    case GET_ALL_SEASONS_SUCCESS:
      return {
        ...state,
        allSeasons: action.payload.allSeasons
      }

    case GET_ALL_SEASON_TEAMS_SUCCESS:
      return {
        ...state,
        allTeamsOnSeason: action.payload.allTeams
      }

    case SET_ACTIVE_SEASON:
      return {
        ...state,
        activeSeason: action.payload.activeSeason
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

    case GET_TEAM_DETAILS_SUCCESS:
      return {
        ...state,
        teamInDetail: action.payload.team
      }

    case SET_ACTIVE_TEAM:
      const { payload: { team } } = action;
      const { allTeamsOnSeason } = state;

      return {
        ...state,
        activeTeam: team,
        activeTeamInSeason: allTeamsOnSeason && allTeamsOnSeason.some(item => item.id === team.id)
      }

    default:
      return state;
  }
}
