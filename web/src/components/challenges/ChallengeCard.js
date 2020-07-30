/* eslint-disable react/prop-types */
import React, { useState } from "react";
import { Button, Card, Grid, Tag } from "tabler-react";
import { useQueryParam, StringParam } from "use-query-params";
import { css } from "@emotion/core";
import ChallengeModal from "./ChallengeModal";

const cardTag = css`
  margin-right: 0.5rem;
`;

const notPurchased = css`
  box-shadow: none;
  opacity: 0.8;
`;

const closedStyle = css`
  opacity: 0.4;
  pointer-events: none;
  box-shadow: none;
`;

const ChallengeCard = ({ challenge }) => {
  const [modalQueryId, setModalQueryId] = useQueryParam("modal", StringParam);
  const [modalOpen, setModalOpen] = useState(modalQueryId === challenge.id);

  const { flavor, subscriptions, id: challengeID } = challenge;
  const purchased = subscriptions;
  const isClosed = subscriptions && subscriptions[0].status === "Closed";

  const onCloseModal = function() {
    setModalOpen(false);
    setModalQueryId(undefined);
  };

  const openModal = function() {
    setModalOpen(true);
    setModalQueryId(challengeID);
  };

  const cardColor = isClosed ? "red" : purchased ? "blue" : "gray";

  return (
    <Card css={[!purchased && notPurchased, isClosed && closedStyle]}>
      <Card.Status color={cardColor} side />
      <Card.Header>
        <Card.Title>{flavor.challenge.name}</Card.Title>
        <Card.Options>
          <Button
            color={isClosed ? "red" : "indigo"}
            size="sm"
            icon={isClosed ? "x-circle" : "compass"}
            onClick={openModal}
          >
            {isClosed ? "Closed" : "View"}
          </Button>
        </Card.Options>
      </Card.Header>
      <Card.Body>
        <>
          <Grid.Row>
            <Grid.Col width={12}>
              <p>
                {flavor.body ||
                  `Hello Ol'salt! Try to beat the ${flavor.challenge.name} challenge.
              Heave ho!`}
              </p>
            </Grid.Col>
          </Grid.Row>
          <Grid.Row>
            <Grid.Col width={12}>
              <Tag
                css={cardTag}
                color="lime"
                addOn={
                  flavor.purchase_price ? `$${flavor.purchase_price}` : "$0"
                }
                addOnColor="green"
              >
                Price
              </Tag>
              <Tag
                css={cardTag}
                color="yellow"
                addOn={flavor.validation_reward}
                addOnColor="yellow"
              >
                Reward
              </Tag>
            </Grid.Col>
          </Grid.Row>
          <ChallengeModal
            open={modalOpen}
            onClose={onCloseModal}
            challengeID={challengeID}
          />
        </>
      </Card.Body>
    </Card>
  );
};

export default React.memo(ChallengeCard);
