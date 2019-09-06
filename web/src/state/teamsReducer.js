import {
  SET_ACTIVE_TEAM,
  GET_USER_TEAMS,
  SET_TEAMS_LIST,
  GET_USER_TEAMS_SUCCESS
} from '../constants/actionTypes';

const initialState = {
  teams: {
      fetchingUserTeams: null,
      userTeamsList: null,
      activeTeam: null,
      allTeamsList: null
  }
};

export default function teamsReducer(state = initialState.teams, action) {

  switch (action.type) {
    case SET_ACTIVE_TEAM:
      return {
        ...state,
        activeTeam: action.payload.team
      }

    case GET_USER_TEAMS:
      return {
        ...state,
        fetchingUserTeams: true
      }

    case GET_USER_TEAMS_SUCCESS:
      return {
        ...state,
        fetchingUserTeams: false,
        userTeamsList: action.payload.userTeamsList
      }

    case SET_TEAMS_LIST:
      return {
        ...state,
        allTeamsList: action.payload.allTeamsList
      };

    default:
      return state;
  }
}
