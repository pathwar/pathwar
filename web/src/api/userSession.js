/* eslint-disable no-unused-vars */
import { baseApi } from "./index";

export function pingUser() {
    return baseApi.get("/ping")
}
