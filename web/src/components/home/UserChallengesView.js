import moment from "moment";
import React from "react";
import { Card, Grid } from "tabler-react";
import iconPwn from "../../images/icon-pwn-small.svg";

const UserChallengesView = ({ challenges }) => {
  const closed =
    challenges &&
    challenges.filter(challenge => {
      const { subscriptions } = challenge;
      const isClosed = subscriptions && subscriptions[0].status === "Closed";

      if (subscriptions && isClosed) {
        return challenge;
      }
    });

  return (
    <div>
      <h3>Challenges done</h3>
      <Grid.Row>
        {closed.map(({ flavor, id, subscriptions }) => (
          <Grid.Col lg={4} key={id}>
            <Card>
              <Card.Header>{flavor.challenge.name}</Card.Header>
              <Card.Body>
                Closed at: {moment(subscriptions[0].closed_at).calendar()}{" "}
                <br />
                Reward: {flavor.validation_reward}{" "}
                <img
                  css={{ display: "inline-block" }}
                  src={iconPwn}
                  className="img-responsive"
                />
              </Card.Body>
            </Card>
          </Grid.Col>
        ))}
      </Grid.Row>
    </div>
  );
};

export default UserChallengesView;
