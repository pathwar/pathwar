import {
  SET_ACTIVE_ORGANIZATION,
  GET_USER_ORGANIZATIONS,
  SET_ORGANIZATIONS_LIST,
  GET_USER_ORGANIZATIONS_SUCCESS
} from '../constants/actionTypes';

const initialState = {
  organizations: {
      fetchingUserOrganizations: undefined,
      userOrganizationsList: undefined,
      activeOrganization: undefined,
      allOrganizationsList: undefined
  }
};

export default function teamsReducer(state = initialState.organizations, action) {

  switch (action.type) {
    case SET_ACTIVE_ORGANIZATION:
      return {
        ...state,
        activeOrganization: action.payload.team
      }

    case GET_USER_ORGANIZATIONS:
      return {
        ...state,
        fetchingUserOrganizations: true
      }

    case GET_USER_ORGANIZATIONS_SUCCESS:
      return {
        ...state,
        fetchingUserOrganizations: false,
        userOrganizationsList: action.payload.userOrganizationsList
      }

    case SET_ORGANIZATIONS_LIST:
      return {
        ...state,
        allOrganizationsList: action.payload.allOrganizationsList
      };

    default:
      return state;
  }
}
