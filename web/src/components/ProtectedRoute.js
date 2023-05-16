import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { Dimmer } from "tabler-react";
import { useAuth0 } from "@auth0/auth0-react";
import {setAuthSession} from "../actions/userSession";


//TODO: update redux state with isAuthenticated value
//TODO: dispatch action to update redux state with userSession
//TODO: Then if userSession.isAuthenticated is true & userSession.accesToken is set, render the component
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
    loginWithRedirect()
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
