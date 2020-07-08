/* eslint-disable no-unused-vars */
import { baseApi } from "./index";

export function getAllOrganizations() {
  return baseApi.get(`/organizations`);
}
