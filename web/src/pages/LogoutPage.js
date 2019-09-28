import * as React from "react";
import { connect } from "react-redux";
import { navigate } from "gatsby";


class LogoutPage extends React.PureComponent {

    componentDidMount() {
        const { userSession: { activeKeycloakSession } } = this.props;

        if (!activeKeycloakSession) {
            navigate("/");
        } else {
          activeKeycloakSession.logout();
        }
    }

    render() {
        return null;
    }
}

const mapStateToProps = state => ({
    userSession: state.userSession
});


const mapDispatchToProps = {};

export default connect(
mapStateToProps,
mapDispatchToProps
)(LogoutPage);
