import * as React from "react";
import { performLogout } from "../../api/userSession";

class LogoutPage extends React.PureComponent {

    logoutRoutine() {
        performLogout();
    }

    render() {
        return (<React.Fragment>{ this.logoutRoutine() }</React.Fragment>)
    }
}

export default LogoutPage;