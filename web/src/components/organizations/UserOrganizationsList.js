import * as React from "react";
import { connect } from "react-redux"
import { Card, Table, Avatar, Button } from "tabler-react";
import PropTypes from "prop-types";
import { fetchUserOrganizations as fetUserTeamsListAction } from "../../actions/organizations"

const OrgRows = ({teamsList, activeOrganization}) => {
    return teamsList.map((team) => {

        const isActive = team.id === activeOrganization.id;

        return (
            <Table.Row key={team.id}>
            <Table.Col className="w-1">
                <Avatar imageURL={team.gravatar_url} />
            </Table.Col>
            <Table.Col>{team.name}</Table.Col>
            <Table.Col>{team.locale}</Table.Col>
            {isActive && <Table.Col>
                Active
            </Table.Col>}
            {!isActive && <Table.Col>
                <Button color="info" size="sm">Set Active</Button>
            </Table.Col>}
        </Table.Row>
        )
    })

}

class UserOrganizationsList extends React.PureComponent {

    componentDidMount() {
        const { fetUserTeamsListAction } = this.props;
        fetUserTeamsListAction();
    }

    render() {
        const { userOrganizationsList, activeOrganization } = this.props;
        return (
            <Card>
                  <Card.Header>
                    <Card.Title>My Teams</Card.Title>
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
                        {userOrganizationsList && <OrgRows teamsList={userOrganizationsList} activeOrganization={activeOrganization} />}
                    </Table.Body>
                  </Table>
                </Card>
        )
    }
}

UserOrganizationsList.propTypes = {
    activeOrganization: PropTypes.object,
    userOrganizationsList: PropTypes.array,
    fetUserTeamsListAction: PropTypes.func
};

const mapStateToProps = state => ({
    userOrganizationsList: state.organizations.userOrganizationsList,
    activeOrganization: state.organizations.activeOrganization
});

const mapDispatchToProps = {
    fetUserTeamsListAction: () => fetUserTeamsListAction()
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(UserOrganizationsList);
