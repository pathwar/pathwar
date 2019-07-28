import * as React from "react";
import Cookies from "js-cookie";
import { navigate } from "gatsby";
import { USER_SESSION_TOKEN_NAME } from "../constants/userSession";

class LogoutPage extends React.PureComponent {

    componentDidMount() {
        Cookies.remove(USER_SESSION_TOKEN_NAME);
        navigate("/app/login");
    }

    render() {
        return null;
    }
}

export default LogoutPage;