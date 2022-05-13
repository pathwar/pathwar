import moment from "moment";
import React from "react";
import { Link } from "gatsby";
import { Card, Dimmer, Grid, TabbedCard, Tab } from "tabler-react";
import { FormattedMessage, useIntl } from "react-intl";
import iconPwn from "../../images/icon-pwn-small.svg";
import ShadowBox from "../ShadowBox";
import ChallengeCard from "../challenges/ChallengeCard";

const UserChallengesView = ({ challenges }) => {
  const intl = useIntl();

  const purchasedNotSolved =
    challenges &&
    challenges.filter(challenge => {
      const { subscriptions } = challenge;
      const isClosed = subscriptions && subscriptions[0].status === "Closed";

      if (subscriptions && !isClosed) {
        return challenge;
      }
    });

  const closed =
    challenges &&
    challenges.filter(challenge => {
      const { subscriptions } = challenge;
      const isClosed = subscriptions && subscriptions[0].status === "Closed";

      if (subscriptions && isClosed) {
        return challenge;
      }
    });

  if (!challenges) {
    return <Dimmer active loader />;
  }

  return (
    <ShadowBox>
      <div
        css={theme => ({
          ".nav-link.active": {
            color: theme.colors.primary,
            borderColor: theme.colors.primary,
            fontWeight: "bold",
          },
        })}
      >
        <TabbedCard
          initialTab={intl.formatMessage({ id: "UserChallengesView.tab1" })}
          style={{ margin: 0 }}
        >
          <Tab title={intl.formatMessage({ id: "UserChallengesView.tab1" })}>
            {purchasedNotSolved && purchasedNotSolved.length ? (
              <Grid.Row>
                {purchasedNotSolved.map(challenge => (
                  <Grid.Col xl={4} lg={6} sm={12} xs={12} key={challenge.id}>
                    <ChallengeCard challenge={challenge} columnMode={true} />
                  </Grid.Col>
                ))}
              </Grid.Row>
            ) : (
              <Grid.Col lg={12}>
                <p>
                  <FormattedMessage id="UserChallengesView.myChallengesEmpty" />{" "}
                  <Link to="/challenges">
                    <FormattedMessage id="UserChallengesView.myChallengesEmptyCTA" />
                  </Link>
                  .
                </p>
              </Grid.Col>
            )}
          </Tab>
          <Tab title={intl.formatMessage({ id: "UserChallengesView.tab2" })}>
            <Grid.Row>
              {closed && closed.length ? (
                closed.map(({ flavor, id, subscriptions }) => (
                  <Grid.Col lg={4} key={id}>
                    <Card>
                      <Card.Header>{flavor.challenge.name}</Card.Header>
                      <Card.Body>
                        <FormattedMessage id="UserChallengesView.done" />:{" "}
                        {moment(subscriptions[0].closed_at).calendar()} <br />
                        <FormattedMessage id="UserChallengesView.reward" />:{" "}
                        {flavor.validation_reward}{" "}
                        <img
                          css={{ display: "inline-block" }}
                          src={iconPwn}
                          className="img-responsive"
                        />
                      </Card.Body>
                    </Card>
                  </Grid.Col>
                ))
              ) : (
                <Grid.Col lg={12}>
                  <p>
                    <FormattedMessage id="UserChallengesView.doneEmpty" />
                  </p>
                </Grid.Col>
              )}
            </Grid.Row>
          </Tab>
        </TabbedCard>
      </div>
    </ShadowBox>
  );
};

export default UserChallengesView;
