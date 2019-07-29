import React from "react"
import PropTypes from "prop-types"
import { navigate } from "gatsby"
import { connect } from "react-redux";
import Keycloak from 'keycloak-js';

class ProtectedRoute extends React.PureComponent {

    constructor(props) {
        super(props);
        this.state = { keycloak: null, authenticated: false };
    }

    componentDidMount() {
        const keycloak = Keycloak("/keycloak.json");

        keycloak.init({onLoad: 'login-required'}).then(authenticated => {
          this.setState({ keycloak: keycloak, authenticated: authenticated })
        })

    }

    render() {
        const { component: Component, location, userSession, ...rest } = this.props;

        if (this.state.keycloak) {
          if (this.state.authenticated) {
            return <Component {...rest} />
          } else return (<h2>Auth error, please try again!</h2>)

        // if (!userSession.isAuthenticated  && location.pathname !== `/app/login`) {
        //   navigate(`/app/login`)
        //   return null
        }
        
        return (
          //  <Component {...rest} />
            <h2>Checking auth...</h2>
        );
    }
}

ProtectedRoute.propTypes = {
  component: PropTypes.any.isRequired,
}

const mapStateToProps = state => ({
  userSession: state.userSession
});

  
const mapDispatchToProps = {};
  
export default connect(
  mapStateToProps,
  mapDispatchToProps
)(ProtectedRoute);