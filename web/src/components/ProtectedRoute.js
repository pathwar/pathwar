import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { Dimmer } from "tabler-react";
import { useAuth0 } from "@auth0/auth0-react";
import {setAuthSession} from "../actions/userSession";

const ProtectedRoute = ({ component: Component, ...rest }) => {
  const dispatch = useDispatch();
  const userSession = useSelector(state => state.userSession);
  const {
    isLoading,
    isAuthenticated,
    loginWithRedirect,
    getIdTokenClaims,
  } = useAuth0()

  if (!isLoading && !isAuthenticated) {
    loginWithRedirect()
      .then(() => {getIdTokenClaims()
        .then((token) => {
          console.log(token.__raw);
          dispatch(setAuthSession(token.__raw))
        })
      })
  } else if (!isLoading && isAuthenticated && !userSession.accessToken) {
    getIdTokenClaims().then((token) => {
        console.log(token.__raw);
        dispatch(setAuthSession(token.__raw))
     })
  }

  if (userSession.isAuthenticated && userSession.accessToken) {
    return <Component {...rest} />;
  }

  return <Dimmer active loader />;
};

export default ProtectedRoute;
