import * as React from "react";
import { connect } from "react-redux";
import { Card, Table, Avatar, Button } from "tabler-react";
import PropTypes from "prop-types";
import {
  fetchOrganizationsList as fetchOrganizationsListAction,
  joinOrganization as joinOrganizationAction,
} from "../../actions/organizations";

const OrgRows = ({
  allOrganizationsList,
  userOrganizationsList,
  joinOrganization,
}) => {
  return allOrganizationsList.map(team => {
    const isUserOnTem = userOrganizationsList
      ? userOrganizationsList.find(userTeam => team.id === userTeam.id)
      : undefined;

    return (
      <Table.Row key={team.id}>
        <Table.Col className="w-1">
          <Avatar imageURL={team.gravatar_url} />
        </Table.Col>
        <Table.Col>{team.name}</Table.Col>
        <Table.Col>{team.locale}</Table.Col>
        {isUserOnTem && <Table.Col>Joined</Table.Col>}
        {!isUserOnTem && (
          <Table.Col>
            <Button
              color="info"
              size="sm"
              onClick={() => joinOrganization(team.id)}
            >
              Join
            </Button>
          </Table.Col>
        )}
      </Table.Row>
    );
  });
};

class AllOrganizationsCard extends React.PureComponent {
  componentDidMount() {
    const { fetchOrganizationsListAction } = this.props;
    fetchOrganizationsListAction();
  }

  render() {
    const { organizations, joinOrganizationAction } = this.props;
    return (
      <Card>
        <Card.Header>
          <Card.Title>All Organizations</Card.Title>
        </Card.Header>
        <Table
          cards={true}
          striped={true}
          responsive={true}
          className="table-vcenter"
        >
          <Table.Header>
            <Table.Row>
              <Table.ColHeader colSpan={2}>Name</Table.ColHeader>
              <Table.ColHeader>Locale</Table.ColHeader>
              <Table.ColHeader></Table.ColHeader>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {organizations.allOrganizationsList && (
              <OrgRows
                allOrganizationsList={organizations.allOrganizationsList}
                userOrganizationsList={organizations.userOrganizationsList}
                joinOrganization={joinOrganizationAction}
              />
            )}
          </Table.Body>
        </Table>
      </Card>
    );
  }
}

AllOrganizationsCard.propTypes = {
  organizations: PropTypes.object,
  fetchOrganizationsListAction: PropTypes.func,
  joinOrganizationAction: PropTypes.func,
};

const mapStateToProps = state => ({
  organizations: state.organizations,
});

const mapDispatchToProps = {
  fetchOrganizationsListAction: () => fetchOrganizationsListAction(),
  joinOrganizationAction: teamID => joinOrganizationAction(teamID),
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(AllOrganizationsCard);
