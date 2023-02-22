/* eslint-disable no-unused-vars */
import { baseApi } from "./index";

export function getAllOrganizations() {
  return baseApi.get(`/organizations`);
}

export function getOrganizationDetails(organizationID) {
  return baseApi.get(`/organization?organization_id=${organizationID}`);
}

export function postInviteUserToOrganization(organizationID, name) {
  return baseApi.post(`/organization/invite`, {
    organization_id: organizationID,
    user_id: name,
  });
}
