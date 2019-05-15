import * as React from "react";
import { connect } from "react-redux";
import { Route, Redirect } from "react-router-dom";

class ProtectedRoute extends React.PureComponent {


    render() {
        const { component: Component, userSession, ...rest } = this.props;
        
        return(
            <Route {...rest} render={props => (
                userSession.isAuthenticated 
                ? <Component {...props}/> 
                : <Redirect to={{
                    pathname: "/login",
                    state: { from: props.location }
                }} />
            )} />
        
        )
    }
}


const mapStateToProps = state => ({
    userSession: state.userSession
});
  
const mapDispatchToProps = {};
  
  export default connect(
  mapStateToProps,
  mapDispatchToProps
  )(ProtectedRoute);