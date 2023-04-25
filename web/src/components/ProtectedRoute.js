import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { Dimmer } from "tabler-react";
import {setAuthSession} from "../actions/userSession";

const ProtectedRoute = ({ component: Component, ...rest }) => {
  const dispatch = useDispatch();
  const userSession = useSelector(state => state.userSession);

  useEffect(() => {
    const { access_token } = userSession;
    if (!access_token) {
      dispatch(setAuthSession());
    }
  }, []);

    if (userSession.isAuthenticated) {
      return <Component {...rest} />;
    }

  return <Dimmer active loader />;
};

export default ProtectedRoute;
