import * as React from "react"
import { connect } from "react-redux"
import PropTypes from "prop-types"
import { Page, Grid, Button } from "tabler-react"
import { isNil } from "ramda"

import AllTeamsOnSeasonList from "../components/season/AllTeamsOnSeasonList"
import ChallengeList from "../components/challenges/ChallengeList"
import ValidateCouponButton from "../components/coupon/ValidateCouponButton"
import CreateTeamStampCard from "../components/team/CreateTeamStampCard"

import {
  fetchChallenges as fetchChallengesAction,
  fetchAllSeasonTeams as fetchAllSeasonTeamsAction,
  buyChallenge as buyChallengeAction,
  createTeam as createTeamAction,
} from "../actions/seasons"

class SeasonPage extends React.Component {
  componentDidUpdate(prevProps) {
    const {
      fetchAllSeasonTeamsAction,
      fetchChallengesAction,
      seasons: { activeSeason },
    } = this.props
    const {
      seasons: { activeSeason: prevActiveSeason },
    } = prevProps

    if (isNil(prevActiveSeason) && activeSeason) {
      fetchAllSeasonTeamsAction(activeSeason.id)
      fetchChallengesAction(activeSeason.id)
    }
  }

  render() {
    const {
      buyChallengeAction,
      createTeamAction,
      seasons: {
        activeSeason,
        activeChallenges,
        allTeamsOnSeason,
        activeTeamInSeason
      },
    } = this.props
    const name = activeSeason ? activeSeason.name : undefined

    return (
      <Page.Content title="Season" subTitle={name}>
        <Grid.Row>
          <Grid.Col lg={4} md={4} sm={4} xs={4}>
            <Button.List>
            <ValidateCouponButton />
            </Button.List>
          </Grid.Col>
        </Grid.Row>
        <hr />
        <Grid.Row>
          <Grid.Col xs={12} sm={3} lg={3}>
            <h4>Teams</h4>
            <CreateTeamStampCard
              activeSeason={activeSeason}
              createTeam={createTeamAction}
              activeTeamInSeason={activeTeamInSeason}
            />
            <AllTeamsOnSeasonList
              activeSeason={activeSeason}
              allTeamsOnSeason={allTeamsOnSeason}
            />
          </Grid.Col>
          <Grid.Col xs={12} sm={9} lg={9}>
            <h4>Challenges</h4>
            <ChallengeList
              challenges={activeChallenges}
              buyChallenge={buyChallengeAction}
            />
          </Grid.Col>
        </Grid.Row>
      </Page.Content>
    )
  }
}

SeasonPage.propTypes = {
  seasons: PropTypes.object,
  fetchChallengesAction: PropTypes.func,
}

const mapStateToProps = state => ({
  seasons: state.seasons,
  activeOrganization: state.organizations.activeOrganization,
})

const mapDispatchToProps = {
  fetchChallengesAction: seasonID => fetchChallengesAction(seasonID),
  fetchAllSeasonTeamsAction: seasonID => fetchAllSeasonTeamsAction(seasonID),
  buyChallengeAction: (seasonID, teamID) =>
    buyChallengeAction(seasonID, teamID),
  createTeamAction: (seasonID, name) => createTeamAction(seasonID, name),
}

export default connect(mapStateToProps, mapDispatchToProps)(SeasonPage)
