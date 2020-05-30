import React from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import Keycloak from "keycloak-js";
import { Redirect } from "@reach/router";
import { Dimmer } from "tabler-react";
import { setKeycloakSession } from "../actions/userSession";

import styles from "../styles/layout/loader.module.css";

class ProtectedRoute extends React.PureComponent {
  async componentDidMount() {
    const { setKeycloakSession } = this.props;
    const keycloak = await Keycloak("/keycloak.json");

    keycloak
      .init({ onLoad: "login-required", checkLoginIframe: false })
      .then(authenticated => {
        setKeycloakSession(keycloak, authenticated);
      });
  }

  render() {
    const {
      component: Component,
      userSession,
      path,
      to,
      redirect,
      ...rest
    } = this.props;

    if (userSession.activeKeycloakSession) {
      if (redirect) {
        return <Redirect from={path} to={to} />;
      }
      if (userSession.isAuthenticated) {
        return <Component {...rest} />;
      } else return <h3>Auth error, please try again!</h3>;
    }

    return <Dimmer className={styles.dimmer} active loader />;
  }
}

ProtectedRoute.propTypes = {
  component: PropTypes.any.isRequired
};

const mapStateToProps = state => ({
  userSession: state.userSession
});

const mapDispatchToProps = {
  setKeycloakSession: (keycloakInstance, authenticated) =>
    setKeycloakSession(keycloakInstance, authenticated)
};

export default connect(mapStateToProps, mapDispatchToProps)(ProtectedRoute);
