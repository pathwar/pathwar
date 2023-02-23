import {
  ACCEPT_ORGANIZATION_INVITATION_SUCCESS,
  CLEAN_CHALLENGE_DETAIL, CLEAN_ORGANIZATION_DETAILS,
  GET_ORGANIZATION_DETAILS_SUCCESS, LIST_USER_ORGANIZATIONS_INVITATIONS_SUCCESS, REJECT_ORGANIZATION_INVITATION_SUCCESS,
  SET_ACTIVE_ORGANIZATION,
  SET_ORGANIZATIONS_LIST,
  SET_USER_ORGANIZATIONS_LIST,
} from "../constants/actionTypes";

const initialState = {
  organizations: {
    fetchingUserOrganizations: undefined,
    userOrganizationsList: undefined,
    activeOrganization: undefined,
    allOrganizationsList: undefined,
    organizationInDetail: undefined,
    userOrganizationsInvitations: undefined,
  },
};

export default function organizationsReducer(
  state = initialState.organizations,
  action
) {
  switch (action.type) {
    case SET_ACTIVE_ORGANIZATION:
      return {
        ...state,
        activeOrganization: action.payload.team,
      };

    case SET_ORGANIZATIONS_LIST:
      return {
        ...state,
        allOrganizationsList: action.payload.allOrganizationsList,
      };

    case SET_USER_ORGANIZATIONS_LIST:
      return {
        ...state,
        userOrganizationsList: action.payload.userOrganizationsList,
      }

    case GET_ORGANIZATION_DETAILS_SUCCESS:
      return {
        ...state,
        organizationInDetail: action.payload.organization,
      }

    case LIST_USER_ORGANIZATIONS_INVITATIONS_SUCCESS:
      return {
        ...state,
        userOrganizationsInvitations: action.payload.userOrganizationsInvitations,
      }

    case ACCEPT_ORGANIZATION_INVITATION_SUCCESS:
      const newAcceptInvitations = state.userOrganizationsInvitations.filter(invitation => invitation.id !== action.payload.organizationInviteID)
      return {
        ...state,
        userOrganizationsInvitations: newAcceptInvitations,
      }

    case REJECT_ORGANIZATION_INVITATION_SUCCESS:
      const newRejectInvitations = state.userOrganizationsInvitations.filter(invitation => invitation.id !== action.payload.organizationInviteID)
      return {
        ...state,
        userOrganizationsInvitations: newRejectInvitations,
      }

      //Don't want to reload the page when we change the organization subpage
    case CLEAN_ORGANIZATION_DETAILS:
      return {
        ...state,
        /*organizationInDetail: undefined,*/
      };

    default:
      return state;
  }
}
