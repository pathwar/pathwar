import React from "react"
import PropTypes from "prop-types"
import { connect } from "react-redux";
import { Route, Redirect } from "react-router-dom";
import Keycloak from 'keycloak-js';

class ProtectedRoute extends React.PureComponent {

    constructor(props) {
        super(props);
        this.state = { keycloak: null, authenticated: false };
    }

    componentDidMount() {
        const keycloak = Keycloak({
            "realm": "Pathwar",
            "auth-server-url": "https://sso.pathwar.land/auth",
            "ssl-required": "external",
            "resource": "platform-front",
            "public-client": true,
            "confidential-port": 0
          });

        keycloak.init({onLoad: 'login-required'}).then(authenticated => {
          this.setState({ keycloak: keycloak, authenticated: authenticated })
        })

    }

    render() {
        const { component: Component, ...rest } = this.props;

        if (this.state.keycloak) {
            return (
                <Route {...rest} render={props => {
                    if (this.state.authenticated) return (
                        <Component {...props}/> 
                    ); else return (<h2>Auth error, please try again!</h2>)
                }} />
            )
        }
        
        return (
            <h2>Checking auth...</h2>
        );
    }
import { navigate } from "gatsby"

class ProtectedRoute extends React.PureComponent {

  render() {
    const { component: Component, location, userSession, ...rest } = this.props;
    
      if (!userSession.isAuthenticated  && location.pathname !== `/app/login`) {
        navigate(`/app/login`)
        return null
      }

    return <Component {...rest} />
  }
}

ProtectedRoute.propTypes = {
  component: PropTypes.any.isRequired,
}

const mapStateToProps = state => ({});
  
const mapDispatchToProps = {};
  
export default connect(
  mapStateToProps,
  mapDispatchToProps
)(ProtectedRoute);