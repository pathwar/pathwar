import * as React from "react";
import { connect } from "react-redux";
import PropTypes from "prop-types";

import {
  Page,
  Grid
} from "tabler-react";

class TeamDetailsPage extends React.PureComponent {

    componentDidMount() {

    }

    render() {
        const {} = this.props;
        return (
            <Page.Content title="Team XYZ">
                <Grid.Row cards={true}>
                  <Grid.Col xs={12} sm={12} lg={6}>
                    <h3>Teste</h3>
                  </Grid.Col>
                </Grid.Row>
              </Page.Content>
          );
    }
}

TeamDetailsPage.propTypes = {

};

const mapStateToProps = state => ({
});

const mapDispatchToProps = {
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(TeamDetailsPage);
