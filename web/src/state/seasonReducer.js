/* eslint-disable no-case-declarations */
import { clone, update, findIndex, propEq } from "ramda";

import {
  GET_ALL_SEASONS_SUCCESS,
  SET_ACTIVE_SEASON,
  GET_ALL_SEASON_TEAMS_SUCCESS,
  SET_CHALLENGES_LIST,
  GET_CHALLENGE_DETAILS_SUCCESS,
  CLEAN_CHALLENGE_DETAIL,
  GET_TEAM_DETAILS_SUCCESS,
  SET_ACTIVE_TEAM,
  CLOSE_CHALLENGE_SUCCESS,
  BUY_CHALLENGE_SUCCESS,
  VALIDATE_CHALLENGE_SUCCESS,
} from "../constants/actionTypes";

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
  },
};

export default function seasonReducer(state = initialState.seasons, action) {
  const {
    challengeInDetail,
    allTeamsOnSeason,
    activeChallenges: activeChallengesInState,
    activeTeam: activeTeamInState,
  } = state;
  const { payload: { challengeSubscription } = {} } = action;

  switch (action.type) {
    case GET_ALL_SEASONS_SUCCESS:
      return {
        ...state,
        allSeasons: action.payload.allSeasons,
      };

    case GET_ALL_SEASON_TEAMS_SUCCESS:
      return {
        ...state,
        allTeamsOnSeason: action.payload.allTeams,
        activeTeamInSeason: action.payload.allTeams.some(
          item => item.id === activeTeamInState.id
        ),
      };

    case SET_ACTIVE_SEASON:
      return {
        ...state,
        activeSeason: action.payload.activeSeason,
      };

    case GET_CHALLENGE_DETAILS_SUCCESS:
      return {
        ...state,
        challengeInDetail: action.payload.challenge,
      };

    case CLEAN_CHALLENGE_DETAIL:
      return {
        ...state,
        challengeInDetail: undefined,
      };

    case SET_CHALLENGES_LIST:
      return {
        ...state,
        activeChallenges: action.payload.activeChallenges,
      };

    case GET_TEAM_DETAILS_SUCCESS:
      return {
        ...state,
        teamInDetail: action.payload.team,
      };

    case SET_ACTIVE_TEAM:
      const {
        payload: { team },
      } = action;

      return {
        ...state,
        activeTeam: team,
        activeTeamInSeason:
          allTeamsOnSeason &&
          allTeamsOnSeason.some(item => item.id === team.id),
      };

    case BUY_CHALLENGE_SUCCESS:
      const challengeInDetailCloneBuy = clone(challengeInDetail);

      if (challengeInDetailCloneBuy.subscriptions) {
        challengeInDetailCloneBuy.subscriptions = [
          ...challengeInDetailCloneBuy.subscriptions,
          challengeSubscription,
        ];
      } else {
        challengeInDetailCloneBuy.subscriptions = [challengeSubscription];
      }

      return {
        ...state,
        challengeInDetail: challengeInDetailCloneBuy,
      };

    case VALIDATE_CHALLENGE_SUCCESS:
      const challengeInDetailClone = clone(challengeInDetail);
      const { subscriptions } = challengeInDetailClone;

      const subscriptionIndex = findIndex(
        propEq("id", challengeSubscription.id)
      )(subscriptions);
      const updatedSubscriptions = update(
        subscriptionIndex,
        challengeSubscription,
        subscriptions
      );
      challengeInDetailClone.subscriptions = updatedSubscriptions;

      return {
        ...state,
        challengeInDetail: challengeInDetailClone,
      };

    case CLOSE_CHALLENGE_SUCCESS:
      const challengeInDetailCloneClose = clone(challengeInDetail);
      const { subscriptions: subscriptionsClose } = challengeInDetailCloneClose;

      const subscriptionIndexClose = findIndex(
        propEq("id", challengeSubscription.id)
      )(subscriptionsClose);
      const updatedSubscriptionsClose = update(
        subscriptionIndexClose,
        challengeSubscription,
        subscriptionsClose
      );
      challengeInDetailCloneClose.subscriptions = updatedSubscriptionsClose;

      return {
        ...state,
        challengeInDetail: challengeInDetailCloneClose,
      };

    default:
      return state;
  }
}
