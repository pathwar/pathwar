import * as React from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";
import { Formik } from "formik";
import Cookies from "js-cookie";
import { navigate } from "gatsby"

import { LoginPage as TablerLoginPage } from "tabler-react";
import { performLoginAction, pingUserAction } from "../actions/userSession";
import { USER_SESSION_TOKEN_NAME } from "../constants/userSession";

class LoginPage extends React.PureComponent {

  state = {
    redirectToReferrer: false
  }

  componentDidMount() {
    const { pingUserAction } = this.props;
    const token = Cookies.get(USER_SESSION_TOKEN_NAME);
    if(token) {
      pingUserAction();
    }
  }
  
  render() {
    const self = this;
    const { performLoginAction, userSession } = this.props;
    const { redirectToReferrer } = this.state;

    if (redirectToReferrer === true || userSession.isAuthenticated) {
      navigate("/app/dashboard")
    }

    return (
      <Formik
        initialValues={{
          email: "",
          password: "",
        }}
        validate={values => {
          let errors = {};
          if (!values.email) {
            errors.email = "Required";
          } else if (
            !/^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,4}$/i.test(values.email)
          ) {
            errors.email = "Invalid email address";
          }
          return errors;
        }}
        onSubmit={ async (values) => {
          await performLoginAction(values.email, values.password).then(() => {
            self.setState({ redirectToReferrer: true })
          });
        }}
        render={({
          values,
          errors,
          touched,
          handleChange,
          handleBlur,
          handleSubmit
        }) => (
          <TablerLoginPage
            onSubmit={handleSubmit}
            onChange={handleChange}
            onBlur={handleBlur}
            values={values}
            errors={errors}
            touched={touched}
          />
        )}
      />
    );
  }
}

LoginPage.propTypes = {
  userSession: PropTypes.object,
  performLoginAction: PropTypes.func,
  pingUserAction: PropTypes.func
};

const mapStateToProps = state => ({
  userSession: state.userSession
});

const mapDispatchToProps = {
  performLoginAction: (email, password) => performLoginAction(email, password),
  pingUserAction: () => pingUserAction()
};

export default connect(
mapStateToProps,
mapDispatchToProps
)(LoginPage);
