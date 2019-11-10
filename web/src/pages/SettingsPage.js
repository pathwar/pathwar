import * as React from "react"
import { connect } from "react-redux"
import PropTypes from "prop-types"

import { Page, Grid, Button } from "tabler-react"

import { deleteAccount as deleteAccountAction } from "../actions/userSession"

class SettingsPage extends React.PureComponent {

  state = {
    isFetching: false,
    success: undefined
  }

  deleteAccount = async (reason) => {
    const self = this;
    const { deleteAccountAction } = this.props;
    this.setState({ isFetching: true, success: undefined })
    deleteAccountAction(reason).then((response) => {
      self.setState({ isFetching: false, success: true });
      return response;
    })
  }

  render() {
    const { isFetching, success } = this.state;

    return (
      <Page.Content title={`Settings`}>
        <Grid.Row cards={true}>
          <Grid.Col xs={12} sm={12} lg={6}>
            <Button.List>
              <Button onClick={() => this.deleteAccount("integration test")} loading={isFetching} color="primary">Delete Account</Button>
            </Button.List>
            {success && <p>Success!</p>}
          </Grid.Col>
        </Grid.Row>
      </Page.Content>
    )
  }
}

SettingsPage.propTypes = {
  deleteAccountAction: PropTypes.func,
}

const mapStateToProps = state => ({})

const mapDispatchToProps = {
  deleteAccountAction: reason => deleteAccountAction(reason),
}

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(SettingsPage)
