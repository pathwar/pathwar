import * as React from "react";
import { connect } from "react-redux"
import { Card, Table, Avatar } from "tabler-react";
import PropTypes from "prop-types";
import { fetchTeamsList as fetchTeamsListAction } from "../../actions/teams"

const TeamsRows = ({teamsList}) => {
    return teamsList.map((team) => {
        return (
            <Table.Row key={team.metadata.id}>
            <Table.Col className="w-1">
                <Avatar imageURL={team.gravatar_url} />
            </Table.Col>
            <Table.Col>{team.name}</Table.Col>
            <Table.Col>{team.locale}</Table.Col>
        </Table.Row>
        )
    })
                    
}

class TeamsCard extends React.PureComponent {

    componentDidMount() {
        const { fetchTeamsListAction } = this.props;
        fetchTeamsListAction();
    }
    
    render() {
        const { teams } = this.props;
        return (
            <Card>
                  <Card.Header>
                    <Card.Title>Teams</Card.Title>
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
                      </Table.Row>
                    </Table.Header>
                    <Table.Body>
                        {teams.teamsList && <TeamsRows teamsList={teams.teamsList} />}
                    </Table.Body>
                  </Table>
                </Card>
        )
    }
}

TeamsCard.propTypes = {
    teams: PropTypes.object,
    fetchTeamsListAction: PropTypes.func
};

const mapStateToProps = state => ({
    teams: state.teams
});

const mapDispatchToProps = {
    fetchTeamsListAction: () => fetchTeamsListAction()
};

export default connect(
	mapStateToProps,
	mapDispatchToProps
)(TeamsCard);