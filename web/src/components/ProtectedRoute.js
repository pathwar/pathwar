import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import Keycloak from "keycloak-js";
import { Dimmer } from "tabler-react";
import { toast } from "react-toastify";
import { setKeycloakSession } from "../actions/userSession";

const ProtectedRoute = ({ component: Component, ...rest }) => {
  const dispatch = useDispatch();
  const userSession = useSelector(state => state.userSession);

  // We want to get token & refreshToken if they exist
  // If they don't exist, we want to "provider" login page
  // if they exist, we want to check if they are still valid
  // if they are not valid, we want to refresh them
  // if they are valid, we want to set them in the redux store
  // we also want to set authenticated to true
  useEffect(() => {
    // retrieve token
    const { activeKeycloakSession, access_token } = userSession;
    const keycloak = new Keycloak("/keycloak.json");
    const token = activeKeycloakSession && activeKeycloakSession.token;
    if (!access_token) {

    }
    const refreshToken =
      activeKeycloakSession && activeKeycloakSession.refreshToken;

    // if tokens don't exist, we want to redirect to "provider" login page
    // verify if token is expired
    // set token in cookie
    keycloak
      .init({
        onLoad: "login-required",
        checkLoginIframe: false,
        enableLogging: true,
        token,
        refreshToken,
      })
      .then(authenticated => {
        dispatch(setKeycloakSession(keycloak, authenticated));
      });

    keycloak.onTokenExpired = () => {
      keycloak
        .updateToken(30)
        .success(authenticated => {
          dispatch(setKeycloakSession(keycloak, authenticated));
        })
        .error(() =>
          toast.error(`SESSION EXPIRED! Please refresh the page.`, {
            autoClose: false,
            hideProgressBar: true,
          })
        );
    };
  }, []);

  if (userSession.activeKeycloakSession) {
    if (userSession.isAuthenticated) {
      return <Component {...rest} />;
    } else return <h3>Auth error, please try again!</h3>;
  }

  return <Dimmer active loader />;
};

export default ProtectedRoute;
