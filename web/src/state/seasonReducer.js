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
  LIST_USER_TEAMS_INVITATIONS_SUCCESS,
  ACCEPT_TEAM_INVITATION_SUCCESS,
  DECLINE_TEAM_INVITATION_SUCCESS,
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
    userTeamsInvitations: undefined,
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

    case LIST_USER_TEAMS_INVITATIONS_SUCCESS:
      return {
        ...state,
        userTeamsInvitations: action.payload.userTeamsInvitations,
      }

    case ACCEPT_TEAM_INVITATION_SUCCESS:
      return {
        ...state,
        userTeamsInvitations: state.userTeamsInvitations ?
          state.userTeamsInvitations.filter(invitation => invitation.id !== action.payload.teamInviteID) :
          state.userTeamsInvitations
      }

    case DECLINE_TEAM_INVITATION_SUCCESS:
      return {
        ...state,
        userTeamsInvitations: state.userTeamsInvitations ?
          state.userTeamsInvitations.filter(invitation => invitation.id !== action.payload.teamInviteID) :
          state.userTeamsInvitations
      }

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
      const activeChallengesUpdated = activeChallengesInState.map(item => {
        if (item.id === challengeSubscription.season_challenge_id) {
          return {
            ...item,
            subscriptions: item.subscriptions
              ? [...item.subscriptions, challengeSubscription]
              : [challengeSubscription],
          };
        }

        return item;
      });

      if (challengeInDetail?.id === challengeSubscription.season_challenge_id) {
        const challengeInDetailCloneBuy = {
          ...challengeInDetail,
          subscriptions: challengeInDetail.subscriptions
            ? [...challengeInDetail.subscriptions, challengeSubscription]
            : [challengeSubscription],
        };

        return {
          ...state,
          activeChallenges: activeChallengesUpdated,
          challengeInDetail: challengeInDetailCloneBuy,
        };
      }

      return {
        ...state,
        activeChallenges: activeChallengesUpdated,
      };

    case VALIDATE_CHALLENGE_SUCCESS:
      const challengeInDetailCloneValidate = clone(challengeInDetail);
      const { subscriptions } = challengeInDetailCloneValidate;

      const subscriptionIndex = findIndex(
        propEq("id", challengeSubscription.id)
      )(subscriptions);
      const updatedSubscriptions = update(
        subscriptionIndex,
        challengeSubscription,
        subscriptions
      );
      challengeInDetailCloneValidate.subscriptions = updatedSubscriptions;

      return {
        ...state,
        challengeInDetail: challengeInDetailCloneValidate,
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
