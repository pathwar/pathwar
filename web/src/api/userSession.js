/* eslint-disable no-unused-vars */
import { baseApi } from "./index";

export function getUserSession() {
  return baseApi.get("/user-session")
}
