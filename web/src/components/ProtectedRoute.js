import React from "react";
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
    getAccessTokenSilently,
  } = useAuth0()

  if (!isLoading && !isAuthenticated) {
    loginWithRedirect({
      scope: "openid profile email roles",
    })
      .then(() => {getAccessTokenSilently()
        .then((token) => {
          console.log(token);
          dispatch(setAuthSession(token))
        })
      })
  } else if (!isLoading && isAuthenticated && !userSession.accessToken) {
    getAccessTokenSilently().then((token) => {
        console.log(token);
        dispatch(setAuthSession(token))
     })
  }

  if (userSession.isAuthenticated && userSession.accessToken) {
    return <Component {...rest} />;
  }

  return <Dimmer active loader />;
};

export default ProtectedRoute;
