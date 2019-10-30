import * as React from "react";
import { connect } from "react-redux";
import { navigate } from "gatsby";
import { logoutUser as logoutUserAction } from "../actions/userSession"


class LogoutPage extends React.PureComponent {

    componentDidMount() {
        const { userSession, logoutUserAction } = this.props;

        if (!userSession.activeKeycloakSession) {
            navigate("/");
        } else {
          userSession.activeKeycloakSession.logout();
          logoutUserAction();
        }
    }

    render() {
        return null;
    }
}

const mapStateToProps = state => ({
    userSession: state.userSession
});


const mapDispatchToProps = {
  logoutUserAction: () => logoutUserAction()
};

export default connect(
mapStateToProps,
mapDispatchToProps
)(LogoutPage);
