import {
  CLEAN_CHALLENGE_DETAIL, CLEAN_ORGANIZATION_DETAILS,
  GET_ORGANIZATION_DETAILS_SUCCESS,
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
