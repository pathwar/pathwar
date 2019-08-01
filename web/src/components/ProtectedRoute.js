import React from "react"
import PropTypes from "prop-types"
import { connect } from "react-redux";
import Keycloak from 'keycloak-js';
import { Dimmer } from "tabler-react";
import { setKeycloakSession } from '../actions/userSession'

class ProtectedRoute extends React.PureComponent {

    async componentDidMount() {
        const { setKeycloakSession } = this.props;
        const keycloak = await Keycloak("/keycloak.json");

        keycloak.init({onLoad: 'login-required'}).success(authenticated => {
          setKeycloakSession(keycloak, authenticated);
        })

    }

    render() {
        const { component: Component, location, userSession, ...rest } = this.props;

        if (userSession.activeSession) {
          if (userSession.isAuthenticated) {
            return <Component {...rest} />
          } else return (<h3>Auth error, please try again!</h3>)
        }
      
        return (
          <Dimmer active loader style={{ marginTop: "50px" }} />
        );
    }
}

ProtectedRoute.propTypes = {
  component: PropTypes.any.isRequired,
}

const mapStateToProps = state => ({
  userSession: state.userSession
});

  
const mapDispatchToProps = {
  setKeycloakSession: (keycloakInstance, authenticated) => setKeycloakSession(keycloakInstance, authenticated)
};
  
export default connect(
  mapStateToProps,
  mapDispatchToProps
)(ProtectedRoute);