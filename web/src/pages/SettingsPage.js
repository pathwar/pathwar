import * as React from "react";
import { connect } from "react-redux";
import { Helmet } from "react-helmet";
import PropTypes from "prop-types";

import { Page, Grid, Button } from "tabler-react";

import siteMetaData from "../constants/metadata";
import { deleteAccount as deleteAccountAction } from "../actions/userSession";

class SettingsPage extends React.PureComponent {
  state = {
    isFetching: false,
  };

  deleteAccount = async reason => {
    const self = this;
    const { deleteAccountAction } = this.props;
    this.setState({ isFetching: true });
    deleteAccountAction(reason).then(response => {
      self.setState({ isFetching: false });
      return response;
    });
  };

  render() {
    const { isFetching } = this.state;
    const { title, description } = siteMetaData;

    return (
      <>
        <Helmet>
          <title>{title} - Settings</title>
          <meta name="description" content={description} />
        </Helmet>
        <Page.Content title={`Settings`}>
          <Grid.Row cards={true}>
            <Grid.Col xs={12} sm={12} lg={6}>
              <Button.List>
                <Button
                  onClick={() => this.deleteAccount("integration test")}
                  loading={isFetching}
                  color="primary"
                >
                  Delete Account
                </Button>
              </Button.List>
            </Grid.Col>
          </Grid.Row>
        </Page.Content>
      </>
    );
  }
}

SettingsPage.propTypes = {
  deleteAccountAction: PropTypes.func,
};

const mapStateToProps = () => ({});

const mapDispatchToProps = {
  deleteAccountAction: reason => deleteAccountAction(reason),
};

export default connect(mapStateToProps, mapDispatchToProps)(SettingsPage);
