import React from "react"
import PropTypes from "prop-types"
import { connect } from "react-redux";
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

const mapStateToProps = state => ({
    userSession: state.userSession
});
  
const mapDispatchToProps = {};
  
export default connect(
  mapStateToProps,
  mapDispatchToProps
)(ProtectedRoute);