import React from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import Keycloak from "keycloak-js";
import { Dimmer } from "tabler-react";
import { setKeycloakSession } from "../actions/userSession";

class ProtectedRoute extends React.PureComponent {
  componentDidMount() {
    const { setKeycloakSession, userSession } = this.props;
    const { activeKeycloakSession } = userSession;
    const keycloak = new Keycloak("/keycloak.json");
    const token = activeKeycloakSession && activeKeycloakSession.token;
    const refreshToken =
      activeKeycloakSession && activeKeycloakSession.refreshToken;

    keycloak
      .init({
        onLoad: "login-required",
        checkLoginIframe: false,
        token,
        refreshToken,
      })
      .then(authenticated => {
        setKeycloakSession(keycloak, authenticated);
      });
  }

  render() {
    const { component: Component, userSession, ...rest } = this.props;

    if (userSession.activeKeycloakSession) {
      if (userSession.isAuthenticated) {
        return <Component {...rest} />;
      } else return <h3>Auth error, please try again!</h3>;
    }

    return <Dimmer active loader />;
  }
}

ProtectedRoute.propTypes = {
  component: PropTypes.any.isRequired,
};

const mapStateToProps = state => ({
  userSession: state.userSession,
});

const mapDispatchToProps = {
  setKeycloakSession: (keycloakInstance, authenticated) =>
    setKeycloakSession(keycloakInstance, authenticated),
};

export default connect(mapStateToProps, mapDispatchToProps)(ProtectedRoute);
