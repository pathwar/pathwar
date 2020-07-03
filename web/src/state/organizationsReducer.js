import {
  SET_ACTIVE_ORGANIZATION,
  SET_ORGANIZATIONS_LIST,
} from "../constants/actionTypes";

const initialState = {
  organizations: {
    fetchingUserOrganizations: undefined,
    userOrganizationsList: undefined,
    activeOrganization: undefined,
    allOrganizationsList: undefined,
  },
};

export default function teamsReducer(
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

    default:
      return state;
  }
}
