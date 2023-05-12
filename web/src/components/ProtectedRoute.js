import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { Dimmer } from "tabler-react";
import { useAuth0 } from "@auth0/auth0-react";

const ProtectedRoute = ({ component: Component, ...rest }) => {
  const dispatch = useDispatch();
  const userSession = useSelector(state => state.userSession);
  const { isLoading, isAuthenticated, loginWithRedirect, getIdTokenClaims } = useAuth0()

  if (!isLoading && !isAuthenticated) {
    loginWithRedirect();
  } else if (!isLoading && isAuthenticated) {
    getIdTokenClaims().then((claims) => {
      console.log(claims.__raw);
    });
  }
  // useEffect(() => {
  //   const { access_token } = userSession;
  //   if (!access_token) {
  //     dispatch(setAuthSession());
  //   }
  // }, []);

    // if (userSession.isAuthenticated) {
    //   return (
    //   <Component {...rest} />
    //   );
    // }
  console.log(isLoading, isAuthenticated)
  return <Dimmer active loader />;
};

export default ProtectedRoute;
