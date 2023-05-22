import React, { useEffect } from "react";
import { connect } from "react-redux";
import { navigate } from "gatsby";
import { logoutUser as logoutUserAction } from "../actions/userSession";
import { useAuth0 } from "@auth0/auth0-react";

const LogoutPage = ({ userSession, logoutUserAction }) => {
  const { logout } = useAuth0();

  useEffect(() => {
    const logoutOptions = { logoutParams: { returnTo: window.location.origin + "/challenges"} };

    if (!userSession.accessToken) {
      navigate("/");
    } else {
      logout(logoutOptions)
      logoutUserAction();
    }
  }, [userSession.accessToken, logout, logoutUserAction]);

  return null;
};

const mapStateToProps = state => ({
  userSession: state.userSession,
});

const mapDispatchToProps = {
  logoutUserAction: () => logoutUserAction(),
};

export default connect(mapStateToProps, mapDispatchToProps)(LogoutPage);
