import * as React from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";
import { Formik } from "formik";
import { LoginPage as TablerLoginPage } from "tabler-react";

import { performLoginAction } from "../../actions/userSession";

class LoginPage extends React.PureComponent {
  
  render() {
    const { performLoginAction } = this.props;

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
        onSubmit={(values) => {
          performLoginAction(values.email, values.password);
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
  performLoginAction: PropTypes.func
};

const mapStateToProps = () => ({});

const mapDispatchToProps = {
  performLoginAction: (email, password) => performLoginAction(email, password)
};

export default connect(
mapStateToProps,
mapDispatchToProps
)(LoginPage);
