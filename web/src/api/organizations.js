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

export function postAnswerOrganizationInvitation(organizationInviteID, accept) {
  return baseApi.post(`/organization/invite/accept`, {
    organization_invite_id: organizationInviteID,
    accept: accept,
  });
}

export function postCreateOrganization(name, gravatarEmail) {
  return baseApi.post(`/organization`, {
    name: name,
    gravatar_email: gravatarEmail,
  });
}
