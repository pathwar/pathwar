import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import Keycloak from "keycloak-js";
import { Dimmer } from "tabler-react";
import { toast } from "react-toastify";
import { setKeycloakSession } from "../actions/userSession";

const ProtectedRoute = ({ component: Component, ...rest }) => {
  const dispatch = useDispatch();
  const userSession = useSelector(state => state.userSession);

  useEffect(() => {
    const { activeKeycloakSession } = userSession;
    const keycloak = new Keycloak("/keycloak.json");
    const token = activeKeycloakSession && activeKeycloakSession.token;
    const refreshToken =
      activeKeycloakSession && activeKeycloakSession.refreshToken;

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
