import {
  GET_ORGANIZATION_DETAILS_FAILED,
  GET_ORGANIZATION_DETAILS_SUCCESS,
  SET_ACTIVE_ORGANIZATION,
  SET_ORGANIZATIONS_LIST,
  SET_ORGANIZATIONS_LIST_FAILED,
  SET_USER_ORGANIZATIONS_LIST,
} from "../constants/actionTypes";
import {getAllOrganizations, getOrganizationDetails} from "../api/organizations";

export const setActiveOrganization = teamObjData => async dispatch => {
  dispatch({
    type: SET_ACTIVE_ORGANIZATION,
    payload: { team: teamObjData },
  });
};

export const fetchOrganizationsList = () => async dispatch => {
  try {
    const response = await getAllOrganizations();
    dispatch({
      type: SET_ORGANIZATIONS_LIST,
      payload: { allOrganizationsList: response.data.items },
    });
  } catch (error) {
    dispatch({ type: SET_ORGANIZATIONS_LIST_FAILED, payload: { error } });
  }
};

export const setUserOrganizationsList = organisations => async dispatch => {
  try {
    dispatch({
      type: SET_USER_ORGANIZATIONS_LIST,
      payload: { userOrganizationsList: organisations },
    });
  } catch (error) {
    dispatch({ type: SET_ORGANIZATIONS_LIST_FAILED, payload: { error } });
  }
}

export const fetchOrganizationDetail = organizationID => async dispatch => {
  try {
    const response = await getOrganizationDetails(organizationID);
    dispatch({
      type: GET_ORGANIZATION_DETAILS_SUCCESS,
      payload: { organization: response.data.item },
    });
  } catch (error) {
    dispatch({ type: GET_ORGANIZATION_DETAILS_FAILED, payload: { error } });
  }
}
