import {
  SET_ACTIVE_ORGANIZATION,
  GET_USER_ORGANIZATIONS_SUCCESS,
  GET_USER_ORGANIZATIONS_FAILED,
  SET_ORGANIZATIONS_LIST,
  SET_ORGANIZATIONS_LIST_FAILED,
  JOIN_ORGANIZATION_SUCCESS,
  JOIN_ORGANIZATION_FAILED,
} from "../constants/actionTypes";
import {
  getAllOrganizations,
  joinOrganization as joinOrganizationCall,
} from "../api/organizations";

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

export const joinOrganization = (userID, teamID) => async dispatch => {
  try {
    const response = await joinOrganizationCall(userID, teamID);
    dispatch({
      type: JOIN_ORGANIZATION_SUCCESS,
      payload: response.data,
    });

    dispatch(fetchOrganizationsList());
    dispatch(fetchUserOrganizations(userID));
  } catch (error) {
    dispatch({ type: JOIN_ORGANIZATION_FAILED, payload: { error } });
    alert("Join team failed, please try again!");
  }
};
