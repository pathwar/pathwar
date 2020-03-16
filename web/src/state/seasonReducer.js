import { clone, update, findIndex, propEq } from "ramda";

import {
  GET_ALL_SEASONS_SUCCESS,
  SET_ACTIVE_SEASON,
  GET_ALL_SEASON_TEAMS_SUCCESS,
  SET_CHALLENGES_LIST,
  GET_CHALLENGE_DETAILS_SUCCESS,
  GET_TEAM_DETAILS_SUCCESS,
  SET_ACTIVE_TEAM,
  CLOSE_CHALLENGE_SUCCESS
} from '../constants/actionTypes';

const initialState = {
  seasons: {
    error: undefined,
    allSeasons: undefined,
    activeSeason: undefined,
    activeTeamInSeason: undefined,
    activeTeam: undefined,
    teamInDetail: undefined,
    allTeamsOnSeason: undefined,
    activeChallenges: undefined,
    challengeInDetail: undefined,
  }
};

export default function seasonReducer(state = initialState.seasons, action) {
  const { challengeInDetail, allTeamsOnSeason } = state;

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
      const { payload: { activeChallenges } } = action;
      return {
        ...state,
        activeChallenges: action.payload.activeChallenges,
        challengeInDetail: activeChallenges && challengeInDetail && activeChallenges.find(item => item.id === challengeInDetail.id)
      };

    case GET_TEAM_DETAILS_SUCCESS:
      return {
        ...state,
        teamInDetail: action.payload.team
      }

    case SET_ACTIVE_TEAM:
      const { payload: { team } } = action;

      return {
        ...state,
        activeTeam: team,
        activeTeamInSeason: allTeamsOnSeason && allTeamsOnSeason.some(item => item.id === team.id)
      }

    case CLOSE_CHALLENGE_SUCCESS:
      const { payload: { subscription: { challenge_subscription } } } = action;

      const challengeInDetailClone = clone(challengeInDetail);
      const { subscriptions } = challengeInDetailClone;

      const subscriptionIndex = findIndex(propEq("id", challenge_subscription.id))(subscriptions);
      const updatedSubscriptions = update(subscriptionIndex, challenge_subscription, subscriptions);
      challengeInDetailClone.subscriptions = updatedSubscriptions;

      return {
        ...state,
        challengeInDetail: challengeInDetailClone
      }

    default:
      return state;
  }
}
